package entity

// Business information
type Business struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	GoogleMapUrl string `json:"google_map_url"`
}
