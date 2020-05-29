package jwttoken

import (
	"encoding/pem"
	"io/ioutil"
	"errors"
//	"crypto/rsa"
	"crypto/x509"
)

func loadPubKey(key_file string) error {
	var rsaPub   interface{}

	fcont, err := ioutil.ReadFile(key_file)
	if err != nil {
		return err
	}

	keybuff,_ := pem.Decode(fcont)
	
	if keybuff==nil {
		return errors.New("Failed to decode content of public key file");
	}

	rsaPub, err = x509.ParsePKIXPublicKey(keybuff.Bytes)

	if err != nil {
		return err
	}

	_rsaPub=rsaPub

	return nil
}

