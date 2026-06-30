# Design Specification: Tracking System Fixes

**Date**: 2026-06-26  
**Status**: Approved

---

## 1. Introduction & Goals

This specification details the implementation plan to address the critical and important issues identified during the code review of the web pages (sites), consents, and visitor tracking modules:
1. **Critical**: Replace the synchronous blocking third-party IP geolocation API call with an offline MaxMind GeoLite2 resolver that falls back gracefully.
2. **Critical**: Fix the session token validation logic in `ResolveSession` to reuse the incoming `sessionToken` when valid.
3. **Important**: Switch the rate limiter store from local memory to Redis for distributed support.
4. **Important**: Resolve concurrency race conditions in consent storage using an atomic database-level Upsert.

---

## 2. Detailed Designs

### 2.1 MaxMind GeoLite2 IP Geolocation Resolver
- **Package**: Create a new package `pkg/geoip` to encapsulate the geolocation reader.
- **Dependency**: Introduce `github.com/oschwald/geoip2-golang`.
- **Initialization**: 
  - On startup, the service attempts to open the database file specified by the environment variable `GEOIP_DB_PATH` (e.g. `/app/GeoLite2-City.mmdb`).
  - If the path is empty or the file does not exist, the resolver logs a warning and initializes a `NoopResolver` to avoid failing startup or blocking local development.
- **Integration**:
  - Update `internal/middlewares/tracking.go` to use this new resolver instead of making HTTP requests to `ip-api.com`.

### 2.2 Reusing Incoming Session Token in `ResolveSession`
- **File**: [internal/modules/visitors/v1/service.go](file:///Users/alouaniyoussef/Desktop/osoul/BACKEND/tracking/internal/modules/visitors/v1/service.go)
- **Logic**:
  - Check if `req.SessionToken` is provided and non-empty.
  - If it is a valid session token format/non-empty, use it for the session creation instead of unconditionally generating a new `uuid.New().String()`.

### 2.3 Redis Store for Rate Limiter
- **File**: [internal/middlewares/tracking.go](file:///Users/alouaniyoussef/Desktop/osoul/BACKEND/tracking/internal/middlewares/tracking.go)
- **Dependency**: Ensure `github.com/ulule/limiter/v3/drivers/store/redis` is imported.
- **Logic**:
  - Check if the global `state.Redis` client is initialized.
  - If initialized, initialize the `ulule/limiter` Redis store using the underlying client.
  - If not initialized (e.g., in unit tests), fall back to `memory.NewStore()`.

### 2.4 Atomic Upsert for Visitor Consent
- **Schema Modification**: 
  - Update [internal/ent/schemas/visitor_consent.go](file:///Users/alouaniyoussef/Desktop/osoul/BACKEND/tracking/internal/ent/schemas/visitor_consent.go) to declare a unique index on `(visitor_id, site_id)` using `entgo.io/ent/schema/index`.
  - Run `go generate ./internal/ent/...` to regenerate the Ent ORM database schema mapping and query builders.
- **Repository Modification**:
  - In [internal/modules/consents/v1/repository.go](file:///Users/alouaniyoussef/Desktop/osoul/BACKEND/tracking/internal/modules/consents/v1/repository.go), rewrite `SetVisitorConsent` to use the Ent `OnConflict` API to atomic-upsert based on the unique index.

---

## 3. Scope & Verification

- Verify that all unit and integration tests compile and pass.
- Verify that standard database auto-migrations run successfully on start.
