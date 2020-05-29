package counterdb

import (
	"time"
	"strconv"

	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/routercore/db/redis/redisdb"
	"github.com/gocql/gocql"
	//redis "gopkg.in/redis.v5"
	
	redis "github.com/go-redis/redis"
)

var (
	rdb          *redisdb.RedisDB
	count_script *redis.Script
	max_msg      int
)

func Init() error {
	rdb = redisdb.New()

	// Now testing whether connection is ok.
	err := rdb.Ping()
	if err == nil {
		count_script = Register_Count_Script()
	}
	return err
}

// Register new lua script immediately after establishing a connection
func Register_Count_Script() *redis.Script {
	lua_script := ` redis.replicate_commands()
					if redis.call("EXISTS",KEYS[1]) == 0 then
					   local count = redis.call("SET",KEYS[1],1)
					   return 1
					else
					   local count = redis.call("GET",KEYS[1])
					   if tonumber(count) < tonumber(ARGV[1]) then
					        local count = redis.call("INCR",KEYS[1])
					        return 1
					   end
					return 0
					end`
	return rdb.RegisterScript(lua_script)
}

func RPushB(key string, data []byte) error {
	return rdb.RPush(key, string(data))
}

func RPush(key string, data string) error {
	return rdb.RPush(key, data)
}

func LPopB(key string) (data []byte, err error) {
	return rdb.LPopB(key)
}

func LPop(key string) (data string, err error) {
	return rdb.LPop(key)
}

func Exists(key string) (ok bool, err error) {
	return rdb.Exists(key)
}

func Set(key string, data string, expire time.Duration) error {
	return rdb.Set(key, data, expire)
}

func Del(key string) error {
	return rdb.Del(key)
}

func Get(key string) (data string, err error) {
	return rdb.Get(key)
}

func Inc(key string) (data int64, err error) {
	return rdb.Inc(key)
}

// Validate counter limit with lua script
func Script_Validate_Limit(key string) (result int64, err error) {
	return rdb.Script(count_script, key, max_msg)
}

func Reached_Msg_Max(uaid string, session *gocql.Session, debug bool) (bool, error) {

	if debug {

		l4g.Info("REP node is checking max msg limit for uaid: %v", uaid)

	}

	// GET counter value from Redis
	result, err := Get(uaid)
	if err != nil {
		return false, err
	}
	
	result_int := 0 
	if result == "" {
		if debug {
			l4g.Info("New device goes offline uaid: %s", uaid)
		}
		result_int = 0
	}

	result_int, err = strconv.Atoi(result)
	if err != nil {
		result_int = 0 
	}

	if result_int >= max_msg {
		if debug {
			l4g.Error("maximum counter value reached for uaid: %s", uaid)
		}
		return true, nil
	}

	status_int, err := Script_Validate_Limit(uaid)

	if err != nil {
		return false, err
	}

	if status_int == 1 {
		if debug {
			l4g.Info("Counter Incremented for uaid: %s", uaid)
		}

		return false, nil
	}

	return true, nil

}

func SetMaxMsg(maximum_msg int) {

	max_msg = maximum_msg

}
