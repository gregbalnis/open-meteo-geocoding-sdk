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
