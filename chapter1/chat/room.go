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
