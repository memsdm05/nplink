package main

import (
	// register providers
	_ "github.com/memsdm05/nplink/provider/nightbot"
	"github.com/memsdm05/nplink/setup"
	"runtime/debug"

	"fmt"
	"github.com/memsdm05/nplink/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("An error has occurred\n")
			fmt.Printf("panic: %s\n\n%s", err, string(debug.Stack()))
			fmt.Println("\n< Press Enter to exit >")
			fmt.Scanln()
			os.Exit(1)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM)

	setup.Auth()
	fmt.Println("here")
	go app.MainLoop()
	<-sigs
	fmt.Println("exiting...")
}