package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var broadCast = make(map[string]chan bool) // lưu trữ tín hiệu

// NOTE: broadCast[nameEvent] = make(chan bool, 1) 
// Đảm bảo phải được khai báo đầu tiên nếu muốn sử dụng nếu không sẽ gặp bug không thể thêm tín hiệu

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

	clients      = make(map[*websocket.Conn][]string) // các client đang kết nối
	clientsMutex = &sync.Mutex{} // tạo ra để xử lý vài trường hợp từ goroutine
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
		for _, event := range nameEvent {
			delete(broadCast, event)
		}
		return false
	}
	return true
}

// Tạo tín hiệu 
func create_broadcast(nameEvent string) {
	if broadCast[nameEvent] == nil {
		broadCast[nameEvent] = make(chan bool, 1)
	}
	broadCast[nameEvent] <- true
}

/// kiểm tra client còn sống hay không
func check_client_alive(conn *websocket.Conn) {
	go func() {
		recive_message_client(conn)
	}()
}

// xử lý mỗi khi có tín hiệu
func handle_Websocket(nameEvent string) {
	go func() {
		for <-broadCast[nameEvent] {
			clientsMutex.Lock()
			for conn := range clients {
				send_message_client(conn, "Dữ liệu từ server")
			}
			clientsMutex.Unlock()
		}
	}()
}


// DEMO xử lý
func handle_connection(nameEvent string, conn *websocket.Conn) {
	clients[conn] = append(clients[conn], nameEvent)

	create_broadcast(nameEvent)

	handle_Websocket(nameEvent)
	check_client_alive(conn)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
