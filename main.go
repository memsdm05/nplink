package main

import (
	"fmt"
	"github.com/memsdm05/nplink/app"
	"github.com/memsdm05/nplink/setup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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

	go app.MainLoop()
	<-sigs
	fmt.Println("exiting...")
}