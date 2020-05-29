package redisdb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/config"
//	redis "gopkg.in/redis.v5"
	redis "github.com/go-redis/redis"
)

var (
	_conf  *config.AutoEndpointConfig
	_addrs []string
	_pass  string
)

type ClusterNode struct {
	Id       string
	Addr     string
	IsMaster bool
	Master   *ClusterNode
}

type ClusterSlot struct {
	Start  int64
	End    int64
	Master *ClusterNode
}

func Init(conf *config.AutoEndpointConfig, password string) error {
	_conf = conf
	_addrs = strings.Split(conf.Redis.Host, ",")
	_pass = password
	return nil
}

type RedisDB struct {
	rediscli *redis.ClusterClient
	prefix   string
}

func New() *RedisDB {

	var tmp *RedisDB = &RedisDB{}
	var opt *redis.ClusterOptions = &redis.ClusterOptions{}

	opt.Addrs = _addrs
	opt.MaxRedirects = _conf.Redis.Max_Redirects
	opt.PoolSize = _conf.Redis.Connection_Pool_Size
	opt.Password = _pass
	opt.ReadOnly = true
	
	if _conf.Debug {
		l4g.Info(fmt.Sprintf("Connecting to Redis cluster"))
		l4g.Info(fmt.Sprintf("Redis cluster Option: %v", opt))
	}

	tmp.rediscli = redis.NewClusterClient(opt)
	tmp.prefix = _conf.Redis.Table_Prefix
	return tmp

}

func NewConnection(prefix string, host string, port int) *RedisDB {
	var tmp *RedisDB = &RedisDB{}
	var opt *redis.ClusterOptions = &redis.ClusterOptions{}

	var addrs []string

	addrs = append(addrs, fmt.Sprintf("%s:%v", host, port))
	opt.Addrs = addrs
	opt.MaxRedirects = 16

	if _conf.Debug {
		l4g.Fine(fmt.Sprintf("Connecting to %s:%v", host, port))
	}

	tmp.rediscli = redis.NewClusterClient(opt)
	tmp.prefix = prefix

	return tmp
}

// To test connectivity after the object has been created.
func (r *RedisDB) Ping() error {

	_, err := r.rediscli.Ping().Result()
	return err
}

// Register new lua script immediately after establishing a connection
func (r *RedisDB) RegisterScript(script string) *redis.Script {

	return redis.NewScript(script)
}

func (r *RedisDB) Set(key string, value string, expire time.Duration) error {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.Set(key, value, expire)

	if _conf.Debug {
		l4g.Fine("Req : Set(%s,%s) => %v", key, value, status.Err())
	}

	return status.Err()
}

func (r *RedisDB) SetB(key string, value []byte, expire time.Duration) error {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.Set(key, string(value), expire)

	if _conf.Debug {
		l4g.Fine("Req: SetB(%s,%s) => %v", key, value, status.Err())
	}

	return status.Err()
}

func (r *RedisDB) Del(key string) error {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.Del(key)

	if _conf.Debug {
		l4g.Fine("Req #%v: Del() => %v", key, status.Err())
	}

	return status.Err()
}

func (r *RedisDB) Get(key string) (value string, err error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.Get(key)

	value, err = status.Result()
	if err == redis.Nil {
		if _conf.Debug {
			l4g.Fine("Req #%v: Get() returned redis: nil, no records found", key)
		}
		err = nil
	}

	if _conf.Debug {
		l4g.Fine("Req #%v: Get() => %v,%v", key, value, err)
	}
	return
}

func (r *RedisDB) Inc(key string) (value int64, err error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.Incr(key)

	value, err = status.Result()
	if err == redis.Nil {
		if _conf.Debug {
			l4g.Fine("Req : Inc() returned redis: nil, no records found", key)
		}
		err = nil
	}

	if _conf.Debug {
		l4g.Fine("Req : Inc() => %v,%v", key, value, err)
	}
	return
}

func (r *RedisDB) Script(script *redis.Script, key string, arg int) (int64, error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}
	status := script.Run(r.rediscli, []string{key}, arg)
	value, err := status.Result()
	if err != nil{
		return 1, err
	}
	result := value.(int64)
	return result, err
}

