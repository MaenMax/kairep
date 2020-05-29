package hawk

import (
	"testing"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"empowerthings.com/cumulis/security/hawk"
	"net/url"
)

func Test_Decoding(t *testing.T) {

	key := "KQHdSdDAXvJefwYTfSJ1QETN5D3HTlhhn+cbhtEw1jE="

	converted, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Errorf("failed to convert test to base64 %v", err)
	}

	expected := []byte  {41,1,221,73,208,192,94,242,94,127,6,19,125,34,117,64,68,205,228,61,199,78,88,97,159,231,27,134,209,48,214,49}

	if len(converted) != len(expected) {
		t.Errorf("array lengths don't match converted %d vs expected %d", len(converted), len(expected))
	}

	for idx := 0; idx < len(expected); idx++ {

		if converted[idx] != expected[idx] {
			t.Errorf("Array elemenet %d doesn't match, converted %d vs expected %d", idx, converted[idx], expected[idx])
		}
	}
}

func Test_Decoding_Again(t *testing.T) {

	key := "5ZIU6PPCIxXtGm7M2od2vbp0otYZ45fQLE+9JoILbPA="

	converted, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Errorf("failed to convert test to base64 %v", err)
	}

	expected := []byte  {229,146,20,232,243,194,35,21,237,26,110,204,218,135,118,189,186,116,162,214,25,227,151,208,44,79,189,38,130,11,108,240}

	if len(converted) != len(expected) {
		t.Errorf("array lengths don't match converted %d vs expected %d", len(converted), len(expected))
	}

	for idx := 0; idx < len(expected); idx++ {

		if converted[idx] != expected[idx] {
			t.Errorf("Array elemenet %d doesn't match, converted %d vs expected %d", idx, converted[idx], expected[idx])
		}
	}
}

func Test_Mac_Sign(t *testing.T) {

	text := "testing 1234567890"
	key := []byte("werxhqb98rpaxn39848xrunpaw3489ruxnpa98w4rxn")
	h := hmac.New(sha256.New, key)
	h.Write([]byte(text))

	mac := base64.StdEncoding.EncodeToString(h.Sum(nil))
	expected := "q4dW/HICCanuP3kPF2naHyruTesDQU2tkmEUkh9aEqc="

	if expected != mac {
		t.Errorf("Mac %s isn't expected value %s", mac, expected)
	}
}

func Test_Text_Key(t *testing.T) {

	text := "testing1234567890testing1234567890"
	//key := "iLJgESHNcxaqSCCqH1z/FotkfhG8eycB0zooCzOmTbk="
	//mac_key_bytes, err := base64.StdEncoding.DecodeString(key)
	//if err != nil {
	//	t.Errorf("failed to convert test to base64 %v", err)
	//}

	h := hmac.New(sha256.New, []byte("secret"))
	h.Write([]byte(text))

	mac := base64.StdEncoding.EncodeToString(h.Sum(nil))
	expected := "KyNpcw5HmT1wet0YRaywdcctbyI3xXf9kpsWFICooFY="

	if expected != mac {
		t.Errorf("Mac %s isn't expected value %s", mac, expected)
	}
}

func Test_Not_B64_Key(t *testing.T) {

	text := "testing1234567890testing1234567890"
	key := "iLJgESHNcxaqSCCqH1z/FotkfhG8eycB0zooCzOmTbk="


	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(text))

	mac := base64.StdEncoding.EncodeToString(h.Sum(nil))
	expected := "BMZFP8Xlkc9DgrEKn8q0A6xcVzsRKTZIHPx0tcJ3VkQ="

	if expected != mac {
		t.Errorf("Mac %s isn't expected value %s", mac, expected)
	}
}

// func Canonical_Header(meth_or_code string, http_host string, url *url.URL, kid string, key []byte, nonce string,
// ts string, ext_data string, payload_hash string, is_tls_active bool, req_id string) (ch *CanonicalHeader, err error) {

func Test_Mac_With_Normalized(t *testing.T) {

	header := "hawk.1.header"
	meth := "GET"
	host := "api.dev.kaiostech.com"
	port := "8090"
	uri := "/v3.0/apps/summary"
	kid := "65msL/yeKyonWqW7ZS6AQGs6uI0="
	key := "UEUxgoI/HxkuJ6FqMN+6X19++88t43pP6s/1XOgh7Ek="
	nonce := "1nIu2NsYztjDO0k3DNPh"
	ts := "1513623037"

	chdr := hawk.CanonicalHeader {
		Kid:			kid,
		Key:            []byte(key),
		Mac:            "",
		Header:         header,
		Time:           ts,
		Nonce:          nonce,
		Meth_or_Status: meth,
		Uri:            uri,
		Host:           host,
		Port:           port,
		PayloadHash:    "",
		Ext:            "",
	}

	expected := "hawk.1.header\n" +
		"1513623037\n" +
		"1nIu2NsYztjDO0k3DNPh\n" +
		"GET\n" +
		"/v3.0/apps/summary\n" +
		"api.dev.kaiostech.com\n" +
		"8090\n" +
		"\n" +
		"\n"

	if chdr.ToString() != expected {
		t.Errorf("Canonical hdr %s doesn't match expected value %s", chdr.ToString(), expected)
	}
}

