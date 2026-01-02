# Open Meteo Geocoding SDK

<p align="center">
    <a href="https://github.com/gregbalnis/open-meteo-geocoding-sdk/actions/workflows/release.yml"><img src="https://github.com/gregbalnis/open-meteo-geocoding-sdk/actions/workflows/release.yml/badge.svg" alt="Release"/></a>
    <a href="https://github.com/gregbalnis/open-meteo-geocoding-sdk/blob/main/LICENSE"><img src="https://img.shields.io/github/license/gregbalnis/open-meteo-geocoding-sdk" alt="License"/></a>
</p>

A Go SDK for the [Open Meteo Geocoding API](https://open-meteo.com/en/docs/geocoding-api).

## Installation

```bash
go get github.com/gregbalnis/open-meteo-geocoding-sdk
```

**Requirements**: Go 1.21 or later

## Quick Start

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

## API Reference

See [GoDoc](https://pkg.go.dev/github.com/gregbalnis/open-meteo-geocoding-sdk) for complete API documentation.

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Check coverage (requires 80%)
make coverage

# Clean artifacts
make clean
```

## License

This project is licensed under the terms of the MIT open source license. Please refer to the [LICENSE](https://github.com/gregbalnis/open-meteo-geocoding-sdk/blob/main/LICENSE) file for the full terms.

## Contributing

Contributions welcome! Please open an issue or pull request.