func (r *RedisDB) GetB(key string) (value []byte, err error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.Get(key)

	svalue, serr := status.Result()
	if serr == redis.Nil {
		if _conf.Debug {
			l4g.Fine("Req #%v: GetB() returned redis: nil, no records found", key)
		}
		serr = nil
	}


	if _conf.Debug {
		l4g.Fine("Req #%v: GetB() => %v,%v", key, svalue, serr)
	}

	return []byte(svalue), serr
}

func (r *RedisDB) Exists(key string) (ok bool, err error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.Exists(key)

//	ok, err = status.Result()

	res, err:= status.Result()

	if res != 0 {

		
		ok = true
		
	}
	

	if _conf.Debug {
		l4g.Fine("Req: Exists(%s) => %v,%v", key, ok, err)
	}

	return
}

func (r *RedisDB) RPush(key string, value string) error {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.RPush(key, value)

	if _conf.Debug {
		l4g.Fine("Req: RPush(%s,%s) => %v", key, value, status.Err())
	}

	return status.Err()
}

func (r *RedisDB) LPop(key string) (value string, err error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.LPop(key)

	value, err = status.Result()

	if _conf.Debug {
		l4g.Fine("Req: LPop(%s) => %v,%v", key, value, err)
	}

	return
}

func (r *RedisDB) RPushB(key string, value []byte) error {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.RPush(key, string(value))

	if _conf.Debug {
		l4g.Fine("Req: RPush(%s,%v) => %v", key, value, status.Err())
	}

	return status.Err()
}

func (r *RedisDB) LPopB(key string) (value []byte, err error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.LPop(key)

	svalue, serr := status.Result()

	if _conf.Debug {
		l4g.Fine("Req : LPop(%s) => %v,%v", key, svalue, serr)
	}

	return []byte(svalue), serr
}

func (r *RedisDB) LIndex(key string, index int64) (value string, err error) {
	if len(r.prefix) > 0 {
		key = fmt.Sprintf("%s-{%s}", r.prefix, key)
	}

	status := r.rediscli.LIndex(key, index)

	value, err = status.Result()

	if _conf.Debug {
		l4g.Fine("Req: LIndex(%s,%v) => %v,%v", key, index, value, err)
	}
	return
}

// Returns ALL the keys of a Cluster in any order.
//
// NOTE: This function is different from the Original Scan/Keys
//       in the scope of their actions.
//
func (r *RedisDB) Scan() (keys []string, err error) {

	l4g.Fine("Scan starts")
	defer l4g.Fine("Scan ends")

	var cursor uint64
	var tmp_keys []string
	var key_prefix string

	if len(r.prefix) > 0 {
		key_prefix = fmt.Sprintf("%s*", r.prefix)
	} else {
		key_prefix = "*"
	}

	l4g.Fine("Scan: matching pattern '%s'", key_prefix)

	// Making sure the cursor is 0 the very first iteration
	cursor = 0

	for {

		l4g.Fine(fmt.Sprintf("Scan: querying with prefix '%s' and cursor %v", key_prefix, cursor))

		tmp_keys, cursor, err = r.rediscli.Scan(cursor, key_prefix, 50).Result()
		//keys, cursor_out, err = r.rediscli.Scan(cursor_in, "", 50).Result()

		if err != nil {
			l4g.Error("Scan: %s ...", err)
			return nil, err
		}

		for _, key := range tmp_keys {
			l4g.Fine(fmt.Sprintf("Scan: node has key '%s'", key))
			keys = append(keys, key)
		}

		if cursor == 0 {
			l4g.Fine(fmt.Sprintf("Scan: completed"))
			break
		}

	} // for {

	return keys, nil
}

