package main

import (
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	// socketはこのクライアントのためのwebsocket
	socket *websocket.Conn

	// メッセージが送られるチャネル
	send chan *message

	// roomはこのクライアントが参加しているチャットルーム
	room *room

	// userDataはユーザーに関する情報を保持
	userData map[string]interface{}
}

// websocketからメッセージを読み込み、roomのforwardチャネルへ送信
func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// sendチャネルからメッセージを受け取り、WriteMessageに書き出す
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
