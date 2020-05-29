package hawk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"empowerthings.com/cumulis/security/nonce"
	redis "gopkg.in/redis.v5"

	l4g "code.google.com/p/log4go"
	"empowerthings.com/cumulis/config"
)

var (
	_conf     *config.Config
	_rediscli *redis.ClusterClient
)

type CanonicalHeader struct {
	Kid            string
	Key            []byte
	Mac            string
	Header         string
	Time           string
	Nonce          string
	Meth_or_Status string
	Uri            string
	Host           string
	Port           string
	PayloadHash    string
	Ext            string
}

type CanonicalPayload struct {
	Header      string
	ContentType string
	Data        []byte
}

func Init(conf *config.Config) error {
	var opt *redis.ClusterOptions = &redis.ClusterOptions{}
	_conf = conf

	//addrs = make([]string,len(conf.RedisService.Servers))


	opt.Addrs = strings.Split(conf.Redis.Host,",")
	opt.MaxRedirects = conf.Redis.Max_Redirects

	_rediscli = redis.NewClusterClient(opt)

	// Now testing whether connection is ok.
	_, err := _rediscli.Ping().Result()

	if err != nil {
		return err
	}

	return nil
}

func (hcontext *HawkContext) Sign(req *http.Request, nonce string, now time.Time, ext string, payload_hash string, req_id string) error {
	var method string

	if len(req.Method) > 0 {
		method = strings.ToUpper(req.Method)
	} else {
		method = "GET"
	}

	is_tls_active := (req.TLS != nil)
	ch, err := Canonical_Header(method, req.Host, req.URL, hcontext.Kid, hcontext.Key, nonce, fmt.Sprintf("%v", now.Unix()), ext, payload_hash, is_tls_active, req_id)

	if err != nil {
		return nil
	}

	_ = ch.Sign(req_id)

	auth := ch.Generate()

	if config.GetConfig().Service.Debug {
		l4g.Debug("Req #%v: AuthHeader='%s'", req_id, auth)
	}
	req.Header.Set("Authorization", auth)

	return nil
}

func Verify(meth_or_code string, http_host string, url *url.URL, auth_head *map[string]string, kid string, hawk_key []byte, is_tls_active bool, check_nonce bool, req_id string) (ch *CanonicalHeader, err error) {
	var actual_mac string
	var hash string
	var ext string

	kid, ok := (*auth_head)["id"]

	if !ok {
		l4g.Error("Req #%v: Invalid Hawk Authorization Header: 'id' is missing.", req_id)
		return nil, errors.New("Invalid Hawk Authorization header.")
	}

	ts, ok := (*auth_head)["ts"]

	if !ok {
		l4g.Error("Req #%v: Invalid Hawk Authorization Header: 'ts' is missing.", req_id)
		return nil, errors.New("Invalid Hawk Authorization header.")
	}

	// Checking time stamp validity.
	now := time.Now().UTC().Unix()
	ts_v, err := strconv.ParseInt(ts, 10, 64)

	if err != nil {
		l4g.Error("Req #%v: Invalid Hawk Authorization Header: 'ts' has invalid format.", req_id)
		return nil, errors.New("Invalid Hawk Authorization header (Invalid 'ts' format).")
	}

	// We accept a time skew of +/- 60 seconds.
	// This means the request should come between
	// [now-60, now+60]
	// If not, then we reject.
	// This is why we will temporary store nonce for 2 minutes only.
	if (ts_v-60) > now || (ts_v+60) < now {
		l4g.Error("Req #%v: Bad Time Stamp in Hawk Authorization Header.", req_id)
		return nil, errors.New("Bad Time Stamp in Hawk Authorization Header.")
	}

	nonce_value, ok := (*auth_head)["nonce"]

	if !ok {
		l4g.Error("Req %v: Invalid Hawk Authorization Header: 'nonce' is missing.", req_id)
		return nil, errors.New("Invalid Hawk Authorization header.")
	}

	if check_nonce  {
		// Checking whether the nonce has been already seen!
		cached_nonce,err:=nonce.Get(nonce_value,req_id)
		
		// If already seen, then no error should be reported from nonce DB!

		if err==nil && len(cached_nonce) > 0 {
			l4g.Error("Req #%v: Invalid Hawk Authorization Header: Already used Nonce '%s'.",req_id, nonce_value)
			return nil, errors.New("Invalid Hawk Authorization header.")
		}

		// Making sure we won't see that nonce again during the coming 2 min.
		// We only need to store nonce for 2 minutes because we check the time.
		// If hackers want to use older nonce, then it is going to fail due to
		// time limit.
		err = nonce.Set(nonce_value, nonce_value, time.Second*120, req_id)
	}

	mac, ok := (*auth_head)["mac"]

	if !ok {
		l4g.Error("Req #%v: Invalid Hawk Authorization Header: 'mac' is missing.", req_id)
		return nil, errors.New("Invalid Hawk Authorization header.")
	}

	// Optional field
	ext, ok = (*auth_head)["ext"]

	// Optional field
	hash, ok = (*auth_head)["hash"]

	ch, err = Canonical_Header(meth_or_code, http_host, url, kid, hawk_key, nonce_value, ts, ext, hash, is_tls_active, req_id)

	if err != nil {
		return nil, err
	}

	actual_mac = ch.Sign(req_id)

	if strings.Compare(actual_mac, mac) != 0 {
		return nil, errors.New(fmt.Sprintf("Signature don't match '%s'!='%s'.", actual_mac, mac))
	}

	return ch, nil
}

