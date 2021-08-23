package main

// Providers
import (
	_ "github.com/memsdm05/nplink/providers/nightbot"
	_ "github.com/memsdm05/nplink/providers/streamelements"

)

import (
	"fmt"
	"github.com/memsdm05/nplink/internal/app"
	"github.com/memsdm05/nplink/internal/setup"
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

	app.Close()
}
