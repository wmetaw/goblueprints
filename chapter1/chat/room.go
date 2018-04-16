package main

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/objx"
	"github.com/wmetaw/goblueprints/chapter1/trace"
	"net/http"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan *message

	// joinは参加しようとしているクライアントのためのチャネル
	join chan *client

	// leaveはチャットルームから退出しようとしているチャネル
	leave chan *client

	// clientsには在室している全てのクライアントが保持される
	clients map[*client]bool

	// tracerはチャットルームで行われた操作のログを受け取る
	tracer trace.Tracer

	// avatar情報
	avatar Avatar
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました:", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value), // cookieデータをデコードしmapへ変換
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
			r.tracer.Trace("新しいクライアントが参加しました")

			// 退室
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")

			// 受信
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました:", msg.Message)

			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {

				// メッセージ送信
				case client.send <- msg:
					r.tracer.Trace(" -- クライアントに送信されました")

					// 失敗
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}
	}
}
