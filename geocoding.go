package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL          = "https://geocoding-api.open-meteo.com/v1/search"
	defaultTimeout          = 10 * time.Second
	defaultConcurrencyLimit = 10
)

// Client is the Open Meteo Geocoding API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	sem        chan struct{} // For concurrency limiting
}

// NewClient creates a new geocoding client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: defaultBaseURL,
		sem:     make(chan struct{}, defaultConcurrencyLimit),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Search performs a location search by name.
// ctx: Context for cancellation and timeouts.
// name: The location name to search for.
// opts: Optional parameters (pass nil for defaults).
func (c *Client) Search(ctx context.Context, name string, opts *SearchOptions) ([]Location, error) {
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	default:
		return nil, ErrConcurrencyLimitExceeded
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	q := u.Query()
	q.Set("name", name)
	q.Set("format", "json")

	// Set defaults if opts is nil, otherwise apply options
	count := 10
	language := "en"
	if opts != nil {
		if opts.Count > 0 {
			if opts.Count > 100 {
				return nil, fmt.Errorf("%w: count must be between 1 and 100", ErrInvalidParameter)
			}
			count = opts.Count
		}
		if opts.Language != "" {
			if len(opts.Language) != 2 {
				return nil, fmt.Errorf("%w: language must be a 2-letter code", ErrInvalidParameter)
			}
			language = opts.Language
		}
	}
	q.Set("count", fmt.Sprintf("%d", count))
	q.Set("language", language)

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Reason != "" {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var searchResp searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if searchResp.Error {
		return nil, &APIError{Reason: searchResp.Reason}
	}

	return searchResp.Results, nil
}
