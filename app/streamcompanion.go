package app

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/memsdm05/nplink/utils"
)

type streamCompanionPacket struct {

}

func (s *streamCompanionPacket) check(b []byte) bool {
	json.Unmarshal(b, s)
	return s == new(streamCompanionPacket)
}

func (s *streamCompanionPacket) fill(conn *websocket.Conn) error {
	return conn.ReadJSON(s)
}

func (s *streamCompanionPacket) flatten(fmap utils.FMap) {
	fmap.Set("test", "test")
}


