package better_version

import (
	"chatroom/util"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func NewMessageStreamHandler(u *websocket.Upgrader) (http.HandlerFunc, error) {
	if broker == nil {
		return nil, errors.New("please run your broker with RunBroker")
	}

	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := u.Upgrade(w, r, nil)
		if err != nil {
			util.LogErr("Upgrade", err)
			return
		}

		defer conn.Close()

		client := newWebSocketClient(conn)
		broker.addClient <- client

		defer func() {
			broker.removeClient <- client
		}()

		util.LogInfo(fmt.Sprintf("client:%s has joined chatroom", client.ID()))
		client.Listen()
		util.LogInfo(fmt.Sprintf("client:%s has left chatroom", client.ID()))
	}, nil
}
