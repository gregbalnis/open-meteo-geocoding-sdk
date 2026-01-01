```go
package geocoding

import (
	"context"
	"net/http"
)

// SearchOptions defines optional parameters for the Search method.
type SearchOptions struct {
	Count    int    // Default: 10, Max: 100
	Language string // Default: "en"
}

// Client defines the interface for the geocoding client.
type Client interface {
	// Search performs a location search by name.
	// ctx: Context for cancellation and timeouts.
	// name: The location name to search for.
	// opts: Optional parameters (pass nil for defaults).
	Search(ctx context.Context, name string, opts *SearchOptions) ([]Location, error)
}

// NewClient creates a new geocoding client with the given options.
func NewClient(opts ...Option) *ClientImpl {
	// Implementation details hidden in contract
	return &ClientImpl{}
}
```