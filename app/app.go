package app

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/setup"
	"github.com/memsdm05/nplink/util"
	"log"
)

type commandChange struct {
	name    string
	content string
}

var fmap = make(util.FMap)

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
		GameMode int
		Name     string
		Score    int
		Accuracy float64
		Combo    struct {
			Current int
			Max     int
		}
		Hits struct {
			Val0   int `json:"0"`
			Val50  int `json:"50"`
			Val100 int `json:"100"`
			Val200 int `json:"200"`
			Val300 int `json:"300"`
			Geki   int
			Katu   int
			SB     int `json:"sliderBreaks"`
			Grade  struct {
				Current string
				Max     string `json:"maxThisPlay"`
			}
			UR float64 `json:"unstableRate"`
		}
		PP struct {
			Current int
			FC      int `json:"fc"`
			Max     int `json:"maxThisPlay"`
		} `json:"pp"`
		LB struct {
			Player struct {
				Team     int // Team 2 = Red, Team 1 = Blue, Team 0 = No Team
				Position int
			} `json:"ourplayer"`
		} `json:"leaderboard"`
	}
}

func flatten(p Packet, f util.FMap) {
	f.Set("skin", p.Settings.Folders.Skin)
	//f.Setf("map", "%d", p.Menu.BeatMap.Id)
	//f.Setf("set", "%d", p.Menu.BeatMap.Set)

	f.Set("artist", p.Menu.BeatMap.Metadata.Artist)
	f.Set("title", p.Menu.BeatMap.Metadata.Title)
	f.Set("diff", p.Menu.BeatMap.Metadata.Difficulty)
	f.Set("mapper", p.Menu.BeatMap.Metadata.Mapper)
	f.SetFunc("url", func() string {
		bid := p.Menu.BeatMap.Id
		sid := p.Menu.BeatMap.Set
		mode := p.Menu.GameMode

		if bid > 0 {
			return fmt.Sprintf("https://osu.ppy.sh/b/%d?m=%d", bid, mode)
		} else if sid > 0{
			return fmt.Sprintf("https://osu.ppy.sh/s/%d?m=%d", sid, mode)
		} else {
			return "no url available :/"
		}
	})
	//f.Setf("fullname", "%s - %s [%s] (by %s)",
	//	p.Menu.BeatMap.Metadata.Artist,
	//	p.Menu.BeatMap.Metadata.Title,
	//	p.Menu.BeatMap.Metadata.Difficulty,
	//	p.Menu.BeatMap.Metadata.Mapper,
	//	)
	f.Setf("mods", p.Menu.Mods.Str)

	f.Setf("pp100", "%d", p.Menu.PP.Acc100)
	f.Setf("pp99" , "%d", p.Menu.PP.Acc99)
	f.Setf("pp98" , "%d", p.Menu.PP.Acc98)
	f.Setf("pp97" , "%d", p.Menu.PP.Acc97)
	f.Setf("pp96" , "%d", p.Menu.PP.Acc96)
	f.Setf("pp95" , "%d", p.Menu.PP.Acc95)

	f.Setf("pp", "%d", p.Gameplay.PP.Current)
	f.Setf("ppmax", "%d", p.Gameplay.PP.Max)
	f.Setf("ppfc", "%d", p.Gameplay.PP.FC)

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

	f.Setf("300s", "%d",   p.Gameplay.Hits.Val300)
	f.Setf("200s", "%d",   p.Gameplay.Hits.Val200)
	f.Setf("100s", "%d",   p.Gameplay.Hits.Val100)
	f.Setf("50s", "%d",    p.Gameplay.Hits.Val50)
	f.Setf("misses", "%d", p.Gameplay.Hits.Val0)
	f.Setf("sliderbreaks", "%d", p.Gameplay.Hits.SB)

}

func providerRunner(provider provider.Provider, changes <-chan commandChange) {

}

func processRunner(raw <-chan string) {

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
			return
		}



		flatten(currentPacket, fmap)

		for key, value := range fmap {
			fmt.Println(key, value)
		}
		fmt.Println()
	}
}
