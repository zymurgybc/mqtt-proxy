package conf

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type MysqlConfiguration struct {
	ConnectionString string `toml:"connection-string"`
	Select           string `toml:"select"`
}

type HttpConfiguration struct {
	ListenAddress string `toml:"listen-address"`
}

type MqttConfiguration struct {
	ListenAddress string `toml:"listen-address"`
	Cert          string `toml:"cert"`
	Key           string `toml:"key"`
}

type Configuration struct {
	BackendServers []string `toml:"backend-servers"`
	User           string   `toml:"user"`
	Pass           string   `toml:"pass"`

	Http           HttpConfiguration  `toml:"http"`
	MqttStoreMysql MysqlConfiguration `toml:"mqtt-store"`
	WsStoreMysql   MysqlConfiguration `toml:"ws-store"`
	Mqtt           MqttConfiguration  `toml:"mqtt"`
}

func LoadConfiguration(fileName string) *Configuration {
	config, err := parseTomlConfiguration(fileName)
	if err != nil {
		log.Println("Couldn't parse configuration file: " + fileName)
		panic(err)
	}
	return config
}

func parseTomlConfiguration(filename string) (*Configuration, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	tomlConfiguration := &Configuration{}
	_, err = toml.Decode(string(body), tomlConfiguration)
	if err != nil {
		return nil, err
	}
	if len(tomlConfiguration.BackendServers) == 0 {
		return nil, errors.New("At least one backend servers required.")
	}
	if tomlConfiguration.Http.ListenAddress == "" {
		tomlConfiguration.Http.ListenAddress = ":9000"
	}
	if tomlConfiguration.Mqtt.ListenAddress == "" {
		tomlConfiguration.Mqtt.ListenAddress = ":1883"
	}
	if tomlConfiguration.User == "" {
		tomlConfiguration.User = "guest"
	}
	if tomlConfiguration.Pass == "" {
		tomlConfiguration.Pass = "guest"
	}

	// need a way to merge defaults..
	if tomlConfiguration.MqttStoreMysql.ConnectionString == "" {
		tomlConfiguration.MqttStoreMysql.ConnectionString = "root:@tcp(127.0.0.1:3306)/mqtt"
	}
	if tomlConfiguration.MqttStoreMysql.Select == "" {
		tomlConfiguration.MqttStoreMysql.Select = "select uid, mqtt_id from users where mqtt_id = ?"
	}

	if tomlConfiguration.WsStoreMysql.ConnectionString == "" {
		tomlConfiguration.WsStoreMysql.ConnectionString = "root:@tcp(127.0.0.1:3306)/mqtt"
	}
	if tomlConfiguration.WsStoreMysql.Select == "" {
		tomlConfiguration.WsStoreMysql.Select = "select uid, mqtt_id from users where mqtt_id = ?"
	}

	return tomlConfiguration, nil
}
