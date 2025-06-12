package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"tictactoe/go/redis"

	"github.com/golang-jwt/jwt/v5"

	"tictactoe/go/que"

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

type CustomClaims struct {
	UserName string `json:"UserName"`
	UserId   string `json:"UserId"`
	jwt.RegisteredClaims
}

type QueuedPlayer struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}

func reader(conn *websocket.Conn, ctx context.Context, store redis.Store, claims *CustomClaims) {
	defer func() {
		conn.Close()
		// in here ill remove the user from que if they disconnect
	}()
	go subscribeToChannel(ctx, store, claims.UserId, conn)
	if que.ProcessQue(store, ctx) != nil {
		//Tbh this could be a completely different message  entirely
		conn.WriteMessage(websocket.TextMessage, []byte("Not Enough Players In Que"))
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		//handleMove(letter, conn, store, ctx, uuid, p)

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
		secret := []byte(os.Getenv("JWT_SECRET"))
		auth_token, _ := r.Cookie("auth_token")

		claims := &CustomClaims{}

		_, err := jwt.ParseWithClaims(auth_token.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			log.Printf("JWT parse error: %v\n", err)
			http.Error(w, "Invalid Token", http.StatusUnauthorized)

			return
		}

		if err := quePlayer(store, claims, r.Context()); err != nil {
			http.Error(w, "Failed to queue player", http.StatusInternalServerError)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Client Connected to Websocket")

		reader(ws, r.Context(), store, claims)

	}

}

func quePlayer(store redis.Store, claims *CustomClaims, ctx context.Context) error {

	player := QueuedPlayer{
		UserId:   claims.UserId,
		UserName: claims.UserName,
	}

	data, err := json.Marshal(player)
	if err != nil {
		log.Println("Failed to marshal to JSON:", err)
		return err
	}

	err = store.Que(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
