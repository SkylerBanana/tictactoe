package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"tictactoe/go/redis"

	"github.com/google/uuid"

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
	Letter string   `json:"letter"`
	Matrix []string `json:"matrix"`
}

func reader(conn *websocket.Conn, letter string, ctx context.Context, store redis.Store, uuid string) {
	defer func() {
		conn.Close()
	}()
	go subscribeToChannel(ctx, store, uuid, conn)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("bingus")

		handleMove(letter, conn, store, ctx, uuid, p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func handleMove(letter string, conn *websocket.Conn, store redis.Store, ctx context.Context, uuid string, msg []byte) {

	var move []int

	err := json.Unmarshal(msg, &move)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte("Invalid move format"))
		return
	}

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

	// what if just flatten the array and store it in redis that way the go server dosent have to store the game

	store.Publish(ctx, uuid, data)
	println("After Publish")

}

func startGame(store redis.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		uuid := uuid.NewString()

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

		reader(ws, letter, r.Context(), store, uuid)

	}

}

func playGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func joinGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// this function will subscribe to the redis channel based off of the UUID

	}
}
