package main

import (
	"fmt"
	"log"
	"net/http"
)

func printConnectionInfo() {
	log.Println("=== Thông tin kết nối ===")

	log.Println("-- BroadCast Channels --")
	for key, channel := range broadCast {
		log.Printf("Channel %s: %v", key, channel)
	}

	log.Println("-- Connected Clients --")
	clientsMutex.Lock()
	for client, event := range clients {
		log.Printf("Client: %v - Event: %s", client.RemoteAddr(), event)
	}
	clientsMutex.Unlock()

	log.Println("=====================")
}

func handle_websocket(w http.ResponseWriter, r *http.Request) {
	conn := connect_client(w, r)
	if conn == nil {
		return
	}

	handle_connection("test", conn)
}

func main() {
	http.HandleFunc("/ws/test", handle_websocket)

	fmt.Println("Server đang chạy tại :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