func HashPayload(content_type string, data []byte, req_id string) string {

	cp := canonical_payload(content_type, data)

	return cp.Hash(req_id)
}

func canonical_payload(content_type string, data []byte) *CanonicalPayload {

	cp := &CanonicalPayload{}

	cp.Header = "hawk.1.payload"
	cp.ContentType = parse_content_type(content_type)
	cp.Data = data

	return cp
}

func parse_content_type(content_type string) string {

	if len(content_type) == 0 {
		return content_type
	}

	result := strings.Split(content_type, ";")

	return strings.ToLower(strings.TrimSpace(result[0]))
}

func Canonical_Header(meth_or_code string, http_host string, url *url.URL, kid string, key []byte, nonce string, ts string, ext_data string, payload_hash string, is_tls_active bool, req_id string) (ch *CanonicalHeader, err error) {

	ch = &CanonicalHeader{}

	ch.Kid = kid
	ch.Key = key

	uri, err := url.Parse(url.String())

	if err != nil {
		err2 := l4g.Error("Req #%v: Failed to parse URL: '%s'", req_id, err)
		if config.GetConfig().Service.Debug {
			l4g.Debug("Req #%v: URL leading to error: '%s'.", req_id, url.String())
		}
		return nil, err2
	}

	host, port := ExtractHostPort(http_host, uri.Scheme, is_tls_active)

	if len(host) == 0 {
		err := l4g.Error("Req #%v: Missing or malformed 'Host' field into the HTTP header (value='%s').", req_id, http_host)
		return nil, err
	}

	// Header name
	ch.Header = "hawk.1.header"

	// Timestamp
	ch.Time = ts

	// Nonce
	ch.Nonce = nonce

	// HTTP Method or Status
	ch.Meth_or_Status = meth_or_code

	// Canonical URI + Eventual Query parameters if available.
	if url != nil && len(url.RawQuery) > 0 {
		ch.Uri = fmt.Sprintf("%s?%s", url.Path, url.RawQuery)
	} else {
		ch.Uri = fmt.Sprintf("%s", url.Path)
	}

	// Host
	ch.Host = fmt.Sprintf("%s", host)

	// Host port if available
	if len(port) > 0 {
		ch.Port = fmt.Sprintf("%s", port)
	}

	if len(payload_hash) > 0 {
		ch.PayloadHash = payload_hash
	}

	// App External Data
	if len(ext_data) > 0 {
		ext_data = strings.Replace(ext_data, "\\", "\\\\", -1)
		ext_data = strings.Replace(ext_data, "\n", "\\n", -1)
		ch.Ext = fmt.Sprintf("%s", ext_data)
	}

	return ch, nil
}

