//go:build integration

package geocoding

import (
	"context"
	"testing"
	"time"
)

func TestIntegration_Search(t *testing.T) {
	client := NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	locations, err := client.Search(ctx, "Berlin", nil)
	if err != nil {
		t.Fatalf("Integration search failed: %v", err)
	}

	if len(locations) == 0 {
		t.Fatal("Expected at least one result for 'Berlin', got 0")
	}

	found := false
	for _, loc := range locations {
		if loc.Name == "Berlin" && loc.CountryCode == "DE" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find Berlin, DE in results")
	}
}
