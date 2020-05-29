package config

import (
	"os"

	l4g "code.google.com/p/log4go"
	"github.com/BurntSushi/toml"
)

var (
	_repconf         *AutoEndpointConfig
	_last_modif_time int64
	_config_file     string
)

func init() {

	_repconf = NewREPConfig()

	// For default config, the last modification time is zero.
	_last_modif_time = 0

	// For default config, the file name is empty.
	_config_file = ""

}

func GetREPConfig() *AutoEndpointConfig {
	return _repconf
}

func Load_REPConfig(config_file string) (conf *AutoEndpointConfig, err error) {

	var fileinfo os.FileInfo
	var local_config *AutoEndpointConfig

	fileinfo, err = os.Stat(config_file)
	if err != nil {
		l4g.Error("Specified configuration file '%v' not found, error %s", config_file, err.Error())

		return
	}

	local_config = NewREPConfig()

	_config_file = config_file
	_last_modif_time = fileinfo.ModTime().UnixNano()

	if _, err = toml.DecodeFile(config_file, local_config); err != nil {
		l4g.Error("Reading configuration file '%v' leads to error '%s'!", config_file, err.Error())
		return
	}

	_repconf = local_config
	conf = _repconf
	return

}

func NewREPConfig() *AutoEndpointConfig {
	var config = &AutoEndpointConfig{

		Endpoint_Hostname: "localhost",
		Endpoint_Port:     8080,
		Debug:             false,
		Cass_Address:      "localhost",
		Keyspace:          "autopush",
		Max_Payload:       1000,
		Statsd_Host:       "localhost",
		Statsd_Port:       8125,
		Max_Msg:           1000,
		Max_Msg_Length:    4096,
	}

	return config

}
