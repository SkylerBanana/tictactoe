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

type gameState struct {
	Type   string   `json:"type"`
	Matrix []string `json:"matrix"`
}

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

		handleMessage(p, letter, conn)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func handleMessage(msg []byte, letter string, conn *websocket.Conn) {
	var move []int

	err := json.Unmarshal(msg, &move)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte("Invalid move format"))
		return
	}
	log.Println(move)
	handleMove(letter, move, conn)
}

func handleMove(letter string, move []int, conn *websocket.Conn) {
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
		_ = conn.WriteMessage(websocket.TextMessage, []byte("Space Occupied o.O"))
		return
	}

	// have to flatten the matrix so we can send it as json
	var flatMatrix []string
	for _, row := range matrix {
		flatMatrix = append(flatMatrix, row[:]...)
	}

	var gamestate []gameState
	gamestate = append(gamestate, gameState{
		Type:   "Game Update",
		Matrix: flatMatrix,
	})

	data, err := json.Marshal(gamestate)
	if err != nil {
		log.Println("Error Marshalling to Json")
	}

	conn.WriteMessage(websocket.TextMessage, data)

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