// Returns ALL the keys of a Cluster in any order.
//
// NOTE: This function is different from the Original Scan/Keys
//       in the scope of their actions.
//
/*
func (r *RedisDB) Scan() (keys []string, err error) {

	l4g.Fine("Scan starts")
	defer l4g.Fine("Scan ends")

	var cursor uint64
	var tmp_keys []string
	var nodes []*ClusterNode
	var node *ClusterNode
	var opt *redis.ClusterOptions
	var rediscli *redis.ClusterClient
	var key_prefix string

	nodes,err=r.ClusterNodes()

	if err!=nil {
		return nil,err
	}

	if len(r.prefix)>0 {
		key_prefix=fmt.Sprintf("%s*",r.prefix)
	} else {
		key_prefix="*"
	}


	l4g.Fine("Scan: matching pattern '%s'",key_prefix)

	// We are going to query each Master node in order to retrieve all
	// the keys of the cluster. To do that, we must request each slot,
	// and then for each master in charge of a specific slot, we scan
	// the available keys.
	for _,node = range nodes {

		if !node.IsMaster {
			l4g.Fine("Scan: Skipping slave '%s'",node.Addr)
			continue
		}

		addrs:=make([]string,1)
		addrs=append(addrs,node.Addr)

		opt=&redis.ClusterOptions{}

		opt.Addrs=addrs
		opt.MaxRedirects=16
		rediscli = redis.NewClusterClient(opt)

		_, err = rediscli.Ping().Result()

		if err!=nil {
			l4g.Error("Scan: Error while connecting to Master server '%s' : '%s'",node.Addr,err)
			return nil,err

		}

		// Making sure the cursor is 0 the very first iteration
		cursor=0


		for {

			l4g.Fine(fmt.Sprintf("Scan: querying master '%s' with prefix '%s' and cursor %v",node.Addr,key_prefix,cursor))

			tmp_keys, cursor, err = rediscli.Scan(cursor, key_prefix, 50).Result()
			//keys, cursor_out, err = r.rediscli.Scan(cursor_in, "", 50).Result()

			if err != nil {
				l4g.Error("Scan: %s ...",err)
				return nil,err
			}

			for _,key:=range tmp_keys {
				l4g.Fine(fmt.Sprintf("Scan: master '%s' has key '%s'",node.Addr,key))
				keys=append(keys,key)
			}

			if cursor==0 {
				l4g.Fine(fmt.Sprintf("Scan: master '%s' has completed",node.Addr))
				break
			}

		} // for {

		err=rediscli.Close()

		if err!=nil {
			l4g.Warn("Scan: error while closing temporary client: '%s' ...",err)
		}

		rediscli=nil


	} // for _,node = range nodes {
	return keys,nil
}
*/

func (r *RedisDB) ClusterNodes() (nodes []*ClusterNode, err error) {
	l4g.Fine("ClusterNodes starts")
	defer l4g.Fine("ClusterNodes ends")

	status := r.rediscli.ClusterNodes()

	var str_nodes string

	str_nodes, err = status.Result()

	if err != nil {
		l4g.Error("ClusterNodes: %s ...", err)
		return nil, err
	}

	nodes, _ = parseNodes(str_nodes)

	return nodes, nil

}

