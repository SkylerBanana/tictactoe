package main

import (
	"context"
	"encoding/json"
	"log"
	"tictactoe/go/que"
	"tictactoe/go/redis"

	"github.com/gorilla/websocket"
)

func subscribeToChannel(ctx context.Context, store redis.Store, uuid string, conn *websocket.Conn) {

	store.Subscribe(ctx, uuid, func(message string) {
		msgJSON, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshalling message to JSON: %v", err)
			return
		}

		// when i receieve match found ill subscribe to that instead

		// This sends notifications with websocket back to the user
		writeBackToClient(msgJSON, conn)
	})

}

func writeBackToClient(msgJSON []byte, conn *websocket.Conn) {

	_ = conn.WriteMessage(websocket.TextMessage, msgJSON)

}

func reader(conn *websocket.Conn, ctx context.Context, store redis.Store, claims *CustomClaims) {
	defer func() {
		conn.Close()
		//store.Deque(ctx,claims)
		// in here ill remove the user from que if they disconnect
	}()

	go subscribeToChannel(ctx, store, claims.UserId, conn)

	if err := quePlayer(store, claims, ctx); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to join que"))
		return
	}

	if que.ProcessQue(store, ctx) != nil {
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
