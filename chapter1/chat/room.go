package main

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"net/http"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte

	// joinは参加しようとしているクライアントのためのチャネル
	join chan *client

	// leaveはチャットルームから退出しようとしているチャネル
	leave chan *client

	// clientsには在室している全てのクライアントが保持される
	clients map[*client]bool
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func (r *room) run() {
	for {
		select {
		// 参加
		case client := <-r.join:
			r.clients[client] = true

			// 退室
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)

			// 全てのクライアントにメッセージを転送
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
