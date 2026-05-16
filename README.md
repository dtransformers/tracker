# Tracking API

The purpose is to master the tracking of the website navigation to provide this service on the guest experience application.

## Vistors Detection

in order to ensure we track users, we need first of all to make sure we do have efficient tools to block crawlers

using **Cloudflare blocks AI crawlers and bots by default for new websites**, offering a one-click solution to protect content from unauthorized training data scraping. Users can manage or block AI bots (e.g., GPTBot, ClaudeBot) via **AI Crawl Control** in the dashboard. This service operates at the infrastructure level, detecting and stopping bots before they reach your server.

the visitors are confirmed once they reach the website, which means

- **Cloudflare** CAPTCHA passed ✅

## Tracking

once the users are on the site we start by communicating with the tracking SDK which allow the website to connect with the tracking module

### **Tracking Steps**

1. Visitor lands on website
2. Generate / retrieve anonymous_id (cookie)
3. Create or update Visitor
4. Capture request context (session)

## Authentication

since the tracking module is an open module, for every website that uses it, the mechanism used it’s a hybrid authentication process which merges:

- Public API Key: Each website has unique API Key, the API key is public
- HMAC signature: Each request will have HMAC fields which will be used to verify if the received requests are valid or not.

### Data Modeling:

#### Sites

```
TABLE sites {
  id bigint
  uuid string
  domaine_name string
  is_active string
  is_deleted string
  key string
  created_at timestamp
  updated_at timestamp
}
```

#### Visitors

```json

TABLE visitors {
  id bigint
  uuid string

  anonymous_id string
  primary_user_id string nullable

  first_seen_at timestamp
  last_seen_at timestamp

  is_bot boolean
  is_verified_human boolean

  created_at timestamp
  updated_at timestamp
}
```

#### Sites Visitors

```
TABLE site_visotrs {
  id bigint
  uuid string bigint
  visitor_id bigint
  session_id bigint
  site_id bigint

  created_at timestamp
  updated_at timestamp

}
```

#### Devices

```
TABLE devices {
  id bigint
  visitor_id bigint

  fingerprint_hash string

  device_type string

  os string

  address_ip string
  country string // (ISO Format)
  region string
  city string
  browser string
  language string
  timezone string

  first_seen_at timestamp
  last_seen_at timestamp
}
```

### Sessions:

```
TABLE visitor_sessions {
  id bigint
  visitor_id bigint
  device_id bigint

  session_token string

  started_at timestamp
  ended_at timestamp
  last_activity_at timestamp

  entry_page string
  exit_page string

  acquisition_source string
  acquisition_medium string
  acquisition_campaign string

  referrer_host string

  created_at timestamp
}
```

#### Visitors Events

```
TABLE visitor_events {
  id bigint

  visitor_id bigint
  session_id bigint
  device_id bigint

  type string
  name string

  page_url string
  referrer_url string

  properties jsonb

  created_at timestamp
}
```

#### Page Views

```
TABLE pageviews {
  id bigint

  visitor_id bigint
  session_id bigint

  path string
  title string

  viewed_at timestamp
}
```

#### Identity Links

```
TABLE visitor_identities {
  id bigint

  visitor_id bigint

  type string
  value string

  verified boolean
  confidence_score decimal

  created_at timestamp
}
```

#### Visitors Consent

```
visitor_consents {
  id
  visitor_id

  analytics_allowed boolean
  marketing_allowed boolean
  personalization_allowed boolean

  tracking_allowed boolean

  ad_storage_allowed boolean
  ad_personalization_allowed boolean

  consent_source string   -- cookie banner, login, POS, etc.
  consent_version string

  region string           -- EU, US, etc.

  created_at
  updated_at
}
```

## Consent Awareness:

In legal theory, modern tracking requires an informed "notice-and-choice" framework before data collection can begin. In practice, however, online service providers layer their data collection practices within dense, overly complex Terms of Service (ToS) and Privacy Policies (PP) that the vast majority of users accept without reading—a pattern famously dubbed "the biggest lie on the internet" (Bergram et al., 2020).

Research indicates that users typically prioritize immediate access to a service over their long-term data privacy, resulting in a systemic mismatch between user understanding and actual corporate data collection (Bergram et al., 2020; Raibulet & Zhang, 2025).

### Pathologies of Digital Consent

Consent models on the web suffer from deep architectural and behavioral defects, frequently categorized as "pathologies of consent" (Richards & Hartzog, 2021). These include:

- **Unwitting Consent:** Accepting terms because the sheer cognitive overhead of understanding what data is being tracked (such as cookies, background processes, or cached tracking vectors) is too high (Welzel et al., 2025).
- **Incapacitated Consent:** The systemic overload of encountering dozens of unique consent pop-ups daily, which triggers **"consent fatigue"**—causing users to blindly click "Accept All" to clear their screens.

### Emerging Vectors: Wearables & Cross-Device Tracking

As web-based tracking moves past standard desktop cookies, consent awareness faces new frontiers:

- **Cross-Device Authentication (XDAuth):** Modern workflows often decouple authentication across multiple devices (e.g., using a phone push notification or QR code to log into a browser). Studies reveal a massive gap here: over half of major web services fail to provide any context about the target device requesting authorization, meaning users frequently approve web sessions without a clear "right to know" what they are consenting to (Zhang, 2026).
- **Wearable Ecosystems:** Connected health applications and consumer wearables generate highly detailed biometric tracking data that is often mirrored or handled through web platforms. Consumer data protection groups frequently warn that these ecosystems lack explicit boundaries; background syncs and alternative permissions make revoking web-tracking consent incredibly difficult for the average user (Doherty, 2025; Welzel et al., 2025).
