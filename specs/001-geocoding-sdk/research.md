# Research: Open Meteo Geocoding SDK

**Feature**: `001-geocoding-sdk`
**Date**: 2026-01-01

## 1. Open Meteo Geocoding API

**Endpoint**: `https://geocoding-api.open-meteo.com/v1/search`

**Parameters**:
- `name` (string, required): Search term.
- `count` (int, optional): 1-100, default 10.
- `language` (string, optional): 2-char code, default "en".
- `format` (string, optional): "json" or "protobuf". We will hardcode "json".

**Response Structure**:
```json
{
  "results": [
    {
      "id": 2950159,
      "name": "Berlin",
      "latitude": 52.52437,
      "longitude": 13.41053,
      "elevation": 74.0,
      "country_code": "DE",
      "country": "Deutschland",
      "admin1": "Berlin",
      "admin2": "",
      "admin3": "Berlin, Stadt",
      "admin4": "Berlin"
    }
  ],
  "generationtime_ms": 0.5
}
```

**Error Structure**:
```json
{
  "error": true,
  "reason": "Parameter count must be between 1 and 100."
}
```

**Decision**:
- Map `results` to `[]Location`.
- Map `error: true` to a custom `APIError` type.
- Ignore `generationtime_ms` in the public API (internal use only if needed).

## 2. Concurrency Limiting

**Requirement**: Max 10 concurrent requests per instance.

**Options**:
1.  **Weighted Semaphore (`golang.org/x/sync/semaphore`)**: Robust, but adds a dependency.
2.  **Buffered Channel**: Simple, standard library only. `make(chan struct{}, 10)`.

**Decision**: **Buffered Channel**.
- **Rationale**: Keeps dependencies minimal (Constitution Principle I). Sufficient for the requirement.
- **Implementation**:
    ```go
    type Client struct {
        sem chan struct{}
        // ...
    }
    
    func (c *Client) Search(...) {
        select {
        case c.sem <- struct{}{}:
            defer func() { <-c.sem }()
        default:
            return nil, ErrConcurrencyLimitExceeded
        }
        // ...
    }
    ```
    *Correction*: The requirement says "return an error if this limit is exceeded". The `default` case in `select` achieves this non-blocking behavior. If blocking was required, we'd just wait on the channel. The prompt says "return an error if this limit is exceeded", so non-blocking is correct.

## 3. Timeouts

**Requirement**: Default 10s timeout.

**Decision**:
- The `NewClient` constructor will set a default `http.Client` with `Timeout: 10 * time.Second`.
- Users can override this via `WithHTTPClient` option.
- `Search` method will respect the passed `context.Context`. If the context is cancelled or times out before the HTTP client timeout, the request aborts.

## 4. JSON Parsing

**Decision**:
- Use standard `encoding/json`.
- Define a private `searchResponse` struct to match the API envelope.
- Define a public `Location` struct for the items in `results`.
- Unmarshal into `searchResponse`, then return `searchResponse.Results`.
