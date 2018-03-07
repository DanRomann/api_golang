package Utils

import (
	"time"
	"log"
	"encoding/json"
	"io/ioutil"
)

type ServerConf struct{
	Port			string			`json:"port"`
	StaticDir		string			`json:"static_dir"`
	WriteTimeout	time.Duration	`json:"write_timeout"`
	ReadTimeout		time.Duration	`json:"read_timeout"`
	Address			string			`json:"address"`
}

type DBConf struct{
	Address		string	`json:"address"`
	Port		string	`json:"port"`
	User		string	`json:"user"`
	Password	string	`json:"password"`
	Dbname		string	`json:"dbname"`
	SSLmode		string	`json:"ssl_mode"`
}

type FrontConf struct {
	Addr		string	`json:"addr"`
}

type EmailConf struct {
	Addr	string	`json:"addr"`
}

type Config struct{
	ServerConf	`json:"server"`
	DBConf		`json:"db"`
	FrontConf	`json:"front"`
	EmailConf	`json:"email"`
}

var MainConfig Config

func init() {
	config := new(Config)
	raw, err := ioutil.ReadFile("./dev_config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		log.Fatal("Init config ", err)
	}
	MainConfig = *config
}