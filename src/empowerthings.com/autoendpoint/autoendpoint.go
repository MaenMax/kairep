package main

import (
	"errors"
	fernet "fernet-go"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"empowerthings.com/autoendpoint/config"
	"empowerthings.com/autoendpoint/security/jwttoken"
	"github.com/gocql/gocql"
	"github.com/statsd"
	"golang.org/x/exp/utf8string"

	//	"runtime"
	"context"
	"net"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"

	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/push_routes"

	"empowerthings.com/autoendpoint/routercore/db/redis/counterdb"
	"empowerthings.com/autoendpoint/routercore/db/redis/redisdb"
)

var (
	_conf          *config.AutoEndpointConfig
	config_file    *string = flag.String("config-endpoint", "configs/autopush_endpoint.conf", "Config file to use.")
	log_file       *string = flag.String("log", "configs/autoendpoint_log.xml", "Log L4G config file to use.")
	http_server    *http.Server
	crypto_key     string
	cass_user      string
	cass_password  string
	cass_session   *gocql.Session
	redis_password string
	c_kill         chan os.Signal
	c_int          chan os.Signal
	cpuprofile     = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile     = flag.String("memprofile", "", "write memory profile to this file")
	redis          *redisdb.RedisDB
)

const (
	CRYPTO_KEY       string        = "CRYPTO_KEY"
	CASS_USERNAME    string        = "CASS_USERNAME"
	CASS_PASSWORD    string        = "CASS_PASSWORD"
	TTL_DURATION     time.Duration = 1000 * 1000 * 1000 * 60 * 60 * 24 * 365
	MAX_ITER         int           = 10000
	SHUTDOWN_TIMEOUT               = 10
)

func Init(config_file string, log_file string) error {

	// Initiating environement variables.

	c_kill = make(chan os.Signal, 1)
	c_int = make(chan os.Signal, 1)

	crypto_key = os.Getenv("CRYPTO_KEY")
	if crypto_key == "" {
		panic("[ERROR]: No Crypto Key Found. Aborting ...")
	}

	cass_user = os.Getenv("CASS_USERNAME")
	if crypto_key == "" {
		panic("[ERROR]: Cassandra username is not set. Aborting ...")
	}

	cass_password = os.Getenv("CASS_PASSWORD")
	if cass_password == "" {
		panic("[ERROR]: Cassandra password is not set. Aborting ...")
	}

	var config_read bool = false

	l4g.LoadConfiguration(log_file)

	if file, err := os.OpenFile(config_file, os.O_RDONLY, 0666); err == nil {
		_conf, err = config.Load_REPConfig(config_file)
		if err == nil {
			l4g.Info("Read config from '%s' ...", config_file)
			config_read = true
		} else {
			l4g.Error("Failed to Read config from '%s': '%s'.", config_file, err)
			panic(err)
		}
		file.Close()
	}

	if !config_read {
		panic("No 'autopush_endpoint.conf' configuration file found!")
	}

	// Making sure to have the latest configuration.
	_conf = config.GetREPConfig()
	// Initiating Cassandra DB
	err := start_cassandra()
	if err != nil {
		panic(err)
	}

	l4g.Info("Connected to Cassandra DB successfully.")

	if _conf.Max_Msg != 0 { // No msg limitation applied. Don't connect to redis !!

		redis_password = os.Getenv("REDIS_PASSWORD")
		if redis_password == "" {
			panic("[ERROR]: Redis password is not set. Aborting ...")
		}

		l4g.Info("Connecting to Redis...")

		decrypted_redis_pass, err := decrypt_password(redis_password, crypto_key, "Redis")

		if err != nil {
			msg := fmt.Sprintf("Redis error - %s", err)
			l4g.Error(msg)
			panic(msg)
		}

		if err = redisdb.Init(_conf, string(decrypted_redis_pass)); err != nil {
			msg := fmt.Sprintf("Redis error - %s", err)
			l4g.Error(msg)
			panic(msg)
		}

		if err = counterdb.Init(); err != nil {
			msg := fmt.Sprintf("Redis error - %s", err)
			l4g.Error(msg)
			panic(msg)
		}

		l4g.Info("Connected to Redis DB successfully.")

	} else {
		l4g.Info("Message limitation is disabled. No Redis connection is established.")
	}

	return nil
}

