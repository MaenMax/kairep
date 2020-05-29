package pwreset

import (
	"time"
	"empowerthings.com/cumulis/db/redisdb"
)

var (
	rdb *redisdb.RedisDB
)

func Init() error {
	rdb=redisdb.New("pwreset")
	// Now testing whether connection is ok.
	err := rdb.Ping()
	return err
}

func Set(key string, data string, expire time.Duration, req_id string) error {
	return rdb.Set(key,data,expire,req_id)
}

func Del(key string, req_id string) error {
	return rdb.Del(key,req_id)
}

func Get(key string, req_id string) (data string, err error) {
	return rdb.Get(key,req_id)
}


func SetB(key string, data []byte, expire time.Duration, req_id string) error {
	return rdb.SetB(key,data,expire,req_id)
}

func GetB(key string, req_id string) (data []byte, err error) {
	return rdb.GetB(key,req_id)
}