func (ch *CanonicalHeader) ToString() string {
	var canonical_header_b bytes.Buffer
	canonical_header_b.WriteString(ch.Header + "\n")
	canonical_header_b.WriteString(ch.Time + "\n")
	canonical_header_b.WriteString(ch.Nonce + "\n")
	canonical_header_b.WriteString(ch.Meth_or_Status + "\n")
	canonical_header_b.WriteString(ch.Uri + "\n")
	canonical_header_b.WriteString(ch.Host + "\n")
	canonical_header_b.WriteString(ch.Port + "\n")
	canonical_header_b.WriteString(ch.PayloadHash + "\n")
	canonical_header_b.WriteString(ch.Ext + "\n")
	return canonical_header_b.String()
}

func (cp *CanonicalPayload) Bytes() []byte {
	var canonical_payload_b bytes.Buffer
	canonical_payload_b.WriteString(cp.Header + "\n")
	canonical_payload_b.WriteString(cp.ContentType + "\n")
	canonical_payload_b.Write(cp.Data)
	canonical_payload_b.WriteString("\n")
	return canonical_payload_b.Bytes()
}

// Return the Base64 of the HMAC-SHA256 signature with the given secret key
// of the provided canonical header.
func (ch *CanonicalHeader) Sign(req_id string) string {
	h := hmac.New(sha256.New, ch.Key)
	canonical_header_str := ch.ToString()
	h.Write([]byte(canonical_header_str))
	ch.Mac = base64.StdEncoding.EncodeToString(h.Sum(nil))

	if config.GetConfig().Service.Debug {
		l4g.Debug("Req #%v: key=%v", req_id, ch.Key)
		l4g.Debug("Req #%v:\ncanonical_header=--------------\n%s\n-------------", req_id, canonical_header_str)
		l4g.Debug("Req #%v: HMAC=%s", req_id, ch.Mac)
	}
	return ch.Mac
}

func (ch *CanonicalHeader) Get_Info() string {
	h := hmac.New(sha256.New, ch.Key)
	return fmt.Sprintf("block size=%d, size=%d", h.BlockSize(), h.Size());
}

func (ch *CanonicalHeader) Generate() string {
	var auth bytes.Buffer

	auth.WriteString("Hawk")
	auth.WriteString(fmt.Sprintf(" id=\"%s\"", ch.Kid))
	auth.WriteString(fmt.Sprintf(", ts=%v", ch.Time))
	auth.WriteString(fmt.Sprintf(", nonce=\"%s\"", ch.Nonce))

	if len(ch.Ext) > 0 {
		auth.WriteString(fmt.Sprintf(", ext=\"%s\"", ch.Ext))
	}

	if len(ch.PayloadHash) > 0 {
		auth.WriteString(fmt.Sprintf(", hash=\"%s\"", ch.PayloadHash))
	}

	auth.WriteString(fmt.Sprintf(", mac=\"%s\"", ch.Mac))

	return auth.String()
}

// Return the Base64 of the SHA256 Hash of the provided canonical payload.
func (cp *CanonicalPayload) Hash(req_id string) string {
	canonical_payload := cp.Bytes()

	hash_b := sha256.Sum256(canonical_payload)

	hash_str := base64.StdEncoding.EncodeToString(hash_b[:])

	if config.GetConfig().Service.Debug {
		l4g.Debug("Req #%v: \ncanonical_payload=--------------%s-------------", req_id, canonical_payload)
		l4g.Debug("Req #%v: HASH=%s", req_id, hash_str)
	}
	return hash_str
}

func ExtractHostPort(hostport string, scheme string, is_tls_active bool) (host string, port string) {

	host, port, _ = net.SplitHostPort(hostport)

	if port == "" {
		switch scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"

		default:
			// RS - 2016/12/01 Bug #46  http://dev.empowerthings.com/tct/empowerthings/issues/46
			// When Firefox is making request on default HTTP or HTTPS port, then we will have empty
			// scheme and port.
			// Then to make the difference between port 80 or 443, we need to use is_tls_active
			// parameter.
			if is_tls_active {
				port = "443"
			} else {
				port = "80"
			}
		}
	}

	if host == "" {
		host = hostport
	}

	return host, port

}
