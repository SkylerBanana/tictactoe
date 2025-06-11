package main

import (
	"context"
	"encoding/json"
	"log"
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

		// This sends notifications with websocket back to the user
		writeBackToClient(msgJSON, conn)
	})

}

func writeBackToClient(msgJSON []byte, conn *websocket.Conn) {

	_ = conn.WriteMessage(websocket.TextMessage, msgJSON)

}
