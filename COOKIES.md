# Tracking cookies and session identifiers

The visitor-creation endpoint sets three tracking cookies after
`visitors/v1.Service.CreateSession` returns successfully:

```go
middlewares.SetTrackingCookies(
    c,
    response.SessionToken,
    response.UUID,
    response.DeviceUUID,
    h.appCfg.SessionExpiration,
    h.appCfg.VisitorExpiration,
    h.appCfg.CookiesDomain,
)
```

`response.UUID` is the **visitor UUID**. It is not the UUID of the
`visitor_sessions` database row.

| Cookie | Database value represented | Created when | Purpose |
| --- | --- | --- | --- |
| `session_token` | `visitor_sessions.session_token` | A new visitor session is created | Identifies one active session. |
| `visitor_uuid` | `visitors.uuid` | A new visitor is created | Identifies the visitor across sessions. |
| `device_uuid` | `devices.uuid` | A new device is created | Identifies the device record across sessions. |

## Value generation and persistence

All three values are generated with `uuid.New().String()` (UUID v4):

- A visitor UUID is created in `createVisitor` and stored in `visitors.uuid`.
- A device UUID is created in `createDevice` and stored in `devices.uuid`.
- A session token is created for a new session and stored in the unique,
  non-empty `visitor_sessions.session_token` column.

Each `visitor_sessions` row also has its own distinct `uuid`; the session
token is the client-facing bearer value.

When the session service creates a session, it deactivates previous sessions
for the same visitor, sets `is_active` to `true`, and, when configured, sets
`ended_at` to `now + SESSION_EXPIRATION`. Defaults are 7,200 seconds (two
hours) for `SESSION_EXPIRATION` and 31,536,000 seconds (one year) for
`VISITOR_EXPIRATION`.

An active token is looked up by `session_token`. A token whose `ended_at` is
in the past is treated as invalid and is asynchronously marked inactive.

## CreateSession decision flow

`RequestTracking` reads incoming `visitor_uuid`, `session_token`, and
`device_uuid` cookies into `TrackingInfo`. `CreateSession` then behaves as
follows:

1. If both `visitor_uuid` and `session_token` are present, it attempts to
   retrieve that active, unexpired session. On success it returns the existing
   visitor UUID, session token, and the device UUID associated with the
   session.
2. If a visitor UUID is present but the token is absent, invalid, inactive, or
   expired, it creates a new session for that visitor, using the visitor's
   last active/last-session device and acquisition attributes.
3. If there is no usable visitor cookie, the device is resolved by the
   required `fingerprint_hash`. A missing fingerprint match creates both a new
   visitor and device; a matching fingerprint reuses that device and its
   associated visitor. The service then creates or updates a session as
   appropriate.

## Cookie attributes

`SetTrackingCookies` sets all three cookies with:

- `Path=/`
- configured `Domain` (default `osoul.local`)
- `Secure=true`
- `HttpOnly=true`
- `SameSite=None`

`session_token` uses `SESSION_EXPIRATION` as `Max-Age`. `visitor_uuid` and
`device_uuid` use `VISITOR_EXPIRATION` as `Max-Age`.

The cookies are not readable by JavaScript because they are `HttpOnly`; they
require HTTPS because they are `Secure`; and `SameSite=None` permits
cross-site sending, subject to browser third-party-cookie restrictions.

## Current implementation caveats

- `device_uuid` is issued and read by `RequestTracking`, but it is not used to
  find a device, choose a session, or validate ownership. Device resolution is
  based on `fingerprint_hash`.
- Session ownership validation verifies that the session's `visitor_id` has
  the supplied `visitor_uuid`. Its `deviceID` argument is currently unused.
- `GetExistingSession` retrieves a session and visitor but does not compare
  `session.VisitorID` with the resolved visitor ID before returning a response.
  Other validation paths do perform this ownership comparison.