/**
80ad7ec5ad1457facd82d2e97d19f3a0edff8d36 127.0.0.1:7002 myself,master - 0 0 3 connected 10923-16383
0aa1418e56b477e194d59006c1a879561542e032 127.0.0.1:7005 slave 80ad7ec5ad1457facd82d2e97d19f3a0edff8d36 0 1472320896108 6 connected
0862efc4cf6054b9b8c992c8bef6eb0a88dda358 127.0.0.1:7003 slave 8b60f92dd5c88422999be118e74f4af6d91fd96c 0 1472320897109 4 connected
8b60f92dd5c88422999be118e74f4af6d91fd96c 127.0.0.1:7000 master - 0 1472320897109 1 connected 0-5460
06ff5b9d7bc7c702e83f62bc193b8df12397949f 127.0.0.1:7004 slave 973993d7bf548c2ae05d2ef2864bda4ea41861bc 0 1472320895607 5 connected
973993d7bf548c2ae05d2ef2864bda4ea41861bc 127.0.0.1:7001 master - 0 1472320896608 2 connected 5461-10922

*/
func parseNodes(str_nodes string) (nodes []*ClusterNode, slots []*ClusterSlot) {
	var id2node map[string]*ClusterNode
	var id2master map[string]string
	var node *ClusterNode
	var token_nb int

	l4g.Fine("parseNodes starts")
	defer l4g.Fine("parseNodes ends")

	id2node = make(map[string]*ClusterNode)
	id2master = make(map[string]string)

	l4g.Fine(fmt.Sprintf("parseNodes: Analyzing  '%s'", str_nodes))

	str_nodes_list := strings.Split(str_nodes, "\n")

	for _, line := range str_nodes_list {
		l4g.Fine(fmt.Sprintf("parseNodes: parsing line '%s'", line))

		if len(line) == 0 {
			continue
		}

		tokens := strings.Split(line, " ")

		token_nb = len(tokens)

		if token_nb < 8 {
			l4g.Error(fmt.Sprintf("parseNodes: Invalid line '%s' as at least 8 tokens are expected when only %v are found.", line, token_nb))
			return nil, nil
		}

		node = NewClusterNode()

		node.Id = tokens[0]
		node.Addr = tokens[1]

		id2node[node.Id] = node
		nodes = append(nodes, node)

		subtokens := strings.Split(tokens[2], ",")

		if (strings.Compare(tokens[2], "master") == 0) || (len(subtokens) == 2 && strings.Compare(subtokens[1], "master") == 0) {
			l4g.Fine(fmt.Sprintf("parseNodes: '%s' is master", node.Id))

			if token_nb < 9 {
				l4g.Error(fmt.Sprintf("parseNodes: Invalid line '%s' as 9 tokens are expected when only %v are found", line, token_nb))
				return nil, nil
			}

			node.IsMaster = true
			slot := NewClusterSlot()

			ranges := strings.Split(tokens[8], "-")

			start, err := strconv.ParseInt(ranges[0], 10, 64)

			if err != nil {
				l4g.Error(fmt.Sprintf("parseNodes: failed to parse int from '%s' from token '%s' and from line '%s'", ranges[0], tokens[8], line))
				return nil, nil
			}

			slot.Start = start

			end, err := strconv.ParseInt(ranges[1], 10, 64)

			if err != nil {
				l4g.Error(fmt.Sprintf("parseNodes: failed to parse int from '%s' from token '%s' and from line '%s'", ranges[1], tokens[8], line))
				return nil, nil
			}

			slot.End = end
			slot.Master = node

			slots = append(slots, slot)

		} else {
			// We have a slave.
			// Mapping the ID of the slave with the ID of the master
			// so we can connect with the right node after the loop
			// completes.
			id2master[node.Id] = tokens[3]

			l4g.Fine(fmt.Sprintf("parseNodes: '%s' is slave", node.Id))
		}

	} // for i,line:= range str_nodes_list {

	for _, node = range nodes {
		if node.IsMaster {
			continue
		}

		master_id, ok := id2master[node.Id]

		if !ok {
			l4g.Error(fmt.Sprintf("parseNodes: failed to find the master node for slave node '%s'", node.Id))
			return nil, nil
		}
		master_node, ok := id2node[master_id]

		if !ok {
			l4g.Error(fmt.Sprintf("parseNodes: failed to find the master node with id '%s' for slave node '%s'", master_id, node.Id))
			return nil, nil
		}

		node.Master = master_node
	}

	return nodes, slots
}

func (r *RedisDB) ClusterSlots() (slots []*ClusterSlot, err error) {
	l4g.Fine("ClusterSlots starts")
	defer l4g.Fine("ClusterSlots ends")

	status := r.rediscli.ClusterNodes()

	var str_nodes string

	str_nodes, err = status.Result()

	if err != nil {
		l4g.Error("ClusterSlots: %s ...", err)
		return nil, err
	}

	_, slots = parseNodes(str_nodes)

	return slots, nil
}

func NewClusterNode() *ClusterNode {
	var tmp *ClusterNode = &ClusterNode{}

	return tmp
}

func NewClusterSlot() *ClusterSlot {
	var tmp *ClusterSlot = &ClusterSlot{}

	return tmp
}
