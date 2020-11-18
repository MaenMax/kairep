package model

type Headers struct {	
	Encoding           string           `json:"encoding"` 
	Encryption         string           `json:"encryption"` 
	CryptoKey          string           `json:"crypto_key"`
}
type Vapid_Headers struct {	
	Encoding           string           `json:"encoding"` 
	Encryption         string           `json:"encryption"` 
	CryptoKey          string           `json:"crypto_key"`
	Authorization      string           `json:"authorization"` 
}
type T_CryptoKey struct {
	Dh               string             `json:"dh"`
	P256ecdsa        string             `json:"p256ecdsa"`
}