func main() {

	flag.Parse()

	var err error

	if *cpuprofile != "" {
		file, err := os.Create(*cpuprofile)
		if err != nil {
			panic("Must provide CPU profile")
		}
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	err = jwttoken.Init()

	if err != nil {

		l4g.Error("Error initializing token checker. %v", err.Error())

	}

	err = Init(*config_file, *log_file)

	if err != nil {
		fmt.Printf("Error during initialization: '%s'! Aborting ...", err)
		return
	}

	statsd_addr := fmt.Sprintf("%s:%v", _conf.Statsd_Host, _conf.Statsd_Port)

	l4g.Info("GoREP is connecting to StatsD Host: %s", statsd_addr)

	c, err := statsd.New(statsd.Prefix("gorep"), statsd.Address(statsd_addr)) // Connect to the UDP port 8125 by default.

	if err != nil {

		l4g.Error("[WARNING] GoREP was not able to start StatsD client : %s", err)

	} else {

		l4g.Info("GoREP is connected to StatsD UDP port: %v", _conf.Statsd_Port)

	}

	signal.Notify(c_kill, syscall.SIGTERM)
	signal.Notify(c_int, syscall.SIGINT)

	err_chan := start_server(_conf, crypto_key, c)

	select {
	case <-c_kill:
		l4g.Info("SIGTERM signal received! Graceful shutdown initiated ...")
		stop_server(_conf)

	case <-c_int:
		l4g.Info("SIGINT (Ctrl+c) detected! Graceful shutdown initiated ...")
		stop_server(_conf)

	case err = <-err_chan:
		l4g.Error("Server graceful shutdown due to error '%s'", err.Error())
		stop_server(_conf)
	}

	if *memprofile != "" {
		time.Sleep(time.Second * 5)
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Printf("Error while creating memory profile: '%s'.", err)
		}

		pprof.WriteHeapProfile(f)
		f.Close()
		return

	}

	l4g.Info("Daemon graceful server shutdown completed")
	l4g.Close()

}

func start_server(conf *config.AutoEndpointConfig, crypto_key string, stats *statsd.Client) chan error {

	errs := make(chan error)

	start_http(conf, crypto_key, stats, &errs)

	return nil

}

func start_http(conf *config.AutoEndpointConfig, crypto_key string, stats *statsd.Client, errs *chan error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.Endpoint_Port))

	if err != nil {
		l4g.Error("start_http  err: %v", err.Error())
		*errs <- err
		return
	}

	router := push_routes.NewRouter(crypto_key, conf.Debug, cass_session, conf.Keyspace, conf.Max_Payload, stats, conf.Max_Msg,conf.Router_Table_Name)

	http_server = &http.Server{
		// 2017/08/16 - RS - Now using a Listener instead (see above).
		// The reason for using a listener is to listen on both IPv6 and IPv4
		// if available.
		//Addr:      fmt.Sprintf(":%v", conf.FrontLayer.HttpPort),
		Handler: router,
	}

	go func() {

		l4g.Info("KAIOS push HTTP service has started on port %v.", conf.Endpoint_Port)

		if err := http_server.Serve(listener); err != nil {
			l4g.Error("start_http err: %v", err.Error())
			*errs <- err
		}

		fmt.Println("************KAIOS PUSH SERVICE STARTED************")
		fmt.Println("*                                                *")
		fmt.Println("*                                                *")
		fmt.Println("*                 REP NODE                       *")
		fmt.Println("*                                                *")
		fmt.Println("*                                                *")
		fmt.Println("*                                                *")
		fmt.Println("*                                                *")
		fmt.Println("*                                                *")
		fmt.Println("*                                                *")
		fmt.Println("**************************************************")

	}()

}

func start_cassandra() (err error) {

	hosts := strings.Split(_conf.Cass_Address, ",")

	cluster := gocql.NewCluster(hosts...)

	cluster.Keyspace = _conf.Keyspace
	cluster.Consistency = gocql.Quorum

	decrypted_cass_pass, err := decrypt_password(cass_password, crypto_key, "Cassandra")
	if err != nil {

		return err
	}

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cass_user,
		Password: string(decrypted_cass_pass),
	}
	cass_session, err = cluster.CreateSession()
	if err != nil {

		l4g.Error("start_cassandra: Failed due to error: '%s'", err)
		return err
	}
	//	defer cass_session.Close()
	return nil
}

func stop_server(conf *config.AutoEndpointConfig) {

	l4g.Info("Shutting down HTTP service on %v ...", conf.Endpoint_Port)

	ctx, _ := context.WithTimeout(context.Background(), SHUTDOWN_TIMEOUT*time.Second)
	http_server.Shutdown(ctx)

	l4g.Info("HTTPS service on %v shutdown ...", conf.Endpoint_Port)
}
func decrypt_password(password string, crypto_key string, db string) ([]byte, error) {

	crypto_key_utf8_encoded := utf8string.NewString(crypto_key)

	key := fernet.MustDecodeKeys(crypto_key_utf8_encoded.String())

	decrypted_pass := fernet.VerifyAndDecrypt([]byte(password), TTL_DURATION, key)
	if decrypted_pass == nil {

		err_str := fmt.Sprintf("[Error]: Can not decrypt %s password.", db)
		decryption_err := errors.New(err_str)
		return nil, decryption_err
	}
	return decrypted_pass, nil

}
