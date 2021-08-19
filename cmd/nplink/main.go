package main

import (
	"github.com/memsdm05/nplink/internal/app"
	"github.com/memsdm05/nplink/internal/setup"
	_ "github.com/memsdm05/nplink/providers/nightbot"
)

import (
	"fmt"
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

	setup.Load()
	setup.Auth()

	go app.MainLoop()
	<-sigs
	fmt.Println("exiting...")
}
