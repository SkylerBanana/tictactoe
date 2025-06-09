package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	matrix [3][3]string
	mu     sync.Mutex
)

func reader(conn *websocket.Conn, letter string) {
	defer func() {
		conn.Close()
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		handleMessage(p, letter)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func handleMessage(msg []byte, letter string) {
	var move []int

	err := json.Unmarshal(msg, &move)
	if err != nil {
		log.Println("Invalid move format:", err)
		return
	}
	log.Println(move)
	handleMove(letter, move)
}

func handleMove(letter string, move []int) {
	log.Println(letter)
	log.Println(matrix)

	row := move[0]
	column := move[1]
	if row < 0 || row > 2 || column < 0 || column > 2 {
		log.Println("Invalid move")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if matrix[row][column] == "" {
		matrix[row][column] = letter
	} else {
		log.Println("Space Occupied o.O")
	}
	log.Println(matrix[row][column])
	// prob gonna need a sync map or do redis PUB/SUB thats a thinker
	// redis PUB/SUB is more scalable but does scalability matter for this project

}

func startGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		letter := r.FormValue("letter")
		if letter != "X" && letter != "Y" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Client Connected to Websocket")

		reader(ws, letter)

	}

}
