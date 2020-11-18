/*******************************************************
 * Copyright (C) 2018  KAIOS Technologies INC.
 * 
 * The first Golang library to decode and verify VAPID JWT
 * 
 *******************************************************/

package vapid
import(
	
	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/model"
	"strings"
	"errors"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"crypto/elliptic"
	"crypto/ecdsa"
	"math/big"
	"encoding/hex"
	"crypto"
	"encoding/json"
	"time"
	"empowerthings.com/autoendpoint/utils/uuid"
	
)

var debug bool
var curve256 = elliptic.P256()


func Verify_AppServer_ID(headers model.Vapid_Headers, public_key_in_headers string, processed_key []byte) (err error) {	
	if debug {	
		l4g.Info("Verifying application server's identity.")
	}
	AUTH_SCHEME := [7]string{"webpush", "Webpush","Bearer","bearer","Hawk","hawk", "WebPush"}
	//Getting the public key of application server form the headers of POST request.
	authorization_header := headers.Authorization
	if len(authorization_header) == 0 {
		err = errors.New("Received VAPID request without authorization header.")
		return err
	}
	/*
 From authorization header we need to check three things: 

 [1] That Authorization scheme is webpush, or Bearer or bearer or Hawk or hawk. If none of them is included in the header, that will all cause authentication  fail.
 [2] That the application server's key included in the POST request header, maches with the one obtained by fernet-decrypting the subscription information.
 [3] That the JWT included in this header is a valid one. We will be using
     the public_key extracted in the step above to test/validate/decode the JWT.
	  */	
	auth_header_seperated := strings.Fields(authorization_header)
	/*
  Length should be no more than 2, one for the auth scheme, and one for the JWT. A valid header will look like:
 "Authorization: webpush <JWT_TOKEN>"   
 "Authorization: Webpush <JWT_TOKEN>"
 "Authorization: WebPush <JWT_TOKEN>"
 "Authorization: bearer  <JWT_TOKEN>"
 "Authorization: Bearer  <JWT_TOKEN>"
 "Authorization: Hawk    <JWT_TOKEN>"
 "Authorization: hawk    <JWT_TOKEN>"

*/
	if len(auth_header_seperated) != 2 {		
		err = errors.New("Invalid VAPID authorization header content.")
		return err
	}
	var found bool
	found = false
	//Searching in the defined set of schemes for any matches.If not,return an error.
	for i:= range AUTH_SCHEME {
		if auth_header_seperated[0] == AUTH_SCHEME[i] {
			found = true
		}	
	}
	if found == false{
		err_msg := fmt.Sprintf("Invalid VAPID authorization scheme received: %s",auth_header_seperated[0])
		err = errors.New(err_msg)
		return err	
	}
	//Extract the JWT part.
	token := auth_header_seperated[1]
	if token == "" {
		err_msg := fmt.Sprintf("VAPID request did not include JWT")
		err = errors.New(err_msg)
		return err
	}
	if debug {	
		l4g.Info("Validating VAPID JWT")
		l4g.Info(token)
	}
	err = validate_vapid_token(token,processed_key)
	if err!= nil {
		return err
	}
	if debug {
		l4g.Info("VAPID JWT is valid.")
	}
	return nil
}
// Verifies signature of VAPID JWT.
func validate_vapid_token (token string, public_key []byte) (err error) {
	// First step: Check token expiration.
	token_arr := strings.Split(token, ".") //Splits the three parts of jwt.
	if len(token_arr) !=3 { 
		err =errors.New("Wrong JWT format")
		return err		
	}
	base64_decoded_body, err := base64.RawURLEncoding.DecodeString(token_arr[1])
	if err != nil {
		l4g.Error("decode error: %s", err)
		return err
	}
	var token_body model.VAPID_JWT_BODY
	err = json.Unmarshal(base64_decoded_body, &token_body)
	expiration := token_body.Expiration
	// Check if token is expired.
	now:= time.Now()
	if now.Unix() > int64(expiration) {
		err = errors.New("Expired VAPID JWT")
		return err
		
	}
	//We need to process the public key (We can not just use it directly to verify the signature of JWT).
	// Baically, we need to trasform it into the *ecdsa.PublicKey structure.
	x,y := elliptic.Unmarshal(curve256, public_key)
	// Extract the point coordinates on the elliptic graph. Any point on the graph can be the public key. Golang's ecdsa library
	// requires the value of x and y as to build a public key instance which can be "understood" and used futher to validate the
	// data which was previously signed with the VAPID private key (Private key is the elliptic curve itslef). 
	//	l4g.Debug("value of X is: %s", x)
	//	l4g.Debug("value of Y is: %s", y)
	pubkey := ecdsa.PublicKey{Curve: curve256, X: x, Y: y}
	if debug{
		l4g.Info("Public key object has been built successfully")
	}
	sig_material, r, s, err := extract_signature(token)
	if err != nil{
		return err
	}
	//	l4g.Debug("Signature material is: %s", sig_material)
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(sig_material))
	valid := ecdsa.Verify(&pubkey,hasher.Sum(nil), s,r)
	if valid == false {
		err = errors.New("Invalid signature")
		return err
	}
	return nil
}

