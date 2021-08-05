package main

import (
	"fmt"
	"github.com/memsdm05/nplink/app"
	// register providers
	_ "github.com/memsdm05/nplink/provider/nightbot"
	"github.com/memsdm05/nplink/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM)


	util.SetCred("foo", "bar")
	fmt.Println(util.GetCred("foo"))

	go app.MainLoop()
	<-sigs
	fmt.Println("exiting...")
}