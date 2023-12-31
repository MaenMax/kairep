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
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	gcmURL     = "https://android.googleapis.com/gcm/send"
	tempGcmURL = "https://gcm-http.googleapis.com/gcm"
)

// NewPushRequest creates a valid Web Push HTTP request for sending a message
// to a subscriber. If the push service requires an authentication header
// (notably Google Cloud Messaging, used by Chrome) then you can add that as the
// token parameter.
// Deprecated - token auth is not part of the spec.
func NewPushRequest(sub *Subscription, message string, token string) (*http.Request, error) {
	// If the endpoint is GCM then we temporarily need to rewrite it, as not all
	// GCM servers support the Web Push protocol. This should go away in the
	// future.
	endpoint := strings.Replace(sub.Endpoint, gcmURL, tempGcmURL, 1)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// TODO: Make the TTL variable
	req.Header.Add("TTL", "0")

	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf(`key=%s`, token))
	}

	// If there is no payload then we don't actually need encryption
	if message == "" {
		return req, nil
	}

	payload, err := Encrypt(sub, message)
	if err != nil {
		return nil, err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(payload.Ciphertext))
	req.ContentLength = int64(len(payload.Ciphertext))
	req.Header.Add("Encryption", headerField("salt", payload.Salt))
	req.Header.Add("Crypto-Key", headerField("dh", payload.ServerPublicKey))
	req.Header.Add("Content-Encoding", "aesgcm")

	return req, nil
}

// NewVapidRequest creates a valid Web Push HTTP request for sending a message
// to a subscriber, using Vapid authentication. You can add more headers to
// configure collapsing, TTL.
func NewRequest(to *Subscription, message string, ttlSec int, vapid *Vapid, as_private_key []byte, as_public_key []byte, as_public_key_str string, receipt bool, urgency string, topic bool) (*http.Request, error) {
	// If the endpoint is GCM then we temporarily need to rewrite it, as not all
	// GCM servers support the Web Push protocol. This should go away in the
	// future.
	req, err := http.NewRequest("POST", to.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("ttl", strconv.Itoa(ttlSec))
	if receipt == true {
		req.Header.Add("Prefer", "respond-async")
	}
	if topic == true {
		req.Header.Add("Topic", "clock")
	}
	req.Header.Add("Urgency", urgency)

	if vapid != nil {
		tok := vapid.Token(to.Endpoint)
		//~ req.Header.Add("Authorization", fmt.Sprintf(`vapid t=%s k=%s`, tok, as_public_key_str))
		req.Header.Add("Authorization", fmt.Sprintf(`WebPush %s`, tok))
	}

	// If there is no payload then we don't actually need encryption
	if message == "" {
		return req, nil
	}

	//~ payload, err := Encrypt(to, message)
        plaintext := []byte(message)
	payload, err := EncryptWithTempKey(to, plaintext, as_private_key, as_public_key)

	if err != nil {
		return nil, err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(payload.Ciphertext))
	req.ContentLength = int64(len(payload.Ciphertext))
	req.Header.Add("Encryption", headerField("salt", payload.Salt))
	if vapid != nil {
		req.Header.Add("Crypto-Key",
			headerField("dh", payload.ServerPublicKey)+",p256ecdsa="+
				vapid.PublicKey)
	} else {
		req.Header.Add("Crypto-Key",
			headerField("dh", payload.ServerPublicKey))
	}
	req.Header.Add("Content-Encoding", "aesgcm")

	return req, nil
}

// Send a message using the Web Push protocol to the recipient identified by the
// given subscription object. If the client is nil then the default HTTP client
// will be used. If the push service requires an authentication header (notably
// Google Cloud Messaging, used by Chrome) then you can add that as the token
// parameter.
func Send(client *http.Client, sub *Subscription, message, token string) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := NewPushRequest(sub, message, token)
	if err != nil {
		return nil, err
	}

	// Default TTL
	req.Header.Add("ttl", "0")
	return client.Do(req)
}

// A helper for creating the value part of the HTTP encryption headers
func headerField(headerType string, value []byte) string {
	return fmt.Sprintf(`%s=%s`, headerType, base64.RawURLEncoding.EncodeToString(value))
}
