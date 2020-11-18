package subscription

import(
	"empowerthings.com/autoendpoint/vapid"
	"encoding/hex"
	l4g "code.google.com/p/log4go"
	fernet "fernet-go"
	"errors"
	"time"
	"empowerthings.com/autoendpoint/utils/uuid"
)
const (
	TTL_DURATION time.Duration = 1000 * 1000 * 1000 * 60 * 60 * 24 * 365 * 100
)
func  DecryptSubscription(encrypted_endpoint string, version int, crypto_key []*fernet.Key,debug bool) (uaid string, chid string, pub_key string, err error) {
	if debug && version == 1 {
		l4g.Info("Decrypting received endpoint.")
	}
	if debug && version == 2 {
		l4g.Info("Decrypting received VAPID endpoint.")
	}
	decoded_endpoint := fernet.VerifyAndDecrypt([]byte(encrypted_endpoint), TTL_DURATION, crypto_key)
	if decoded_endpoint == nil {
		decryption_err := errors.New("[Error]: Invalid Endpoint.")
		return "", "", "", decryption_err
	}
	hex_encoded_endpoint := hex.EncodeToString([]byte(decoded_endpoint))
	if version == 1 {
		uaid,chid = ExtractSubscription(hex_encoded_endpoint,debug)
		return 
	} else {
		uaid,chid,pub_key = vapid.ExtractVapidSubscription(hex_encoded_endpoint,debug)
		return
	}
}
func RepadSubscription(subscription string, debug bool) (paded_str string) {
	/*
		Adds padding to strings for base64 decoding.
		base64 encoding requires 'padding' to 4 octet boundries. 'padding' is a '='. So a string like 'abcde' is 5 octets, and would need to be padded to 'abcde==='.
	*/
	if debug {
		l4g.Info("Repading received encrypted endpoint.")
	}
	padings := len(subscription) % 4
	if padings != 0 {
		for i := 0; i < padings; i++ {
			subscription = subscription + "="
		}
	}
	paded_str = subscription
	return
}
func ExtractSubscription(subscription string, debug bool) (uaid string, chid string) {
	if debug {
		l4g.Info("Exracting chid and uaid from the received CassandraDB result object")
	}
	uaid = subscription[0:32]
	unformatted_chid := subscription[32:64]
	chid = uuid.FormatId(unformatted_chid,debug)
	return
}
