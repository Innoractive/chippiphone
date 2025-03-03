package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Innoractive/chippiphone/cache"
	"github.com/Innoractive/chippiphone/entity"
	"github.com/Innoractive/chippiphone/google"
	"github.com/Innoractive/chippiphone/google/nearby"

	"github.com/gofiber/fiber/v2/log"
)

type SearchService struct{}

// Search for restaurants within a given area.
// - `area` - Address in staring. Stating of state/city is recommended. E.g. "mount austin, johor bahru"
// Result will be cached for 24 hours.
func (s *SearchService) Search(area string) ([]entity.Business, error) {
	/* Steps:
	- Normalize `area` string
	- Check cache for `area`
	- Resolve `area` to lat/lon
	- Search nearby resolved lat/lon
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
		lat, lng, err := area2LatLng(area)
		if err != nil {
			return "", fmt.Errorf("unable to determine location of %v: %v", area, err)
		}

		places, err := google.Nearby(lat, lng, os.Getenv("GOOGLE_API_KEY"))
		if err != nil {
			return "", fmt.Errorf("unable to search nearby %v: %v", area, err)
		}

		businesses := places2Businesses(places)

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

// Convert Nearby Search JSON to `Business` entity
func places2Businesses(places []nearby.Place) []entity.Business {
	var businesses []entity.Business
	for _, place := range places {
		businesses = append(businesses, entity.Business{
			Name:         place.DisplayName.Text,
			Phone:        place.InternationalPhoneNumber,
			Address:      place.FormattedAddress,
			GoogleMapUrl: place.GoogleMapsUri,
		})
	}

	return businesses
}

// Determine the center (lat, lng) of the area.
// If the area resolves to multiple locations, only first 1 is chosen.
func area2LatLng(area string) (lat float64, lng float64, err error) {
	addresses, err := google.Geocode(area, os.Getenv("GOOGLE_API_KEY"))
	if err != nil {
		return 0, 0, fmt.Errorf("google geocoding failed: %v", err)
	} else if len(addresses) == 0 {
		return 0, 0, nil
	}

	// XXX Pick the first one, there might be multiple matching places
	lat = addresses[0].Geometry.Location.Lat
	lng = addresses[0].Geometry.Location.Lng

	return lat, lng, err
}
