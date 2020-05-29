package token_blacklist

import (
	"time"
	"empowerthings.com/cumulis/db/redisdb"
//	l4g "code.google.com/p/log4go"
)

var (
	rdb *redisdb.RedisDB
)

func Init() error {
	rdb=redisdb.New("bl_token")
	// Now testing whether connection is ok.
	err := rdb.Ping()
	return err
}

func Set(jti string, expire time.Duration, req_id string) error {
	// Putting a dummy value "x"
	return rdb.Set(jti,jti,expire,req_id)
}

func Exists(jti string,req_id string) (ok bool, err error) {
	return rdb.Exists(jti,req_id)
}
