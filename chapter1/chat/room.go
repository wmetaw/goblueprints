package main

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
