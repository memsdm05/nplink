package app

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/memsdm05/nplink/internal/provider"
	"github.com/memsdm05/nplink/internal/setup"
	"github.com/memsdm05/nplink/internal/utils"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	shutdown = make(chan bool, 1)
	changes = make(chan commandChange, 100)
	services = []packet{
		new(gosumemoryPacket),
		//new(streamCompanionPacket),
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

func providerRunner(prov provider.Provider, changes <-chan commandChange) {
	for change := range changes {
		for {
			if err := prov.SetCommand(change.name, change.content); err == nil {
				break
			}
			log.Println("Something bad happened with command set, retrying")
			time.Sleep(1 * time.Millisecond)
		}

		log.Printf("set %s to \"%s\"\n", change.name, change.content)

		time.Sleep(700 * time.Millisecond)
	}

	shutdown <- true
}

func newConnection(p packet) *websocket.Conn {
	connUrl := fmt.Sprintf("ws://%s/%s", setup.Config.Address, service.path())
	for retry := 1;; retry++ {
		if ok, _ := <-shutdown; ok {
			return nil
		}

		c, _, err := websocket.DefaultDialer.Dial(connUrl, nil)

		if err != nil {
			return c
		}

		if retry == 1 {
			fmt.Printf("Connection problem, make sure you have a memory scanner running \n")
		}

		fmt.Printf("have tried %d time%s, retrying again in 3 seconds...\n", retry,
			func() string{
				if retry == 1 {
					return ""
				}
				return "s"
			}())
		time.Sleep(3 * time.Second)
	}
}

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

	go providerRunner(setup.SelectedProvider, changes)

	c := newConnection(service)
	defer c.Close()

	mainloop:
	for {
		select {
		case <-shutdown:
			break mainloop
		default:
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

				// we first "mark" when a command's content is changed
				// is this stored in the format tracker
				if newC != old.c {
					old.t = time.Now()
					old.c = newC
					old.s = true
				}

				// once a sufficient amount of time has passed, consume that tracked change
				if old.s && time.Since(old.t) > changeWait {
					changes <- commandChange{
						name:    command.Name,
						content: newC,
					}
					old.s = false
				}
			}

			// we want strings to be instantly formatted on fist iteration
			// but need and timeout on later ones
			// quick and dirty

			if first {
				changeWait = time.Duration(setup.Config.Timeout*1000) * time.Millisecond
				first = false
			}
		}
	}
}

func Close() {
	shutdown <- true

	for _, command := range setup.Config.Commands {
		changes <- commandChange{
			name:    command.Name,
			content: "nplink is offline",
		}
	}
	close(changes)

	<-shutdown
}