package app

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/memsdm05/nplink/setup"
	"github.com/memsdm05/nplink/utils"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	changes = make(chan commandChange, 100)
	services = []packet{
		new(gosumemoryPacket),
		new(streamCompanionPacket),
	}
	service packet
)

type commandChange struct {
	name    string
	content string
}

type packet interface {
	check(data []byte) bool
	fill(conn *websocket.Conn) error
	flatten(fmap utils.FMap)
	path() string
}

func providerRunner(changes <-chan commandChange) {
	prov := setup.SelectedProvider
	for change := range changes {
		fmt.Printf("set %s to \"%s\"\n", change.name, change.content)
		if err := prov.SetCommand(change.name, change.content); err != nil{
			panic(err)
		}
		time.Sleep(700 * time.Millisecond)
	}
}

/*
           bid    sid
Main menu   0      0
Ranked Map  >1     >1
N Submitted 0      0
Practice    0      0
*/

func MainLoop() {
	fmap := make(utils.FMap)

	formatTracker := make([]struct {
		t time.Time
		c string
		s bool
	}, len(setup.Config.Commands))

	changeWait := time.Duration(-1)
	first := true

	// check which service we are using
	resp, err := http.Get(fmt.Sprintf("http://%s/json",
		setup.Config.Address))

	if err == nil {
		b, _ := io.ReadAll(resp.Body)
		for _, s := range services {
			if s.check(b) {
				service = s
				break
			}
		}
	} else {
		panic(err)
	}

	go providerRunner(changes)

	c, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://%s/%s", setup.Config.Address, service.path()),
		nil)

	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		err = service.fill(c)
		if err != nil {
			log.Println(err)
			// todo real error handling
			continue
		}

		service.flatten(fmap)

		for i, command := range setup.Config.Commands {
			newC := command.Format.Format(fmap)
			old := &formatTracker[i]

			// change only after content is different for changeWait milliseconds

			if newC != old.c {
				old.t = time.Now()
				old.c = newC
				old.s = true
			}

			if old.s && time.Since(old.t) > changeWait {
				changes <- commandChange{
					name:    command.Name,
					content: newC,
				}
				old.s = false
			}
		}

		//fmt.Println(time.Since(formatTracker[0].t))

		// we want strings to be instantly formatted on fist iteration
		// but need and timeout on later ones
		// quick and dirty

		if first {
			changeWait = time.Duration(setup.Config.Timeout*1000) * time.Millisecond
			first = false
		}
	}
}

func Close() {
	
}