/*
get_label() method basically fetches the lable value included in header.
Example: 

str1 := "dh=BJ4QcSWOYelDxnwLkqoHDY80n8ZpcDtDgUoYFv5,p256ecdsa=BF93fPhAnJYn-QdqlFz4CbdVDSxWR"

get_label("dh", str1) >> BJ4QcSWOYelDxnwLkqoHDY80n8ZpcDtDgUoYFv5

get_label("P256", str1) >> BF93fPhAnJYn-QdqlFz4CbdVDSxWR

*/

func GetLabel(key string, body string) (label string) {	
	index := strings.Index(body, key)
	label = ""
	sub_body := body[index+ len(key) +1:]
	if index != -1 {
		if endidx := strings.IndexAny(sub_body,";,"); endidx != -1 {
			label = sub_body[:endidx]
		}else {
			label = sub_body
		}
		return label
	}else{
		// Not found !	
		return ""	
 	}
} 

func extract_signature(token string) (sig_material string, r,s *big.Int,err error) {
	
	rbig := new(big.Int)
	sbig := new(big.Int)
	
	r = rbig
	s = sbig
	if debug {
		l4g.Info("Extracting VAPID signature." )
		
	}
	if !strings.Contains(token,".") {	
		err =errors.New("Wrong JWT format")
		return "",r,s,err
	}
	token_arr := strings.Split(token, ".") //Splits the three parts of jwt.
	if len(token_arr) !=3 { 
		err =errors.New("Wrong JWT format")
		return "",r,s,err		
	}
	/*Based on VAPID specs, verification will happen as follows:
	 
[1] Received VAPID JWT is like any normal JWT which consists of three parts:
<header>.<body>.<signature>

[2] We need to extract signature material. Signature material is the first two parts combined and offcourse spearated by the dot.

[3] We need to extract two values from the signature (or from the third part). After decoding it to URL Base64 format, the first 
    part would be the first half of the result and the second part would be the second half, both hex encoded. Now each of those
    two hex-encoded parts will have to be converted to base 16 integer literal and returned as a big int.
    Note: The reason why I chosed big.Int type for those two integer literals is that even int64 would overflow if we assign the 
    result to it.       

[4] Now the the first three steps above, we have a) signature material  b)the big int values of the two parts.We are ready now
    to varify the signature.

    Verifying signature is simple and basic and can be understood by knowing how signing was done in the first place.


     ::::::::Sining a signature :::::::

    A) Hash the signature material.

       SHA256(signature_material)

    B) Sign the data with your elliptic curve (Private key) !

    s,r = sign (private key, hashed_signature_material) 




     ::::::::Verifying a signature :::::::   The inverse of Signing a signature.

    A) s,r = extract_r_and_s_values(signature)
    B) Verify_signature(  public_key ,SHA256(signature_material),s,r)

    That's it !

 */
	sig_material = token_arr[0] + "." + token_arr[1]  // first part and second part of JWT form the signature material together.	
	//	l4g.Debug("Third part (Signature) to decode: >%v< ", token_arr[2])
	base64_encoded_signature, err := base64.RawURLEncoding.DecodeString(token_arr[2])
	if err != nil {
		return  "",r,s,err
	}
	base64_encoded_signature, err = base64.RawURLEncoding.DecodeString(token_arr[2])
	if err != nil {
		return  "",r,s,err
	}
	if len(base64_encoded_signature)!= 64 {	
		err =errors.New("Invalid Signature.Length is not correct.")
		return "",r,s,err						
	}
	// Extracting s ans r values from the base64 decoded signature.
	
	//	l4g.Debug("Extracting r and s values from the base64 decoded signature.")
	
	first_part := []byte(base64_encoded_signature[:32])
	second_part:= []byte(base64_encoded_signature[32:])

	unfinished_s := make([]byte, hex.EncodedLen(len(first_part)))
	hex.Encode(unfinished_s, first_part)

	unfinished_r:= make([]byte, hex.EncodedLen(len(second_part)))
	hex.Encode(unfinished_r, second_part)
	
	// l4g.Debug("Unfinished s is %s: ", string(unfinished_s))
	// l4g.Debug("Unfinished r is:%s ", string(unfinished_r))

	rbig, ok := rbig.SetString(string(unfinished_r), 16)
        if !ok {
		err = errors.New("Cannot Parse string value of r to big int.")
		return "",r,s,err
	}
	sbig, ok = sbig.SetString(string(unfinished_s), 16)
        if !ok {
		err = errors.New("Cannot Parse string value of s to big int.")
		return "",r,s,err
	}
	r= rbig
	s= sbig
	// l4g.Debug("finished s is %s: ", string(s))
	// l4g.Debug("finished r is:%s ", string(r0)
	err = nil
	return  
} 

