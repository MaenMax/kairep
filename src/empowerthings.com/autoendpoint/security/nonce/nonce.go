package nonce

import (
	"time"
	"empowerthings.com/cumulis/db/redisdb"
//	l4g "code.google.com/p/log4go"
)

var (
	rdb *redisdb.RedisDB
)

func Init() error {
	rdb=redisdb.New("nonce")
	// Now testing whether connection is ok.
	err := rdb.Ping()
	return err
}

func Set(nonce_id string, nonce_data string, expire time.Duration, req_id string) error {
	return rdb.Set(nonce_id,nonce_data,expire,req_id)
}

func Del(nonce_id string, req_id string) error {
	return rdb.Del(nonce_id,req_id)
}

func Get(nonce_id string, req_id string) (nonce_data string, err error) {
	return rdb.Get(nonce_id,req_id)
}

