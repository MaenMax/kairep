// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webpush

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/url"
	"time"
)

var (
	// encoded {typ:JWT, alg: ES256}
	vapidPrefix = []byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.")
	dot         = []byte(".")
)

//TODO: if vapidPrefix is missing, decode the prefix into jwtPrefix
// and check if it has the correct values.
type jwtPrefix struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type jwtBody struct {
	Aud string `json:"aud"`
	Sub string `json:"sub,omitempty"`
	Exp int64  `json:"exp"`
}

// Vapid represents a sender identity.
type Vapid struct {
	// The EC256 public key, base64 urlsafe, uncompressed. This value should be
	// used in 'subscribe' requests and is included as p256ecdsa in
	// the Crypto-Key header.
	PublicKey string

	// The private key used to sign tokens
	pkey *ecdsa.PrivateKey

	// Sub should be an email or URL, for identification
	Sub string
}

// TODO: cache for tokens

// Token creates a token with the specified endpoint, using configured Sub id
// and a default expiration (1h).
func (vapid *Vapid) Token(aud string) (res string) {
	url, _ := url.Parse(aud)
	host := url.Host
	jwt := jwtBody{Aud: "https://" + host}
	if vapid.Sub != "" {
		jwt.Sub = vapid.Sub
	}
	jwt.Exp = int64(time.Now().Unix() + 3600)
	t, _ := json.Marshal(jwt)
	enc := base64.RawURLEncoding

	t64 := make([]byte, enc.EncodedLen(len(t)))
	enc.Encode(t64, t)

	token := make([]byte, len(t)+len(vapidPrefix)+100)
	token = append(token[:0], vapidPrefix...)
	token = append(token, t64...)

	hasher := crypto.SHA256.New()
	hasher.Write(token)

	if r, s, err := ecdsa.Sign(rand.Reader, vapid.pkey, hasher.Sum(nil)); err == nil {
		// Vapid key is 32 bytes
		keyBytes := 32
		sig := make([]byte, 2*keyBytes)

		rBytes := r.Bytes()
		rBytesPadded := make([]byte, keyBytes)
		copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

		sBytes := s.Bytes()
		sBytesPadded := make([]byte, keyBytes)
		copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

		sig = append(sig[:0], rBytesPadded...)
		sig = append(sig, sBytesPadded...)

		sigB64 := make([]byte, enc.EncodedLen(len(sig)))
		enc.Encode(sigB64, sig)

		token = append(token, dot...)
		token = append(token, sigB64...)
	}
	res = string(token)
	return
}

// NewVapid constructs a new Vapid generator from EC256 public and private keys,
// in uncompressed format.
func NewVapid(publicKey, privateKey string) (v *Vapid) {
	publicUncomp, _ := base64.RawURLEncoding.DecodeString(publicKey)
	privateUncomp, _ := base64.RawURLEncoding.DecodeString(privateKey)

	x, y := elliptic.Unmarshal(curve, publicUncomp)
	d := new(big.Int).SetBytes(privateUncomp)
	pubkey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	pkey := ecdsa.PrivateKey{PublicKey: pubkey, D: d}

	v = &Vapid{
		PublicKey: publicKey,
		pkey:      &pkey}

	return
}
