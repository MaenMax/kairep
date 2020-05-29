package ratelimit

import (
	"time"
	"empowerthings.com/cumulis/db/redisdb"
//	l4g "code.google.com/p/log4go"
)

var (
	rdb *redisdb.RedisDB
)

func Init() error {
	rdb=redisdb.New("ratelimit")
	// Now testing whether connection is ok.
	err := rdb.Ping()
	return err
}

func Set(ratelimit_id string, ratelimit string, expire time.Duration, req_id string) error {
	return rdb.Set(ratelimit_id,ratelimit,expire,req_id)
}

func Del(ratelimit_id string, req_id string) error {
	return rdb.Del(ratelimit_id,req_id)
}

func Get(ratelimit_id string, req_id string) (ratelimit string, err error) {
	return rdb.Get(ratelimit_id,req_id)
}

