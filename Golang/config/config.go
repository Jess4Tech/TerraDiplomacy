package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type redisConfig struct {
	Host     string `json:"host"`
	Password string `json:"pass"`
}

type mySQLConfig struct {
	Address  string `json:"address"`
	User     string `json:"host"`
	Password string `json:"pass"`
	Database string `json:"db"`
}

type sslConfig struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

type addressConfig struct {
	Scheme string `json:"scheme"`
	Domain string `json:"domain"`
	Port   int    `json:"port"`
}

func (ac *addressConfig) Address() string {
	return fmt.Sprintf("%s:%d", ac.Domain, ac.Port)
}

func (ac *addressConfig) SchemedAddress() string {
	return fmt.Sprintf("%s://%s:%d", ac.Scheme, ac.Domain, ac.Port)
}

func (ac *addressConfig) Secure() bool {
	return ac.Scheme == "https"
}

type config struct {
	ListenAddress   addressConfig `json:"listenAddress"`
	FrontendAddress addressConfig `json:"frontendAddress"`
	TestKey         string        `json:"testKey"`
	Administrators  []string      `json:"admins"`
	RedisConfig     redisConfig   `json:"redis"`
	MySQLConfig     mySQLConfig   `json:"mySql"`
	SslConfig       sslConfig     `json:"ssl"`
}

var configFile = "config.json"

var Config = config{}

func LoadConfiguration() {
	configFileHandle, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(configFile)
			if err != nil {
				panic(err)
			}
			log.Println("Creating default configuration file (config.json)")
			json, err := json.MarshalIndent(config{
				ListenAddress: addressConfig{
					Scheme: "http",
					Domain: "localhost",
					Port:   10000,
				},
				FrontendAddress: addressConfig{
					Scheme: "http",
					Domain: "localhost",
					Port:   10001,
				},
				TestKey:        "",
				Administrators: []string{},
				RedisConfig: redisConfig{
					Host:     "localhost:6379",
					Password: "password",
				},
				MySQLConfig: mySQLConfig{
					Address:  "localhost:3306",
					User:     "terra_user",
					Password: "password",
					Database: "terra",
				},
				SslConfig: sslConfig{
					Certificate: "localhost.crt",
					Key:         "localhost.key",
				},
			}, "", "    ")
			if err != nil {
				panic(err)
			}
			file.Write(json)
			log.Println("Default configuration file created, edit as necessary")
		}
	} else {
		defer configFileHandle.Close()
		dec := json.NewDecoder(configFileHandle)
		dec.Decode(&Config)
	}
}