func Test_Mac_Preconstructed_Normalized(t *testing.T) {

	key := "PRLSGjwba5OyIFqqH7e5XM3kAdTnXO4crbrAH628CuA="

	converted, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Errorf("failed to convert test to base64 %v", err)
	}

	expected_bytes := []byte  {
		61,18,210,26,60,27,107,147,178,32,90,170,31,183,185,92,205,228,1,212,231,92,238,28,173,186,192,31,173,188,10,224}

	if len(converted) != len(expected_bytes) {
		t.Errorf("array lengths don't match converted %d vs expected %d", len(converted), len(expected_bytes))
	}

	for idx := 0; idx < len(expected_bytes); idx++ {

		if converted[idx] != expected_bytes[idx] {
			t.Errorf("Array elemenet %d doesn't match, converted %d vs expected %d", idx, converted[idx], expected_bytes[idx])
		}
	}

	// From js:
	// Hawk id="p33suMPiWtNlEfzknuC+Z+0mMn0=", ts="1513643773", nonce="7Wb8WRwx1Y6Uk6XOtTpyYin7ARzMi4",
	// mac="Tu+9DbhkJ3NTB8YeNdsOUjhSoLS9pYqYXsTER71dHcw="
	header := "hawk.1.header"
	meth := "GET"
	host := "api.dev.kaiostech.com"
	port := "8090"
	uri := "/v3.0/apps/summary"
	kid := "rjniBDZTk5XPQ4nRRURxLNlCGEw"
	nonce := "TO74yl3UnIQhoZAhMAKf9XDmpae2Qv"
	ts := "1513710493"

	chdr := hawk.CanonicalHeader {
		Kid:			kid,
		Key:            converted,
		Mac:            "",
		Header:         header,
		Time:           ts,
		Nonce:          nonce,
		Meth_or_Status: meth,
		Uri:            uri,
		Host:           host,
		Port:           port,
		PayloadHash:    "",
		Ext:            "",
	}

	expected_canonical_header := 	"hawk.1.header\n" +
		"1513710493\n" +
		"TO74yl3UnIQhoZAhMAKf9XDmpae2Qv\n" +
		"GET\n" +
		"/v3.0/apps/summary\n" +
		"api.dev.kaiostech.com\n" +
		"8090\n\n\n"

	actual_canonical_header := chdr.ToString()
	if expected_canonical_header != actual_canonical_header {
		t.Errorf("Canonical headers do not match, actual\n %s, expected\n %s", actual_canonical_header, expected_canonical_header)
	}

	h := hmac.New(sha256.New, expected_bytes)
	h.Write([]byte(actual_canonical_header))

	mac := base64.StdEncoding.EncodeToString(h.Sum(nil))
	expected := "rRtRT+1XxwFWDdusZc8JKUJe/cOyZ73ATVssBx3KNEo="

	if expected != mac {
		t.Errorf("Mac %s isn't expected value %s", mac, expected)
	}
}

