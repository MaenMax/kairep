package jwttoken

import (
//	"empowerthings.com/cumulis/config"
	"github.com/SermoDigital/jose/jwt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
//	"fmt"
	"errors"

)


var(
	_rsaPub  interface{}
	_valtor  *jwt.Validator
	_crypto  crypto.SigningMethod
)

func Init() error {



	var claims jwt.Claims
	claims=jwt.Claims{}
	claims["scp"]=nil
	claims["typ"]=nil
	claims["uid"]=nil
	claims["aid"]=nil
	claims["did"]=nil
	claims["pid"]=nil
	claims["jti"]=nil
	
// ES256
	_crypto = jws.GetSigningMethod("ES256") //For VAPID JWT
	
	if _crypto==nil {
		return errors.New("Failed to find Signing Method during initialization ...")
	}


	_valtor=&jwt.Validator{}
	_valtor.Expected=claims
	_valtor.EXP=0 // No leeway allowed
	_valtor.NBF=0 // No leeway allowed
	_valtor.Fn=nil	
	

	return nil

}

