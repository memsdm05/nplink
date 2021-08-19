package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/memsdm05/nplink/internal/utils"
	"strconv"
	"strings"
)

type streamCompanionPacket struct{
	hasRequested  bool

	Md5 string
	Skin string
	Acc100 float32 `json:"osu_SSPP"`
	Acc99 float32 `json:"osu_99PP"`
	Acc98 float32 `json:"osu_98PP"`
	Acc97 float32 `json:"osu_97PP"`
	Acc96 float32 `json:"osu_96PP"`
	Acc95 float32 `json:"osu_95PP"`

	Sid int `json:"mapid"`
	Bid int `json:"mapsetid"`

	AR float32 `json:"ar"`
	CS float32 `json:"cs"`
	HP float32 `json:"hp"`
	OD float32 `json:"od"`
	SR float32 `json:"starsNomod"`

	Mods string

	MinBpm float32
	MaxBpm float32
	MainBpm float32

	Artist string `json:"artistRoman"`
	Title string `json:"titleRoman"`
	Diff string `json:"diffName"`
	Mapper string `json:"creator"`

	Mode string
	RankedStatus int
}

var watchTokens = []string {
	"md5", "skin",
	"osu_SSPP", "osu_99PP", "osu_98PP", "osu_97PP", "osu_96PP", "osu_95PP",
	"mapid", "mapsetid",
	"ar", "cs", "hp", "od", "starsNomod",
	"mods",
	"minBpm",
	"maxBpm",
	"artistRoman", "titleRoman", "diffName", "creater",
	"mode",
	"rankedStatus",
}

func (s *streamCompanionPacket) check(b []byte) bool {
	json.Unmarshal(b, s)
	return true
}

func (s *streamCompanionPacket) fill(conn *websocket.Conn) error {
	if !s.hasRequested {
		err := conn.WriteJSON(watchTokens)

		if err != nil {
			return err
		}

		s.hasRequested = true
	}

	return conn.ReadJSON(s)
}

func (s *streamCompanionPacket) flatten(f utils.FMap) {
	f.Set("skin", s.Skin)

	f.SetFunc("mode", func() string {
		// todo use a map
		m, _ := strconv.Atoi(s.Mode)
		switch m {
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

	//f.Set("mode", s.Mode)

	f.Setf("mapid", "%d", s.Bid)
	f.Setf("setid", "%d", s.Sid)

	f.Set("artist", s.Artist)
	f.Set("title", s.Title)
	f.Set("diff", s.Diff)
	f.Set("mapper", s.Mapper)

	f.SetFunc("status", func() string {
		// todo use a map
		switch s.RankedStatus {
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
			return ".net sucks haha"
		}
	})

	f.SetFunc("bpm", func() string {
		low := s.MinBpm
		high := s.MaxBpm

		if low == high {
			return fmt.Sprintf("%.2f", s.MainBpm)
		}

		return fmt.Sprintf("%.2f - %.2f", low, high)
	})

	f.SetFunc("url", func() string {
		bid := s.Bid
		sid := s.Sid
		mode, _ := strconv.Atoi(s.Mode)

		if bid > 0 {
			return fmt.Sprintf("https://osu.ppy.sh/b/%d?m=%d", bid, mode)
		} else if sid > 0 {
			return fmt.Sprintf("https://osu.ppy.sh/s/%d?m=%d", sid, mode)
		} else {
			return "no url available :/"
		}
	})

	f.Setf("fullname", "%s - %s [%s] (by %s)",
		s.Artist,
		s.Title,
		s.Diff,
		s.Mapper)

	// note: mods are not in correct order
	f.SetFunc("mods", func() string {
		return strings.Join(strings.Split(s.Mods, ","), "")
	})
	f.Setf("pp100", "%.0f", s.Acc100)
	f.Setf("pp99", "%.0f", s.Acc99)
	f.Setf("pp98", "%.0f", s.Acc98)
	f.Setf("pp97", "%.0f", s.Acc97)
	f.Setf("pp96", "%.0f", s.Acc96)
	f.Setf("pp95", "%.0f", s.Acc95)

	f.Setf("ar", "%.1f", s.AR)
	f.Setf("cs", "%.1f", s.CS)
	f.Setf("od", "%.1f", s.OD)
	f.Setf("hp", "%.1f", s.HP)
	f.Setf("sr", "%.2f", s.SR)
}

func (s *streamCompanionPacket) path() string {
	return "tokens"
}