func Test_Mac_Verify(t *testing.T) {

	key := "PRLSGjwba5OyIFqqH7e5XM3kAdTnXO4crbrAH628CuA="

	mac_key_bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Errorf("failed to convert test to base64 %v", err)
	}

	expected_bytes := []byte  {
		61,18,210,26,60,27,107,147,178,32,90,170,31,183,185,92,205,228,1,212,231,92,238,28,173,186,192,31,173,188,10,224}

	if len(mac_key_bytes) != len(expected_bytes) {
		t.Errorf("array lengths don't match converted %d vs expected %d", len(mac_key_bytes), len(expected_bytes))
	}

	for idx := 0; idx < len(expected_bytes); idx++ {

		if mac_key_bytes[idx] != expected_bytes[idx] {
			t.Errorf("Array elemenet %d doesn't match, converted %d vs expected %d", idx, mac_key_bytes[idx], expected_bytes[idx])
		}
	}

	// From js:
	// Hawk id="rjniBDZTk5XPQ4nRRURxLNlCGEw=", ts="1513710493", nonce="TO74yl3UnIQhoZAhMAKf9XDmpae2Qv",
	// mac="7gpLEYMy6cSx0SBgDC4G2CQOvbPOWY7gLUg/u8mkGzs="
	meth_or_code := "GET"
	http_host := "api.dev.kaiostech.com:8090"
	end_point_url, err := url.Parse("http://api.dev.kaiostech.com:8090/v3.0/apps/summary")
	kid := "rjniBDZTk5XPQ4nRRURxLNlCGEw"
	nonce_value := "TO74yl3UnIQhoZAhMAKf9XDmpae2Qv"
	ts := "1513710493"
	req_id := "1234"

	ch, err := hawk.Canonical_Header(meth_or_code, http_host, end_point_url, kid, mac_key_bytes,
		nonce_value, ts, "", "", false, req_id)
	if err != nil {
		t.Errorf("Failed to create canonical header, err %v", err)
	}

	expected_canonical_header := 	"hawk.1.header\n" +
									"1513710493\n" +
									"TO74yl3UnIQhoZAhMAKf9XDmpae2Qv\n" +
									"GET\n" +
									"/v3.0/apps/summary\n" +
									"api.dev.kaiostech.com\n" +
									"8090\n\n\n"

	actual_canonical_header := ch.ToString()
	if expected_canonical_header != actual_canonical_header {
		t.Errorf("Canonical headers do not match, actual\n %s, expected\n %s", actual_canonical_header, expected_canonical_header)
	}

	actual_mac := ch.Sign(req_id)

	expected := "rRtRT+1XxwFWDdusZc8JKUJe/cOyZ73ATVssBx3KNEo="

	if expected != actual_mac {
		t.Errorf("Mac %s isn't expected value %s", actual_mac, expected)
	}
}

func Test_Mac_Verify2(t *testing.T) {

	key := "t9aH/iFrQeUsoZJucauTHq1xBvdC/e4/5jH1PDBQuZc="

	mac_key_bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Errorf("failed to convert test to base64 %v", err)
	}

	expected_bytes := []byte  {
		183,214,135,254,33,107,65,229,44,161,146,110,113,171,147,30,173,113,6,247,66,253,238,63,230,49,245,60,48,80,185,151}

	if len(mac_key_bytes) != len(expected_bytes) {
		t.Errorf("array lengths don't match converted %d vs expected %d", len(mac_key_bytes), len(expected_bytes))
	}

	for idx := 0; idx < len(expected_bytes); idx++ {

		if mac_key_bytes[idx] != expected_bytes[idx] {
			t.Errorf("Array elemenet %d doesn't match, converted %d vs expected %d", idx, mac_key_bytes[idx], expected_bytes[idx])
		}
	}

	// From js:
	// Hawk id="TvlmUAvN8wfxFrV7entspPJCDU8=", ts="1513714840", nonce="QmoDeMNygB1sB7YYNazspmdRyrXgwX",
	// mac="MYI8LKo6WSzn+KJ3tFZTtGsCsiCNky9VzU/v06Q1QYc="
	meth_or_code := "GET"
	http_host := "api.dev.kaiostech.com:8090"
	end_point_url, err := url.Parse("http://api.dev.kaiostech.com:8090/v3.0/apps/summary")
	kid := "TvlmUAvN8wfxFrV7entspPJCDU8="
	nonce_value := "QmoDeMNygB1sB7YYNazspmdRyrXgwX"
	ts := "1513714840"
	req_id := "1234"

	ch, err := hawk.Canonical_Header(meth_or_code, http_host, end_point_url, kid, mac_key_bytes,
		nonce_value, ts, "", "", false, req_id)
	if err != nil {
		t.Errorf("Failed to create canonical header, err %v", err)
	}

	expected_canonical_header := 	"hawk.1.header\n" +
		"1513714840\n" +
		"QmoDeMNygB1sB7YYNazspmdRyrXgwX\n" +
		"GET\n" +
		"/v3.0/apps/summary\n" +
		"api.dev.kaiostech.com\n" +
		"8090\n\n\n"

	actual_canonical_header := ch.ToString()
	if expected_canonical_header != actual_canonical_header {
		t.Errorf("Canonical headers do not match, actual\n %s, expected\n %s", actual_canonical_header, expected_canonical_header)
	}

	actual_mac := ch.Sign(req_id)

	t.Logf("Block size %s", ch.Get_Info())

	expected := "sav3w+Bn1RWjU+N9BYuQMrQc/Qy9WPrmWO695g0cx/A="

	if expected != actual_mac {
		t.Errorf("Mac %s isn't expected value %s", actual_mac, expected)
	}
}
