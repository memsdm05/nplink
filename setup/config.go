package setup

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var Config config

func init() {
	meta, _ := toml.DecodeFile("config.toml", &Config)
	fmt.Println(meta.Keys())
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