func DecipherKey(raw_key string) (result []byte, err error) {		
	if debug{
		l4g.Info("Deciphring public key.")
	}
	result,err = base64.RawURLEncoding.DecodeString(raw_key)
	if err !=nil {
		result = nil
		return result, err
	}
	key_len := len(result)
	if key_len ==65 &&  result[0] == '\x04'{	
		if debug{
			l4g.Info("CASE 1")
		}
		return result, nil 
	}
	//Key format is "raw"
	if key_len == 64 {
		result =  append( result, '\x04')
		return  result , nil
		if debug{
			l4g.Info("CASE 2")
		}
	}
	var equals bool 
	equals = false
	if result[0] == '0' && result[1]== 'v' && result[2]== '0'{
		equals = true
	} 
	if key_len == 88 && equals {
		if debug{
			l4g.Info("CASE 3")
		}
		result  := result[len(result)-64:]
		return result,nil	
	}
	err = errors.New("Unknown public key format specified")
	result = nil
	return result,err
}


func Varify_PublicKey(headers model.Vapid_Headers, pub_key string, publlic_key_in_header string, processed_key []byte) (err error) {	
	if debug {
		l4g.Info("Comparing public keys..")
	}
	//[A] SHA256 hash it.  (Defined in FIPS 180-4)
	sha256_pub_key := sha256.New()
	sha256_pub_key.Write([]byte(processed_key))	
	//[B] "Hexlify" the result !
	hex_encoded_pubkey := make([]byte, hex.EncodedLen(len(sha256_pub_key.Sum(nil))))
	hex.Encode(hex_encoded_pubkey, sha256_pub_key.Sum(nil))
	//[C] Compare between hex_encoded_pubkey  and the public key which was previously fernet-decoded
	if debug {
		l4g.Info("Comparing public keys")
	}
	if strings.Compare(string(hex_encoded_pubkey), pub_key) != 0 {		
		if debug {
			l4g.Error("[ERROR] Key Mismatch !")
			err = errors.New("Key Mismatch")
			return err
		}
	}
	if debug{
		l4g.Info("Public Keys matched.")
	}
	// Public keys matched !!
	return nil
}

func  ExtractVapidSubscription(subscription string, debug bool) (uaid string, chid string, pub_key string) {
	if debug {
		l4g.Info("Extracting chid, uaid and application server's public key from the received VAPID subscription data.")
	}
	if debug {
		l4g.Info("Decrypted Vapid endpoint information: %s ", subscription)
	}
	uaid = subscription[0:32]
	unformated_chid := subscription[32:64]
	chid = uuid.FormatId(unformated_chid,debug)
	pub_key = subscription[64:]
	return
}
func SetDebug( debug_flag bool) {
	debug = debug_flag
}
