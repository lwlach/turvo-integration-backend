package models

// TurvoAuthRequest represents the request body for Turvo authentication
type TurvoAuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Scope        string `json:"scope"`
	Type         string `json:"type"`
}

// TurvoAuthResponse represents the response from Turvo authentication endpoint

//	{
//		"access_token": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
//		"token_type": "Bearer",
//		"refresh_token": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
//		"expires_in": 41363,
//		"scope": "read+trust+write",
//		"busId": "2",
//		"country": "US",
//		"temp": {
//		"id": 3923,
//		"key": "1510",
//		"value": "℉"
//		},
//		"distance": "mi",
//		"offset": 0,
//		"timezone": "America/Los_Angeles",
//		"busName": "Turvo",
//		"weight": {
//		"id": 0,
//		"key": "1520",
//		"value": "lb"
//		},
//		"busLogo": null,
//		"language": "en",
//		"type": "BUSUSER",
//		"userId": 3,
//		"version": "1.7.5",
//		"busLogoSmall": null,
//		"tenant_ref": "9Lmh7huZ",
//		"userOffset": 0,
//		"currency": "$",
//		"userTimezone": "",
//		"email": "example@xxx.com"
//		}
type TurvoAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Token expiration time in seconds
	Scope        string `json:"scope"`
	BusId        string `json:"busId"`
	Country      string `json:"country"`
}
