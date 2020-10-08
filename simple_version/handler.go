package simple_version

import (
	"chatroom/util"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// NewMessageStreamHandler returns a streams endpoint handler.
func NewMessageStreamHandler(u *websocket.Upgrader) (http.HandlerFunc, error) {
	if broker == nil {
		return nil, errors.New("please init your broker with InitBroker")
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

		// reg := registration{id: util.RandID(15), conn: conn}
		reg := registration{id: r.RemoteAddr, conn: conn}
		broker.addConn <- reg

		defer func() {
			broker.removeConn <- reg
		}()

		util.LogInfo(fmt.Sprintf("client:%s has joined chatroom", reg.id))

		for {
			p := Payload{}
			err := conn.ReadJSON(&p)
			if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				util.LogInfo(fmt.Sprintf("client:%s has left the chatroom", reg.id))
				return
			}

			if err != nil {
				// This is an abnormal error, important to log it for debugging.
				util.LogErr("conn.ReadJSON", err)
				return
			}

			broker.broadcast <- p
		}
	}, nil
}
