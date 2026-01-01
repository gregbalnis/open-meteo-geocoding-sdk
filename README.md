# Open Meteo Geocoding SDK

A Go SDK for the Open Meteo Geocoding API.

## Installation

```bash
go get github.com/gregbalnis/open-meteo-geocoding-sdk
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	geocoding "github.com/gregbalnis/open-meteo-geocoding-sdk"
)

func main() {
	// 1. Initialize the client
	client := geocoding.NewClient()

	// 2. Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Define search options (optional)
	opts := &geocoding.SearchOptions{
		Count:    5,
		Language: "en",
	}

	// 4. Perform the search
	locations, err := client.Search(ctx, "Berlin", opts)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// 5. Process results
	for _, loc := range locations {
		fmt.Printf("Found: %s (%s) at %f, %f\n", 
			loc.Name, loc.CountryCode, loc.Latitude, loc.Longitude)
	}
}
```
