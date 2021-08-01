package main

import (
	"errors"
	"github.com/memsdm05/nplink/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	util.FatalError(errors.New("here"))

	//go app.MainLoop()
	//<-sigs
	//fmt.Println("exiting...")
}