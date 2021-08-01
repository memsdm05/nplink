package setup

import (
	"github.com/BurntSushi/toml"
	"github.com/memsdm05/nplink/util"
)

var Config config

func init() {
	_, err := toml.DecodeFile("config.toml", &Config)
	if err != nil {}
}

type config struct {
	Provider string
	Timeout float32
	Address string
	SkipAuth bool `toml:"auto_authorize"`
	Commands []struct {
		Name      string
		Format    util.FormatString
	} `toml:"command"`
}
