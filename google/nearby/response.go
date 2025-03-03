// https://places.googleapis.com/v1/places:searchNearby response
package nearby

type NearbyResponse struct {
	Places []Place `json:"places"`
}

type Place struct {
	Id                       string      `json:"id"`
	InternationalPhoneNumber string      `json:"internationalPhoneNumber"`
	FormattedAddress         string      `json:"formattedAddress"`
	Location                 Location    `json:"location"`
	GoogleMapsUri            string      `json:"googleMapsUri"`
	DisplayName              DisplayName `json:"displayName"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DisplayName struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}
