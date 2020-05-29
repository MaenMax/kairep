package context

import (
	"net/http"

	"github.com/gorilla/context"
)

var empty_string string
var empty_buffer []byte

type key int

const KEY_AUTH_TYPE key = 1
const KEY_CLAIMS key = 2
const KEY_CANONIZED_HEADER key = 3
const KEY_PAYLOAD key = 4
const KEY_PAYLOAD_LEN key = 5
const KEY_STATUS_CODE key = 6
const KEY_RESPONSE_LENGTH key = 7
const KEY_PAGINATION key = 8
const KEY_REQ_ID key = 9
const KEY_SOURCE_IP key = 10

//SetPagination stores the pagination Cursor information
func SetPagination(r *http.Request, cursor []byte) {
	context.Set(r, KEY_PAGINATION, cursor)
}

func GetPagination(r *http.Request) []byte {
	if rv := context.Get(r, KEY_PAGINATION); rv != nil {
		return rv.([]byte)
	}
	return empty_buffer
}

func SetAuthType(r *http.Request, auth string) {
	context.Set(r, KEY_AUTH_TYPE, auth)
}

func GetAuthType(r *http.Request) string {
	if rv := context.Get(r, KEY_AUTH_TYPE); rv != nil {
		return rv.(string)
	}
	return empty_string
}

func GetClaims(r *http.Request) string {
	if rv := context.Get(r, KEY_CLAIMS); rv != nil {
		return rv.(string)
	}
	return empty_string
}

func SetClaims(r *http.Request, token string) {
	context.Set(r, KEY_CLAIMS, token)
}

func GetCanonizedHeader(r *http.Request) string {
	if rv := context.Get(r, KEY_CANONIZED_HEADER); rv != nil {
		return rv.(string)
	}
	return empty_string
}

func SetCanonizedHeader(r *http.Request, ch string) {
	context.Set(r, KEY_CANONIZED_HEADER, ch)
}

func GetPayload(r *http.Request) []byte {
	if rv := context.Get(r, KEY_PAYLOAD); rv != nil {
		return rv.([]byte)
	}
	return empty_buffer
}

func SetPayload(r *http.Request, payload []byte) {
	context.Set(r, KEY_PAYLOAD, payload)
}

func GetPayloadLen(r *http.Request) int {
	if rv := context.Get(r, KEY_PAYLOAD_LEN); rv != nil {
		return rv.(int)
	}
	return 0
}

func SetPayloadLen(r *http.Request, payload_len int) {
	context.Set(r, KEY_PAYLOAD_LEN, payload_len)
}

func GetStatusCode(r *http.Request) int {
	if rv := context.Get(r, KEY_STATUS_CODE); rv != nil {
		return rv.(int)
	}
	return 0
}

func SetStatusCode(r *http.Request, status_code int) {
	context.Set(r, KEY_STATUS_CODE, status_code)
}

func GetResponseLength(r *http.Request) int {
	if rv := context.Get(r, KEY_RESPONSE_LENGTH); rv != nil {
		return rv.(int)
	}
	return 0
}

func SetResponseLength(r *http.Request, response_length int) {
	context.Set(r, KEY_RESPONSE_LENGTH, response_length)
}

func SetReqId(r *http.Request, req_id string) {
	context.Set(r, KEY_REQ_ID, req_id)
}

func GetReqId(r *http.Request) string {
	if req_id := context.Get(r, KEY_REQ_ID); req_id != nil {
		return req_id.(string)
	}
	return empty_string
}

func SetSourceIP(r *http.Request, clientIP string) {
	context.Set(r, KEY_SOURCE_IP, clientIP)
}

func GetSourceIP(r *http.Request) string {
	if clientIP := context.Get(r, KEY_SOURCE_IP); clientIP != nil {
		return clientIP.(string)
	}
	return empty_string
}

func ClearContext(r *http.Request) {
	context.Clear(r)
}
