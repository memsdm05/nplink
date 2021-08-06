package setup

import (
	"github.com/BurntSushi/toml"
	"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/utils"
	"strings"
)

var Config config
var SelectedProvider provider.Provider

func init() {
	_, err := toml.DecodeFile("config.toml", &Config)
	if err != nil {} // lol

	prov, _ := provider.Select(strings.ToLower(Config.Provider))
	SelectedProvider = prov
}

type config struct {
	Provider string
	Timeout float32
	Address string
	SkipAuth bool `toml:"auto_authorize"`
	Commands []struct {
		Name   string
		Format utils.FormatString
	} `toml:"command"`
}
