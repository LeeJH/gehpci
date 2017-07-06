package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"

	"gehpci/models"
	//	"samples/WebIM/models"
)

type WebSocketController struct {
	//	baseController
	BaseController
}

// @Title Bind to a web socket bash
// @Description Bind to a web socket bash
// @Param       machine            path    string  true            "The machine name"
// @Success 200 {common.Result} result
// @Failure 403 {err} body is err info
// @Failure 401 need login
// @Failure 400 Not a websocket handshake
// @router /bindbash/:machine [get]
func (this *WebSocketController) BindBash() {
	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	wsm := models.NewWebSocketMan(ws, this.user)
	wsm.BindBash()
	// Join chat room.
	//Join(uname, ws)
	//defer Leave(uname)
	//Bind()

	// Message receive loop.
	// for {
	// 	t, p, err := ws.ReadMessage()
	// 	if err != nil {
	// 		return
	// 	}
	// 	log.Printf("t p e : \n%#v ; %#v ; %#v\n", t, p, err)
	// 	//publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
	// }
}

var (
// Channel for new join users.
//	subscribe = make(chan Subscriber, 10)
// Channel for exit users.
//	unsubscribe = make(chan string, 10)
// Send events here to publish them.
//	publish = make(chan models.Event, 10)
// Long polling waiting list.
//	waitingList = list.New()
//	subscribers = list.New()
)
