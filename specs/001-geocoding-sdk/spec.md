# Feature Specification: Open Meteo Geocoding SDK

**Feature Branch**: `001-geocoding-sdk`
**Created**: 2026-01-01
**Status**: Draft
**Input**: User description: "We are building an SDK that will support interactions with the Open Meteo Geocoding API. We will use open-access API, no apikey will be required or used. We'll always use JSON responses. Any Go program should be able to import this package without worrying about API integration detail."

## Clarifications

### Session 2026-01-01
- Q: Should network methods accept `context.Context`? → A: Yes, all network methods MUST accept `context.Context` as the first argument.
- Q: Should the SDK allow injecting a custom `*http.Client`? → A: Yes, the SDK configuration/constructor MUST allow injecting a custom `*http.Client`.
- Q: What configuration pattern should be used? → A: The Functional Options pattern (e.g., `NewClient(opts ...Option)`) MUST be used.
- Q: What should the search method return? → A: It MUST return a slice of locations `([]Location, error)`, abstracting the API envelope.
- Q: What data type should be used for coordinates? → A: `float64` MUST be used for latitude and longitude.

## User Scenarios & Testing

### User Story 1 - Basic Location Search (Priority: P1)

As a Go developer, I want to search for a location by name so that I can retrieve its geographical coordinates (latitude, longitude) and other metadata without handling HTTP requests manually.

**Why this priority**: This is the core functionality of the SDK. Without this, the package has no value.

**Independent Test**: Can be tested by creating a simple Go program that imports the package, calls the search function with a city name (e.g., "Berlin"), and asserts that valid coordinates are returned.

**Acceptance Scenarios**:

1. **Given** the SDK is initialized, **When** I search for "Berlin", **Then** I receive a list of locations containing "Berlin, Germany" with valid latitude and longitude.
2. **Given** the SDK is initialized, **When** I search for a non-existent location (e.g., "XyZ123Random"), **Then** I receive an empty list and no error.
3. **Given** the network is down, **When** I attempt a search, **Then** I receive a clear error indicating a network/connection issue.

---

### User Story 2 - Search Configuration (Priority: P2)

As a Go developer, I want to limit the number of results and specify the language of the results so that I can tailor the output to my application's needs.

**Why this priority**: Provides necessary flexibility for real-world applications (e.g., autocomplete dropdowns, localized apps).

**Independent Test**: Call the search function with specific options (e.g., count=1, language="de") and verify the response respects these constraints.

**Acceptance Scenarios**:

1. **Given** I want only the top result, **When** I search for "Paris" with `count=1`, **Then** I receive exactly one result.
2. **Given** I want German results, **When** I search for "Munich" with `language="de"`, **Then** the returned location name is "München".

## Functional Requirements

### 1. Client Interface
- The package MUST expose a Client or function to perform searches.
- The Client MUST be initialized using the Functional Options pattern (e.g., `NewClient(opts ...Option)`).
- The Client MUST allow configuration of a custom `*http.Client` (e.g., for timeouts or testing).
- The API MUST NOT require an API key configuration (as per Open Meteo specs).
- The default base URL MUST be `https://geocoding-api.open-meteo.com/v1/search`.

### 2. Search Functionality
- A method/function MUST exist to search by name (string).
- All network-bound methods MUST accept `context.Context` as the first argument.
- The search method MUST return `([]Location, error)`, unwrapping the API response envelope.
- It SHOULD support optional parameters for:
  - `count` (number of results, default 10, max 100).
  - `language` (response language, default "en").
  - `format` MUST be hardcoded to "json".

### 3. Data Models
- The SDK MUST define public Go structs representing the API response.
- Fields MUST include at minimum: `id`, `name`, `latitude` (float64), `longitude` (float64), `elevation`, `country`, `country_code`.
- JSON decoding MUST be handled internally.

### 4. Error Handling
- The SDK MUST return typed or descriptive errors for:
  - Network failures.
  - JSON decoding errors.
  - Invalid input parameters.

## Success Criteria

### Quantitative
- **Test Coverage**: Unit tests MUST cover > 80% of the code (per Constitution).
- **Performance**: Parsing a standard API response (10 items) MUST take < 1ms on standard hardware.

### Qualitative
- **Usability**: A user can implement a "Hello World" search in < 10 lines of code.
- **Correctness**: Integration tests against the live Open Meteo API return accurate data for reference cities (London, New York, Tokyo).
- **Compliance**: No API key is requested or sent in headers.

## Assumptions
- The Open Meteo Geocoding API remains free and open-access.
- The API schema (JSON structure) is stable.
- Rate limiting is handled by the consumer or the API (SDK will not implement complex backoff logic for V1).
