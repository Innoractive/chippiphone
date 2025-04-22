package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/Innoractive/chippiphone/cache"
	"github.com/Innoractive/chippiphone/entity"
	"github.com/Innoractive/chippiphone/google"

	"github.com/gofiber/fiber/v2/log"
)

type SearchService struct{}

// Search for restaurants within a given area.
// - `area` - Address in staring. Stating of state/city is recommended. E.g. "mount austin, johor bahru"
// Result will be cached for 24 hours.
func (s *SearchService) Search(area string) ([]entity.Business, error) {
	/* Steps:
	- Normalize `area` string
	- Return cached result if available
	- Resolve `area` to lat/lng, this is the center.
	- Determine lat/lng of 8 points around this center
	- Search Businesses at each point, filter duplicates.
	- Cache result
	- Return result
	*/

	// Normalize `area`
	area = strings.ToLower(strings.TrimSpace(area))
	if len(area) < 8 {
		return nil, errors.New("area name must be at least 8 characters")
	}

	var businesses []entity.Business
	cache := cache.New()
	jsonString, err := cache.CachedGet(area, func() (string, error) {
		log.Infof("Creating cache for area: %v", area)

		// Resolve area to lat/lon
		center, err := area2LatLng(area)
		if err != nil {
			return "", fmt.Errorf("unable to determine location of %v: %v", area, err)
		}

		// Add 8 points around center
		ninePoints := GetSurroundingCoordinates(center.Lat, center.Lng, 1000)
		ninePoints = append(ninePoints, entity.Coordinate{
			Lat: center.Lat,
			Lng: center.Lng,
		})

		// Find businesses
		businesses, err = findBusinesses(ninePoints)
		if err != nil {
			return "", err
		}

		jsonBytes, err := json.Marshal(businesses)
		if err != nil {
			return "", fmt.Errorf("unable to JSON encode businesses info: %v", err)
		}
		return string(jsonBytes), nil
	}, 24*time.Hour) // Cached for 24 hours
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonString), &businesses)

	return businesses, err
}

// Find businesses at each coordinate. Due to search radius, there might be duplicates.
// Return deduplicated list of businesses.
// - `coordinates` - Array of coordinates to search.
func findBusinesses(coordinates []entity.Coordinate) ([]entity.Business, error) {
	// Key: nearby.Place.Id
	businesses := make(map[string]entity.Business)

	for _, coordinate := range coordinates {
		places, err := google.Nearby(coordinate.Lat, coordinate.Lng, os.Getenv("GOOGLE_API_KEY"))
		if err != nil {
			return nil, fmt.Errorf("unable to search nearby (%v, %v): %v", coordinate.Lat, coordinate.Lng, err)
		}

		// Filter out duplicate businesses
		for _, place := range places {
			if _, ok := businesses[place.Id]; ok {
				continue
			}

			// Convert Nearby Search JSON to `Business` entity
			businesses[place.Id] = entity.Business{
				Name:         place.DisplayName.Text,
				Phone:        place.InternationalPhoneNumber,
				Address:      place.FormattedAddress,
				GoogleMapUrl: place.GoogleMapsUri,
			}
		}
	}

	// Convert map to slice
	result := make([]entity.Business, 0, len(businesses))
	for _, business := range businesses {
		result = append(result, business)
	}

	return result, nil
}

// Determine the center (lat, lng) of the area.
// If the area resolves to multiple locations, only first 1 is chosen.
func area2LatLng(area string) (coordinate entity.Coordinate, err error) {
	addresses, err := google.Geocode(area, os.Getenv("GOOGLE_API_KEY"))
	if err != nil {
		return entity.Coordinate{}, fmt.Errorf("google geocoding failed: %v", err)
	} else if len(addresses) == 0 {
		return entity.Coordinate{}, fmt.Errorf("no matching address found")
	}

	// XXX Pick the first one, there might be multiple matching places
	return entity.Coordinate{
		Lat: addresses[0].Geometry.Location.Lat,
		Lng: addresses[0].Geometry.Location.Lng,
	}, nil
}

// Returns a slice of coordinates representing the points around the given center point.
// The distanceMeters parameter specifies the distance in meters from the center point to the surrounding points.
func GetSurroundingCoordinates(lat, lng float64, distanceMeters float64) []entity.Coordinate {
	// Earth's radius in meters
	earthRadius := 6378137.0

	// Convert distance to radians
	distRadians := distanceMeters / earthRadius

	// Convert degrees to radians
	latRad := lat * (math.Pi / 180)
	lngRad := lng * (math.Pi / 180)

	// Calculate coordinates in 8 directions
	coordinates := make([]entity.Coordinate, 8)

	// Define angles for 8 directions in radians (0 = north, π/4 = northeast, etc.)
	angles := []float64{0, math.Pi / 4, math.Pi / 2, 3 * math.Pi / 4, math.Pi, 5 * math.Pi / 4, 3 * math.Pi / 2, 7 * math.Pi / 4}

	for i, angle := range angles {
		// Calculate new position
		newLat := math.Asin(math.Sin(latRad)*math.Cos(distRadians) +
			math.Cos(latRad)*math.Sin(distRadians)*math.Cos(angle))

		newLng := lngRad + math.Atan2(math.Sin(angle)*math.Sin(distRadians)*math.Cos(latRad),
			math.Cos(distRadians)-math.Sin(latRad)*math.Sin(newLat))

		// Convert back to degrees
		coordinates[i].Lat = newLat * (180 / math.Pi)
		coordinates[i].Lng = newLng * (180 / math.Pi)
	}

	return coordinates
}
