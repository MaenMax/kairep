package model

// The second part of VAPID JWT (BODY).

type VAPID_JWT_BODY struct {
	
	
	Audience       string      `json:"aud"`
	
	Sub            string      `json:"sub"`
	
	Expiration     int         `json:"exp"`
	
	
}
