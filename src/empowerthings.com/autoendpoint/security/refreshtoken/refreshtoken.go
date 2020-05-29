package refreshtoken

import (
	"time"
	"empowerthings.com/cumulis/db/redisdb"
//	l4g "code.google.com/p/log4go"
)

var (
	rdb *redisdb.RedisDB
)

func Init() error {
	rdb=redisdb.New("refreshtoken")
	// Now testing whether connection is ok.
	err := rdb.Ping()
	return err
}

func Set(refresh_id string, key_id string, expire time.Duration, req_id string) error {
	return rdb.Set(refresh_id,key_id,expire,req_id)
}

func Del(refresh_id string, req_id string) error {
	return rdb.Del(refresh_id,req_id)
}

func Get(refresh_id string, req_id string) (key_id string, err error) {
	return rdb.Get(refresh_id,req_id)
}


func SetB(refresh_id string, key_id []byte, expire time.Duration, req_id string) error {
	return rdb.SetB(refresh_id,key_id,expire,req_id)
}
