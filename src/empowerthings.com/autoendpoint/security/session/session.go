package session

import (
	"time"
	"empowerthings.com/cumulis/db/redisdb"
//	l4g "code.google.com/p/log4go"
)

var (
	rdb *redisdb.RedisDB
)

func Init() error {
	rdb=redisdb.New("session")
	// Now testing whether connection is ok.
	err := rdb.Ping()
	return err
}

func Set(session_id string, session string, expire time.Duration, req_id string) error {
	return rdb.Set(session_id,session,expire, req_id)
}

func Del(session_id string, req_id string) error {
	return rdb.Del(session_id,req_id)
}

func Get(session_id string, req_id string) (session string, err error) {
	return rdb.Get(session_id,req_id)
}

