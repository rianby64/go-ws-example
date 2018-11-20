package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// is for serving WS
func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer (func() {
		for i, c := range conns {
			if c == conn {
				conns = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		conn.Close()
	})()

	if err != nil {
		log.Println("connecting", err)
		return
	}

	conns = append(conns, conn)
	for {
		var request JSONRPCrequest
		if err := conn.ReadJSON(&request); err != nil {
			log.Println(err)
			return
		}

		if err := processRequest(&request, conn); err != nil {
			log.Println(err)
			return
		}
	}
}
