package jwttoken

import (
	"empowerthings.com/autoendpoint/utils/uuid"

	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

type Claims map[string]interface{}

func NewClaims() Claims {
	var tmp Claims

	tmp = make(Claims)
	jwtid := uuid.NewUuid()
	tmp.SetJWTID(jwtid)

	return tmp
}

func (c Claims) Set(key string, value interface{}) {
	c[key] = value
}

func (c Claims) Get(key string) interface{} {
	return c[key]
}

func (claims Claims) Has(key string) bool {
	_, ok := claims[key]
	return ok
}

func (c Claims) Del(key string) {
	delete(c, key)
}

func (c Claims) SetTime(key string, t time.Time) {
	c.Set(key, t.Unix())
}

func (c Claims) GetTime(key string) (time.Time, bool) {
	var zero = time.Time{}

	x := c.Get(key)
	if x == nil {
		return zero, false
	}
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		return time.Unix(v.Int(), 0), true
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return time.Unix(int64(v.Uint()), 0), true
	case reflect.Float64:
		return time.Unix(int64(v.Float()), 0), true
	default:
		return zero, false
	}
}

func (c Claims) GetInt(key string) (v int, ok bool) {

	x := c.Get(key)
	if x == nil {
		return 0, false
	}

	switch x.(type) {
	case int:
		return x.(int), true

	case uint32:
		return int(x.(uint32)), true

	case int64:
		return int(x.(int64)), true

	case uint64:
		return int(x.(uint64)), true

	case float32:
		return int(x.(float32)), true

	case float64:
		return int(x.(float64)), true

	default:
		return 0, false
	}

	return 0, false
}

func (c Claims) SetExpiration(t time.Time) {
	c.SetTime("exp", t)
}

func (c Claims) Expiration() (t time.Time, ok bool) {
	return c.GetTime("exp")
}

func (c Claims) SetIssuer(issuer string) {
	c.Set("iss", issuer)
}

func (c Claims) SetSubject(subject string) {
	c.Set("sub", subject)
}

func (c Claims) SetAudience(audience ...string) {
	if len(audience) == 1 {
		c.Set("aud", audience[0])
	} else {
		c.Set("aud", audience)
	}
}

func (c Claims) SetIssuedAt(issuedAt time.Time) {
	c.SetTime("iat", issuedAt)
}

func (c Claims) SetJWTID(uniqueID string) {
	c.Set("jti", uniqueID)
}

func (c Claims) JWTID() (jti string, ok bool) {
	v, ok := c.Get("jti").(string)
	return v, ok
}

func ParseClaims(claims_str string) (claims Claims, err error) {
	var tmp interface{}
	var tmp_map map[string]interface{}

	err = json.Unmarshal([]byte(claims_str), &tmp)

	if err != nil {
		return nil, err
	}

	tmp_map = tmp.(map[string]interface{})

	claims = Claims(tmp_map)

	return claims, nil

}

func ReadClaimsFromFile(claims_file string) (claims Claims, err error) {
	claims_content, err := ioutil.ReadFile(claims_file)

	if err != nil {
		return nil, err
	}

	claims, err = ParseClaims(string(claims_content))

	return claims, err
}

/**
  Returns the clear text JSON representation claims from the encoded token
  (i.e. as received from the HTTP client).

  Example: if claims_token contains:

  eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiIyMzVrcndscmp3ZWZhc2xyMzJxYyI
  sImRpZCI6IjVrMzQ1NjM0a2o0NjMzNDkwZXNsIiwiZXhwIjoxNDYzNTU0MTQwLCJpYXQiOjE0NjM
  1NTA1NDAsImlzcyI6InRlc3QuZW1wb3dlcnRoaW5ncy5jb20iLCJqdGkiOiJWcFlwRVZybGlYMHp
  EQ1phR1A1YUdYemRMZyIsInVpZCI6IjRyZTV3Zmxqc2tmbDhzOWRmazBhIiwidW5tIjoiUmFmZmk
  ifQ.AykKLrOXB1hO_102LbkSBPwGadZaKPR5AvrxsMA10rA1KV5b8z9zKGjkQ8HuA9VQ1h87noWL
  RuIHdOVPxFxqTvGcM7gKeSFrjhXHcNbceh3gFnjXVZHZC3VByM-HbWPl4l6sgmPQ7W0Ibi5elDCF
  HLSvKqeaJS3Ki_T-2PBAaViPehgd42A4LeOuVZbfY3eYqqWuVdBcAf19ofeO2FZwRiI7b8BxN9eH
  NbE35oC58XbLnZVDllLsNlb7ghG4NJOrWApcpRlcp7dZuhtXlJKLm6HnfozzlDtoUFNZfANnHd6M
  HK6JR-7rX5SfV4_WO_B7tJU4m_pQJBdG8MxxGP9b5g

  then this function should extract the substring found between the two '.'
  charaters:  eyJhaWQi....joiUmFmZmkifQ, decode the corresponding base64
  representation and return it as:


*/
func ExtractClaims(enc_token string) (claims_token string, err error) {
	var result []byte

	tok_parts := strings.Split(enc_token, ".")

	if len(tok_parts) != 3 {
		return claims_token, errors.New("Malformed token!")
	}

	result, err = base64.RawURLEncoding.DecodeString(tok_parts[1])

	if err != nil {
		return claims_token, err
	}

	claims_token = string(result)

	return claims_token, nil
}
