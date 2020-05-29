package credkey

import (
	"io"
	"bytes"
	"time"
	"crypto/sha256"
	"encoding/base64"
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/hkdf"

	"empowerthings.com/cumulis/utils/uuid"
	"empowerthings.com/cumulis/utils"
	"empowerthings.com/cumulis/security/rand"

	"empowerthings.com/cumulis/db/redisdb"
//	l4g "code.google.com/p/log4go"
)

const (
	NAMESPACE="credkey.empowerthings.com/v1/"
	SALT="'24;6344 3^%$;6743^(687834254"
)

var (
	_xor *rand.Xoroshiro128Plus
	_mix64 *rand.SplitMix64
	rdb *redisdb.RedisDB
)

func init() {	
	seed:=time.Now().UnixNano()
	_mix64=rand.NewSplitMix64(seed)
	_mix64.Next()
	_xor=rand.NewXoroshiro128Plus(uint64(_mix64.Next()),uint64(_mix64.Next()))

}

func Init() error {
	rdb=redisdb.New("credkey")
	// Now testing whether connection is ok.
	err := rdb.Ping()
	return err
}

func SetB(cred_id string, cred_value []byte, expire time.Duration, req_id string) error {
	return rdb.SetB(cred_id,cred_value,expire,req_id)
}

func Del(cred_id string, req_id string) error {
	return rdb.Del(cred_id,req_id)
}

func Get(cred_id string, req_id string) (cred_value []byte, err error) {
	return rdb.GetB(cred_id,req_id)
}

// Returns a base64 encoded key id but a byte array as key.
func Gen(token string) (key_id string, mac_key []byte) {
	var l int=0
	var tl int = len(token)+len(SALT)+8

	// Generating Key ID
	key_id_seed:=make([]byte,tl,tl)
	
	copy(key_id_seed[l:],utils.UnsafeCastInt64ToBytes(int64(_xor.Next())))
	l+=8
	copy(key_id_seed[l:],token)
	l+=len(token)
	copy(key_id_seed[l:],SALT)
	l+=len(SALT)

	hash:=sha1.New()
	io.Copy(hash,bytes.NewReader(key_id_seed))
	key_id_hash:=hash.Sum(nil)
	
	key_id=base64.StdEncoding.EncodeToString(key_id_hash)

	// Generating key
	quickStretched:=_pbkdf2([]byte(token),_kwe([]byte("quickStretchcredkey"),[]byte(token)),10,32)
	credkey:=_hkdf(quickStretched, _kw([]byte("credkey")), []byte(SALT),32)

	//mac_key=base64.StdEncoding.EncodeToString(credkey)
	mac_key=credkey

	return key_id,mac_key
}

func _kw(mac_id []byte) []byte {
	var buffer bytes.Buffer

	uuid.NewUuid()
	buffer.WriteString(NAMESPACE)
	buffer.Write(mac_id)
	buffer.WriteString(":")
	buffer.WriteString(uuid.NewUuid())
	return buffer.Bytes()
}

func _kwe(mac_id []byte, token []byte) []byte {
	var buffer bytes.Buffer

	uuid.NewUuid()
	buffer.WriteString(NAMESPACE)
	buffer.Write(mac_id)
	buffer.WriteString(":")
	buffer.Write(token)
	return buffer.Bytes()
}

func _hkdf(ikm []byte, info []byte, salt []byte,len int) []byte {
	buf:=make([]byte,len)
	
	hkdf_io:=hkdf.New(sha256.New, ikm, salt, info)

	n,err:=io.ReadFull(hkdf_io,buf)
	
	if err!=nil {
		//panic("Failed to read hkdf ")
		return nil
	}

	if n!=len {
		//panic("Failed to generate required length hkdf output.")
		return nil
	}

	return buf

}

func _pbkdf2(p []byte, s []byte, c int, len int) []byte {
	result:= pbkdf2.Key(p, s, c, len, sha256.New)
	return result;
}

