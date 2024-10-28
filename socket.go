package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var broadCast = make(map[string]chan bool)

// Socket connection
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func get_event_socket(fileNameSocket string, nameEvent string) {
	file, err := os.Open(fileNameSocket)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var context string

	for scanner.Scan() {
		context = scanner.Text()
	}
	if context == "SUCCESS" {
		broadCast[nameEvent] <- true
	} else {
		log.Printf("Lỗi từ công cụ %s", context)
	}
}

// Websocket connection
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients      = make(map[*websocket.Conn]string)
	clientsMutex = &sync.Mutex{}
)

func connect_client(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return nil
	}
	return conn
}

func recive_message_client(conn *websocket.Conn) string {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Client ngắt kết nối đột ngột:", conn.RemoteAddr())
			conn.Close()
			delete(clients, conn)
			return
		}
	}()

	_, messageClient, err := conn.ReadMessage()
	if err != nil {
		log.Println("Client mất kết nối tại:", conn.RemoteAddr())
		conn.Close()
		delete(clients, conn)
		return ""
	}

	log.Printf("Dữ liệu client %s là: %s", conn.RemoteAddr(), string(messageClient))

	return string(messageClient)
}

func send_message_client(conn *websocket.Conn, message interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Client ngắt kết nối đột ngột:", conn.RemoteAddr())
			conn.Close()
			delete(clients, conn)
			return
		}
	}()

	nameEvent := clients[conn]

	err := conn.WriteJSON(message)

	if err != nil {
		log.Println("Client đã mất kết nối tại địa chỉ:", conn.RemoteAddr())
		conn.Close()
		delete(clients, conn)
		delete(broadCast, nameEvent)
		return false
	}
	return true
}

func handle_connection(nameEvent string, conn *websocket.Conn) {
	clients[conn] = nameEvent
	broadCast[nameEvent] = make(chan bool, 1)

	for {
		printConnectionInfo() // kiểm tra dữ liệu

		message := recive_message_client(conn)
		if message == "" { 
			break
		}

		check := send_message_client(conn, message)
		if !check {
			break
		}

		log.Printf("Đã gửi dữ liệu cho client")
	}

	delete(clients, conn)
	delete(broadCast, nameEvent)
	conn.Close()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
