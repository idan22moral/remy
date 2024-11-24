package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"remy/events"
	"remy/static"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run(addr string) error {
	http.HandleFunc("/ws", handleWebSocket)
	serveStaticFiles()

	rootCert, rootKey := generateCACertificate()
	cert, _ := generateTLSCredentials(rootCert, rootKey)

	server := http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	err := server.ListenAndServeTLS("", "")
	if err != nil {
		return err
	}
	return err
}

func serveStaticFiles() {
	fs := http.FileServer(http.FS(static.StaticFolder))
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
