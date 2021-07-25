package setup

import (
	"github.com/BurntSushi/toml"
)

var Config config

func init() {
	toml.DecodeFile("config.toml", &Config)
}

type config struct {
	Provider string
	Timeout float32
	Address string
	SkipAuth bool `toml:"auto_authorize"`
	//Commands []struct {
	//	Name      string
	//	Format    util.FormatString
	//	Cooldown  int
	//	Userlevel string
	//} `toml:"command"`
	Commands []map[string]interface{} `toml:"command"`
}