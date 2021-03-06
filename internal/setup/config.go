package setup

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/memsdm05/nplink/internal/provider"
	"github.com/memsdm05/nplink/internal/utils"
	"strings"
)

var Config config
var SelectedProvider provider.Provider

func Load() {
	_, err := toml.DecodeFile("config.toml", &Config)
	if err != nil {
		panic("this will generate a new config in the future\nnow it fucking dies")
	} // for right now

	prov, err := provider.Select(strings.ToLower(Config.Provider))
	if err != nil {
		panic(fmt.Sprintf("%s is not a valid provider", Config.Provider))
	}
	SelectedProvider = prov
}

type config struct {
	Provider string
	Timeout  float32 `toml:"change_wait"`
	Address  string
	SkipAuth bool `toml:"auto_authorize"`
	Commands []struct {
		Name   string
		Format utils.FormatString
	} `toml:"command"`
}
