package util

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var GlobalConfig Config

func init() {
	initConfig()
}

type Config struct {
	UserList  []User     `yaml:"USERS"`
	Danmaku   []string   `yaml:"DANMU"`
	Endpoints []Endpoint `yaml:"PUSH"`
	Cron      string     `yaml:"CRON"`
}

type Endpoint struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

type User struct {
	BannedUid string `yaml:"banned_uid"`
	AccessKey string `yaml:"access_key"`
	PushName  string `yaml:"push_name"`
}

// initConfig bind endpoints with config file
func initConfig() {
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
	conf.BindStruct("", &GlobalConfig)
}
