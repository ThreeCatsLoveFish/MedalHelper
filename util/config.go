package util

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var GlobalConfig Config

type Config struct {
	UserList  []User     `yaml:"USERS"`
	Danmuku   []string   `yaml:"DANMU"`
	Endpoints []Endpoint `yaml:"PUSH"`
	CD        CoolDown   `yaml:"CD"`
	Cron      string     `yaml:"CRON"`
}

type CoolDown struct {
	Async  int `yaml:"async"`
	Retry  int `yaml:"retry"`
	MaxTry int `yaml:"max_try"`
	Like   int `yaml:"like"`
	Share  int `yaml:"share"`
	Danmu  int `yaml:"danmu"`
}

type Endpoint struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

type User struct {
	AllowedUid string `yaml:"allowed_uid"`
	BannedUid  string `yaml:"banned_uid"`
	AccessKey  string `yaml:"access_key"`
	PushName   string `yaml:"push_name"`
}

// InitConfig bind endpoints with config file
func InitConfig() {
	conf := config.NewWithOptions("push", func(opt *config.Options) {
		opt.DecoderConfig.TagName = "yaml"
		opt.ParseEnv = true
	})
	conf.AddDriver(yaml.Driver)
	err := conf.LoadFiles("users.yaml")
	if err != nil {
		panic(err)
	}
	// Load config file
	err = conf.BindStruct("", &GlobalConfig)
	if err != nil {
		panic(err)
	}
}

// LoadConfig bind endpoints with config file
func LoadConfig(filePath string) {
	conf := config.NewWithOptions("push", func(opt *config.Options) {
		opt.DecoderConfig.TagName = "yaml"
		opt.ParseEnv = true
	})
	conf.AddDriver(yaml.Driver)
	err := conf.LoadFiles(filePath)
	if err != nil {
		panic(err)
	}
	// Load config file
	err = conf.BindStruct("", &GlobalConfig)
	if err != nil {
		panic(err)
	}
}
