package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/memsdm05/nplink/utils"
)

// Packet Not the full packet as some information is unneeded
type gosumemoryPacket struct {
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

func (g *gosumemoryPacket) check(data []byte) bool {
	json.Unmarshal(data, g)
	return *g != gosumemoryPacket{}
}

func (g *gosumemoryPacket) fill(conn *websocket.Conn) error {
	return conn.ReadJSON(g)
}

func (g *gosumemoryPacket) flatten(f utils.FMap) {
	f.Set("skin", g.Settings.Folders.Skin)

	f.SetFunc("mode", func() string {
		// todo use a map
		switch g.Menu.GameMode {
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
	f.Setf("mapid", "%d", g.Menu.BeatMap.Id)
	f.Setf("setid", "%d", g.Menu.BeatMap.Set)

	f.Set("artist", g.Menu.BeatMap.Metadata.Artist)
	f.Set("title", g.Menu.BeatMap.Metadata.Title)
	f.Set("diff", g.Menu.BeatMap.Metadata.Difficulty)
	f.Set("mapper", g.Menu.BeatMap.Metadata.Mapper)

	//unknown, unsubmitted, pending/wip/graveyard, unused, ranked, approved, qualified
	f.SetFunc("status", func() string {
		// todo use a map
		switch g.Menu.BeatMap.RankedStatus {
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
		low := g.Menu.BeatMap.Stats.BPM.Min
		high := g.Menu.BeatMap.Stats.BPM.Max

		if low == high {
			return fmt.Sprintf("%.2f", high)
		}

		return fmt.Sprintf("%.2f - %.2f", low, high)
	})

	f.SetFunc("url", func() string {
		bid := g.Menu.BeatMap.Id
		sid := g.Menu.BeatMap.Set
		mode := g.Menu.GameMode

		if bid > 0 {
			return fmt.Sprintf("https://osu.ppy.sh/b/%d?m=%d", bid, mode)
		} else if sid > 0 {
			return fmt.Sprintf("https://osu.ppy.sh/s/%d?m=%d", sid, mode)
		} else {
			return "no url available :/"
		}
	})

	f.Setf("fullname", "%s - %s [%s] (by %s)",
		g.Menu.BeatMap.Metadata.Artist,
		g.Menu.BeatMap.Metadata.Title,
		g.Menu.BeatMap.Metadata.Difficulty,
		g.Menu.BeatMap.Metadata.Mapper)

	f.Setf("mods", g.Menu.Mods.Str)
	f.Setf("pp100", "%d", g.Menu.PP.Acc100)
	f.Setf("pp99", "%d", g.Menu.PP.Acc99)
	f.Setf("pp98", "%d", g.Menu.PP.Acc98)
	f.Setf("pp97", "%d", g.Menu.PP.Acc97)
	f.Setf("pp96", "%d", g.Menu.PP.Acc96)
	f.Setf("pp95", "%d", g.Menu.PP.Acc95)

	f.SetFunc("team", func() string {
		switch g.Gameplay.LB.Player.Team {
		case 1:
			return "Blue"
		case 2:
			return "Red"
		default:
			return "No Team"
		}
	})

	f.Setf("ar", "%.1f", g.Menu.BeatMap.Stats.AR)
	f.Setf("cs", "%.1f", g.Menu.BeatMap.Stats.CS)
	f.Setf("od", "%.1f", g.Menu.BeatMap.Stats.OD)
	f.Setf("hp", "%.1f", g.Menu.BeatMap.Stats.HP)
	f.Setf("sr", "%.2f", g.Menu.BeatMap.Stats.FullSR)

	// todo the order doesnt make any fucking sense
}
