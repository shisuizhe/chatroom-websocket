package better_version

import (
	"chatroom/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type Client interface {
	ID() string
	Listen()
	Activate(ctx context.Context)
	WriteQueue() chan<- Payload
	SetBroadcast(ch chan Payload)
}

// WebSocketClient reads and writes to broker.
type WebsocketClient struct {
	id        string
	conn      *websocket.Conn
	read      chan json.RawMessage // read queue
	write     chan Payload         // write queue
	broadcast chan Payload
}

// newWebsocketClient
func newWebSocketClient(conn *websocket.Conn) Client {
	return &WebsocketClient{
		id:    util.RandID(8),
		conn:  conn,
		read:  make(chan json.RawMessage, 100),
		write: make(chan Payload, 100),
	}
}

// ID 获取 client's ID
func (client *WebsocketClient) ID() string {
	return client.id
}

// Listen 不断监听websocket连接，将接收到的消息放入WebsocketClient.read
func (client *WebsocketClient) Listen() {
	for {
		_, p, err := client.conn.ReadMessage()
		// 错误为客户端正常关闭或离开
		if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return
		}

		// 不正常的错误，记录日志
		if err != nil {
			util.LogErr("ReadMessage", err)
			return
		}

		client.read <- p
	}
}

// Activate
func (client *WebsocketClient) Activate(ctx context.Context) {
	go client.readLoop(ctx)
	go client.writeLoop(ctx)
}

// WriteQueue 返回写通道
func (client *WebsocketClient) WriteQueue() chan<- Payload {
	return client.write
}

// SetBroadcast 设置广播通道
func (client *WebsocketClient) SetBroadcast(ch chan Payload) {
	client.broadcast = ch
}

// readLoop 从读取队列中抽出并通过将消息分解为结构来处理消息
func (client *WebsocketClient) readLoop(ctx context.Context) {
	client.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	client.conn.SetPongHandler(func(appData string) error {
		client.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			util.LogInfo(fmt.Sprintf("client:%s has termianted read loop", client.ID()))
			return

		case p := <-client.read:
			client.conn.SetReadDeadline(time.Now().Add(2 * time.Second))

			payload := Payload{}
			err := json.Unmarshal(p, &payload)
			if err != nil {
				util.LogErr("Unmarshal", err)
				continue
			}

			client.broadcast <- payload
		}
	}
}

// writeLoop 从写入通道中抽出并通过将消息封送为字节来处理消息，并将字节写入websocket连接
func (client *WebsocketClient) writeLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ctx.Done():
			util.LogInfo(fmt.Sprintf("client %s has terminated write loop", client.id))
			return

		case payload := <-client.write:
			client.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

			pByte, err := json.Marshal(payload)
			if err != nil {
				util.LogErr("Marshal", err)
				continue
			}

			err = client.conn.WriteMessage(websocket.TextMessage, pByte)
			if err != nil {
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

			err := client.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}
