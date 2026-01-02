package geocoding

import (
	"context"
	"encoding/json"
	"errors"

	// "errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearch_Basic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/search" {
			t.Errorf("Expected path /v1/search, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "Berlin" {
			t.Errorf("Expected name=Berlin, got %s", r.URL.Query().Get("name"))
		}
		if r.URL.Query().Get("format") != "json" {
			t.Errorf("Expected format=json, got %s", r.URL.Query().Get("format"))
		}

		resp := searchResponse{
			Results: []Location{
				{
					ID:          2950159,
					Name:        "Berlin",
					Latitude:    52.52437,
					Longitude:   13.41053,
					CountryCode: "DE",
					Country:     "Deutschland",
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL + "/v1/search"))
	locations, err := client.Search(context.Background(), "Berlin", nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(locations) != 1 {
		t.Fatalf("Expected 1 location, got %d", len(locations))
	}

	loc := locations[0]
	if loc.Name != "Berlin" {
		t.Errorf("Expected name Berlin, got %s", loc.Name)
	}
	if loc.CountryCode != "DE" {
		t.Errorf("Expected country code DE, got %s", loc.CountryCode)
	}
}

func TestSearch_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  true,
			"reason": "Parameter count must be between 1 and 100.",
		})
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL))
	_, err := client.Search(context.Background(), "Berlin", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.Reason != "Parameter count must be between 1 and 100." {
		t.Errorf("Unexpected error reason: %s", apiErr.Reason)
	}
}

func TestSearch_NetworkError(t *testing.T) {
	client := NewClient(WithBaseURL("http://invalid-url"))
	_, err := client.Search(context.Background(), "Berlin", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestSearch_OptionsValidation(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name    string
		opts    *SearchOptions
		wantErr bool
	}{
		{
			name: "Valid options",
			opts: &SearchOptions{
				Count:    5,
				Language: "de",
			},
			wantErr: false,
		},
		{
			name: "Invalid count > 100",
			opts: &SearchOptions{
				Count: 101,
			},
			wantErr: true,
		},
		{
			name: "Invalid language length",
			opts: &SearchOptions{
				Language: "eng",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We don't need a real server for validation tests as validation happens before request
			// But Search calls url.Parse(c.baseURL) first.
			// We can use a dummy URL.
			_, err := client.Search(context.Background(), "Berlin", tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearch_OptionsApplied(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("count") != "5" {
			t.Errorf("Expected count=5, got %s", r.URL.Query().Get("count"))
		}
		if r.URL.Query().Get("language") != "fr" {
			t.Errorf("Expected language=fr, got %s", r.URL.Query().Get("language"))
		}
		_ = json.NewEncoder(w).Encode(searchResponse{Results: []Location{}})
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL))
	_, _ = client.Search(context.Background(), "Paris", &SearchOptions{
		Count:    5,
		Language: "fr",
	})
}

func TestClient_WithHTTPClient(t *testing.T) {
	hc := &http.Client{Timeout: 5 * time.Second}
	client := NewClient(WithHTTPClient(hc))
	if client.httpClient != hc {
		t.Error("Expected custom HTTP client to be set")
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		Reason: "Something went wrong",
	}
	if err.Error() != "api error: Something went wrong" {
		t.Errorf("Expected error string 'api error: Something went wrong', got '%s'", err.Error())
	}
}

func TestSearch_EmptyName(t *testing.T) {
	client := NewClient()
	_, err := client.Search(context.Background(), "", nil)
	if err == nil {
		t.Fatal("Expected error for empty name, got nil")
	}
	if !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("Expected ErrInvalidParameter, got %v", err)
	}
}
