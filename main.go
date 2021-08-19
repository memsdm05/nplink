package main

import (
	_ "github.com/memsdm05/nplink/provider/nightbot"
	_ "github.com/memsdm05/nplink/provider/streamelements"
)

import (
	"fmt"
	// register providers
	"github.com/memsdm05/nplink/app"
	"github.com/memsdm05/nplink/setup"
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

	setup.Auth()
	go app.MainLoop()
	<-sigs
	fmt.Println("exiting...")
}
