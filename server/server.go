package server

import (
	"fmt"
	"log"
	"net/http"
	"remy/events"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run() {
	http.HandleFunc("/ws", handleWebSocket)
	serveStaticFiles()

	err := http.ListenAndServeTLS("0.0.0.0:8888", "server.crt", "server.key", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func serveStaticFiles() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
}

func handleWebSocket(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected!")

	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)

			sendResponse(conn, false)
			continue
		}

		err = events.HandleEvent(rawMessage)

		if err != nil {
			sendResponse(conn, false)
			continue
		}

		sendResponse(conn, true)
	}
}

func sendResponse(conn *websocket.Conn, success bool) {
	response := map[string]bool{
		"success": success,
	}
	conn.WriteJSON(response)
}
