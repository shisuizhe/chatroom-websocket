package simple_version

import (
	"chatroom/util"
	"github.com/gorilla/websocket"
)

var broker *MessageBroker

// InitBroker 初始化消息代理并运行它
func InitBroker() {
	broker = &MessageBroker{
		connByID:   make(map[string]*websocket.Conn),
		addConn:    make(chan registration),
		removeConn: make(chan registration),
		broadcast:  make(chan Payload),
	}

	go broker.loop()
}

// registration 发送给代理以注册新连接的Payload
type registration struct {
	id   string
	conn *websocket.Conn
}

// Payload 当用户输入用户名和消息时，服务器收到Payload则将其发送给消息代理，
// 然后消息代理就可以把这个Payload广播给所有其他用户
type Payload struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// MessageBroker 每当收到消息时，消息代理都会接收此消息，并通过侦听此消息将其发送到订单连接
type MessageBroker struct {
	connByID   map[string]*websocket.Conn
	addConn    chan registration
	removeConn chan registration
	broadcast  chan Payload
}

// loop 收听广播并将其写入每个客户端
func (mb *MessageBroker) loop() {
	for {
		select {
		// 收到广播
		case payload := <-mb.broadcast:
			for ID := range mb.connByID {
				err := mb.connByID[ID].WriteJSON(payload) // This is bad practice
				if err != nil {
					util.LogErr("WriteJSON", err)
				}
			}
		// 添加连接（客户端连接上时）
		case r := <-mb.addConn:
			mb.connByID[r.id] = r.conn
		// 移除连接（客户端与服务器断开连接时）
		case r := <-mb.removeConn:
			delete(mb.connByID, r.id)
		}
	}
}
