// https://maps.googleapis.com/maps/api/geocode/json reponse
package geocode

type GeocodeResponse struct {
	Results []Address `json:"results"`
	Status  string    `json:"status"`
}

type Address struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          Geometry           `json:"geometry"`
	PlaceId           string             `json:"place_id"`
	Types             []string           `json:"types"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Bound        Bound    `json:"bounds"`
	Location     Location `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Bound    `json:"viewport"`
}

type Bound struct {
	NorthEast Location `json:"northeast"`
	SouthWest Location `json:"southwest"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
