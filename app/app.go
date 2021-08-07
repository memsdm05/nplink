package app

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/memsdm05/nplink/setup"
	"github.com/memsdm05/nplink/utils"
	"log"
	"time"
)

type commandChange struct {
	name    string
	content string
}

// Packet Not the full packet as some information is unneeded
type Packet struct {
	Settings struct {
		Folders struct {
			Skin string
		}
	}
	Menu struct {
		State    int
		GameMode int
		BeatMap  struct {
			Id           int
			Set          int
			Md5          string
			RankedStatus int
			Metadata     struct {
				Artist     string
				Title      string
				Mapper     string
				Difficulty string
			}
			Stats struct {
				AR  float32
				CS  float32
				OD  float32
				HP  float32
				SR  float32
				BPM struct {
					Min float64
					Max float64
				}
				FullSR float32
			}
		} `json:"bm"`
		Mods struct {
			Str string
		}
		PP struct {
			Acc95  int `json:"95"`
			Acc96  int `json:"96"`
			Acc97  int `json:"97"`
			Acc98  int `json:"98"`
			Acc99  int `json:"99"`
			Acc100 int `json:"100"`
		} `json:"pp"`
	}
	Gameplay struct {
		LB struct {
			Player struct {
				Team int // Team 2 = Red, Team 1 = Blue, Team 0 = No Team
			} `json:"ourplayer"`
		} `json:"leaderboard"`
	}
}

func flatten(p Packet, f utils.FMap) {
	f.Set("skin", p.Settings.Folders.Skin)

	f.SetFunc("mode", func() string {
		// todo use a map
		switch p.Menu.GameMode {
		case 0:
			return "standard"
		case 1:
			return "taiko"
		case 2:
			return "catch"
		case 3:
			return "mania"
		default:
			return "unknown"
		}
	})
	f.Setf("mapid", "%d", p.Menu.BeatMap.Id)
	f.Setf("setid", "%d", p.Menu.BeatMap.Set)

	f.Set("artist", p.Menu.BeatMap.Metadata.Artist)
	f.Set("title", p.Menu.BeatMap.Metadata.Title)
	f.Set("diff", p.Menu.BeatMap.Metadata.Difficulty)
	f.Set("mapper", p.Menu.BeatMap.Metadata.Mapper)

	//unknown, unsubmitted, pending/wip/graveyard, unused, ranked, approved, qualified
	f.SetFunc("status", func() string {
		// todo use a map
		switch p.Menu.BeatMap.RankedStatus {
		case 0:
			return "unknown"
		case 1:
			return "unsubmitted"
		case 2:
			return "unranked"
		case 4:
			return "ranked"
		case 5:
			return "approved"
		case 6:
			return "qualified"
		default:
			return "gosupoop"
		}
	})

	f.SetFunc("bpm", func() string {
		low := p.Menu.BeatMap.Stats.BPM.Min
		high := p.Menu.BeatMap.Stats.BPM.Max

		if low == high {
			return fmt.Sprintf("%.2f", high)
		}

		return fmt.Sprintf("%.2f - %.2f", low, high)
	})

	f.SetFunc("url", func() string {
		bid := p.Menu.BeatMap.Id
		sid := p.Menu.BeatMap.Set
		mode := p.Menu.GameMode

		if bid > 0 {
			return fmt.Sprintf("https://osu.ppy.sh/b/%d?m=%d", bid, mode)
		} else if sid > 0 {
			return fmt.Sprintf("https://osu.ppy.sh/s/%d?m=%d", sid, mode)
		} else {
			return "no url available :/"
		}
	})

	f.Setf("fullname", "%s - %s [%s] (by %s)",
		p.Menu.BeatMap.Metadata.Artist,
		p.Menu.BeatMap.Metadata.Title,
		p.Menu.BeatMap.Metadata.Difficulty,
		p.Menu.BeatMap.Metadata.Mapper)

	f.Setf("mods", p.Menu.Mods.Str)
	f.Setf("pp100", "%d", p.Menu.PP.Acc100)
	f.Setf("pp99", "%d", p.Menu.PP.Acc99)
	f.Setf("pp98", "%d", p.Menu.PP.Acc98)
	f.Setf("pp97", "%d", p.Menu.PP.Acc97)
	f.Setf("pp96", "%d", p.Menu.PP.Acc96)
	f.Setf("pp95", "%d", p.Menu.PP.Acc95)

	f.SetFunc("team", func() string {
		switch p.Gameplay.LB.Player.Team {
		case 1:
			return "Blue"
		case 2:
			return "Red"
		default:
			return "No Team"
		}
	})

	f.Setf("ar", "%.1f", p.Menu.BeatMap.Stats.AR)
	f.Setf("cs", "%.1f", p.Menu.BeatMap.Stats.CS)
	f.Setf("od", "%.1f", p.Menu.BeatMap.Stats.OD)
	f.Setf("hp", "%.1f", p.Menu.BeatMap.Stats.HP)
	f.Setf("sr", "%.2f", p.Menu.BeatMap.Stats.FullSR)

	// todo the order doesnt make any fucking sense
}

func providerRunner(changes <-chan commandChange) {
	prov := setup.SelectedProvider
	for change := range changes {
		fmt.Printf("set %s to \"%s\"\n", change.name, change.content)
		prov.SetCommand(change.name, change.content)
		time.Sleep(1 * time.Second)
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
	var currentPacket Packet
	fmap := make(utils.FMap)

	formatTracker := make([]struct {
		t time.Time
		c string
		s bool
	}, len(setup.Config.Commands))

	changeWait := time.Duration(-1)
	first := true

	changes := make(chan commandChange, 100)
	go providerRunner(changes)

	c, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://%s/ws", setup.Config.Address),
		nil,
	)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	for {
		err := c.ReadJSON(&currentPacket)
		if err != nil {
			log.Println(err)
			// todo real error handling
			continue
		}

		flatten(currentPacket, fmap)

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
