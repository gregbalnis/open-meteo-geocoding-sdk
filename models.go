package geocoding

// Location represents a single geographical location returned by the API.
type Location struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"` // first-order administrative division (https://www.geonames.org/export/codes.html)
	Admin2      string  `json:"admin2"` // second-order administrative division
	Admin3      string  `json:"admin3"` // third-order administrative division
	Admin4      string  `json:"admin4"` // fourth-order administrative division
}

// searchResponse matches the JSON envelope from the API.
type searchResponse struct {
	Results []Location `json:"results"`
	Error   bool       `json:"error,omitempty"`
	Reason  string     `json:"reason,omitempty"`
}

// SearchOptions defines optional parameters for the Search method.
type SearchOptions struct {
	// Count limits the number of results (default: 10, max: 100).
	Count int
	// Language specifies the language for location names (default: "en").
	Language string
}
