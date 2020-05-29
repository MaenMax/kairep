package jwttoken

import (
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"

	//	"empowerthings.com/cumulis/utils/uuid"
	//	"time"

	"crypto/ecdsa"
)

// GenToken returns a JWT using the specified claims
func CheckToken(encToken string, public_key ecdsa.PublicKey) error {
	var token jwt.JWT
	var err error

	if token, err = jws.ParseJWT([]byte(encToken)); err != nil {
		return err
	}

	if err = token.Validate(&public_key, _crypto, _valtor); err != nil {
		return err
	}

	return nil
}

