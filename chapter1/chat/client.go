package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	// socketはこのクライアントのためのwebsocket
	socket *websocket.Conn

	// メッセージが送られるチャネル
	send chan []byte

	// roomはこのクライアントが参加しているチャットルーム
	room *room
}

// websocketからメッセージを読み込み、roomのforwardチャネルへ送信
func (c *Client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

// sendチャネルからメッセージを受け取り、WriteMessageに書き出す
func (c *Client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
