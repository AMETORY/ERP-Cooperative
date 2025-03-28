package services

import (
	"net/http"

	"gopkg.in/olahol/melody.v1"
)

var WS = &melody.Melody{}

func InitWS() *melody.Melody {
	mel := melody.New()
	mel.Config.MaxMessageSize = 2000
	mel.Config.MaxMessageSize = 2000

	mel.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	WS = mel
	return mel
}
