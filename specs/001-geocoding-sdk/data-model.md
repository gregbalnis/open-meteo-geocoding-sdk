# Data Model: Open Meteo Geocoding SDK

**Feature**: `001-geocoding-sdk`

## Public Types

### Location

Represents a single geographical location returned by the API.

```go
type Location struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	// Optional fields can be added later, keeping it minimal for now as per spec
}
```

### Client

The main entry point for the SDK.

```go
type Client struct {
	httpClient *http.Client
	baseURL    string
	sem        chan struct{} // For concurrency limiting
}
```

### Options

Functional options for configuring the client.

```go
type Option func(*Client)
```

## Internal Types

### searchResponse

Matches the JSON envelope from the API.

```go
type searchResponse struct {
	Results []Location `json:"results"`
	Error   bool       `json:"error,omitempty"`
	Reason  string     `json:"reason,omitempty"`
}
```

## Error Types

```go
var (
	ErrConcurrencyLimitExceeded = errors.New("concurrency limit exceeded")
	ErrInvalidParameter         = errors.New("invalid parameter")
)

type APIError struct {
	Reason string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error: %s", e.Reason)
}
```
