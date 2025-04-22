// Execute Google APIs
package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Innoractive/chippiphone/google/geocode"
	"github.com/Innoractive/chippiphone/google/nearby"

	"github.com/gofiber/fiber/v2/log"
)

// Search for F&B within 500m radius of a given position.
// Uses [Nearby Search API](https://bit.ly/4hblK6t).
func Nearby(latitude float64, longitude float64, apiKey string) ([]nearby.Place, error) {
	// Ref: https://bit.ly/4hN95rH
	requestTemplate := `{
  "includedTypes": ["cafe", "cafeteria", "coffee_shop", "hamburger_restaurant", "ice_cream_shop", "juice_shop", "tea_house", "restaurant"],
  "maxResultCount": %d,
  "locationRestriction": {
    "circle": {
      "center": {
        "latitude": %f,
        "longitude": %f},
      "radius": %f
    }
  }
}`
	// Specifies the maximum number of place results to return. Must be between 1 and 20 (default) inclusive.
	maxResultCount := 20
	// Radius in meters
	radius := 500.0
	requestBody := fmt.Sprintf(requestTemplate, maxResultCount, latitude, longitude, radius)
	req, err := http.NewRequest(http.MethodPost, "https://places.googleapis.com/v1/places:searchNearby", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Errorf("Error making Nearby request: %v", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Goog-Api-Key", apiKey)
	// Ref: https://bit.ly/41kE7QA
	req.Header.Add("X-Goog-FieldMask", "places.displayName,places.formattedAddress,places.id,places.location,places.internationalPhoneNumber,places.googleMapsUri")

	// Create an HTTP client and send the request.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending Nearby request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading Nearby response body: %v", err)
		return nil, err
	}

	var nearbyResponse nearby.NearbyResponse
	err = json.Unmarshal(body, &nearbyResponse)
	if err != nil {
		log.Errorf("Failed to unmarshal Nearby JSON: %v", err)
		return nil, err
	}

	return nearbyResponse.Places, nil
}

// Convert address string to lat/lon using [Geocoding API](https://bit.ly/43e5bDP).
// *Restricted* to Malaysia only.
func Geocode(address string, apiKey string) ([]geocode.Address, error) {
	getUrl := "https://maps.googleapis.com/maps/api/geocode/json"

	queryParams := url.Values{}
	queryParams.Add("key", apiKey)
	queryParams.Add("address", address)
	queryParams.Add("components", "country:MY")

	resp, err := http.Get(getUrl + "?" + queryParams.Encode())
	if err != nil {
		log.Errorf("Error making Geocode request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading Geocode response body: %v", err)
		return nil, err
	}

	var geocodeResponse geocode.GeocodeResponse
	err = json.Unmarshal(body, &geocodeResponse)
	if err != nil {
		log.Errorf("Failed to unmarshal Geocode JSON: %v", err)
		return nil, err
	}

	if geocodeResponse.Status != "OK" {
		return nil, fmt.Errorf("Geocode API request failed: %s", geocodeResponse.Status)
	}

	return geocodeResponse.Results, nil
}
