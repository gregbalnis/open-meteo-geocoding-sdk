---
description: "Task list for Open Meteo Geocoding SDK implementation"
---

# Tasks: Open Meteo Geocoding SDK

**Input**: Design documents from `/specs/001-geocoding-sdk/`
**Prerequisites**: plan.md, spec.md, data-model.md, research.md

**Tests**: Unit tests are included in the implementation phases.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel
- **[Story]**: Which user story this task belongs to (US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Project Initialization)

**Purpose**: Initialize the Go module and project structure.

**Note**: Rate limiting is out-of-scope for V1 per spec assumptions.

- [x] T001 Initialize Go module in `go.mod`
- [x] T002 Create project file structure (`geocoding.go`, `models.go`, `options.go`, `errors.go`, `geocoding_test.go`)

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define core data structures, errors, and the client configuration pattern.

- [x] T003 Define `Location` and `searchResponse` structs in `models.go`
- [x] T004 Define `APIError` and sentinel errors (`ErrConcurrencyLimitExceeded`, `ErrInvalidParameter`) in `errors.go`
- [x] T005 Define `Client` struct and `NewClient` constructor with concurrency limit initialization in `geocoding.go`
- [x] T006 [P] Implement Functional Options (`WithHTTPClient`, `WithBaseURL`) in `options.go`

## Phase 3: User Story 1 - Basic Location Search (Priority: P1)

**Goal**: Enable searching for a location by name and retrieving coordinates.
**Independent Test**: `TestSearch_Basic` in `geocoding_test.go` verifies a successful search returns expected coordinates.

- [x] T007 [US1] Implement `Search` method signature and URL construction in `geocoding.go`
- [x] T008 [US1] Implement concurrency limiting logic (non-blocking select) in `Search` method in `geocoding.go`
- [x] T009 [US1] Implement HTTP request execution with context and response handling in `geocoding.go`
- [x] T010 [US1] Implement JSON parsing and error mapping (handling `APIError`) in `geocoding.go`
- [x] T011 [US1] Create unit tests for basic search scenarios (success, network error, api error) in `geocoding_test.go`

## Phase 4: User Story 2 - Search Configuration (Priority: P2)

**Goal**: Allow filtering results by count and language.
**Independent Test**: `TestSearch_Options` in `geocoding_test.go` verifies query parameters are correctly applied.

- [x] T012 [US2] Define `SearchOptions` struct in `models.go`
- [x] T013 [US2] Update `Search` method to apply `SearchOptions` query parameters (`count`, `language`) in `geocoding.go`
- [x] T014 [US2] Add unit tests for search configuration and validation (including `count < 1`, `count > 100`, empty `name`) in `geocoding_test.go`

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final quality checks.

- [x] T015 Update `README.md` with installation instructions and usage examples from `quickstart.md`
- [x] T016 Verify unit test coverage exceeds 80% using `go test -cover`
- [x] T017 Verify no linting issues: run `go fmt` and `golangci-lint run` on all files
- [x] T018 [US1, US2] Create integration tests against live API in `geocoding_test.go` (Berlin search, Paris with count=1, network timeout scenarios)

## Dependencies

- **US1** depends on **Phase 2** (Models, Client)
- **US2** depends on **US1** (Search method existence)

## Implementation Strategy

1.  **Setup & Foundation**: Establish the types and client structure first.
2.  **MVP (US1)**: Implement the core search logic. This is the minimum usable product.
3.  **Enhancement (US2)**: Add configuration options to the existing search method.
4.  **Polish**: Finalize documentation and linting.
