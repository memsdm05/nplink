package main

import (
	"fmt"
	"github.com/memsdm05/nplink/setup"
)

func main() {
	fmt.Println(func(commands []map[string]interface{}) bool {
		for _, m := range commands {
			_, okname := m["name"].(string)
			_, okformat := m["format"].(string)
			if !(okname || okformat) {
				return false
			}
		}
		return true
	}(setup.Config.Commands))
	fmt.Println(setup.Config.Commands[0]["name"].(string))
	//app.MainLoop()
}