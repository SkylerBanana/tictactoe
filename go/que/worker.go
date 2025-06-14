package que

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"tictactoe/go/redis"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	UserName string `json:"userName"`
	UserID   string `json:"userID"`
}

type Message struct {
	MatchId string `json:"matchID"`
	Status  string `json:"status"`
}

// to remove someone from que i can just do a defer in the reader of the websocket i think
func ProcessQue(store redis.Store, ctx context.Context) error {

	length, err := store.Length(ctx, "MatchMaking")
	if err != nil {
		return err
	}

	if length < 2 {
		log.Println("Not enough players to start a match")

		return errors.New("Not enough players to start a match")
	}

	var players []Player

	for i := 0; i < 2; i++ {
		result, err := store.QuePop(ctx, 1*time.Second)
		if err != nil || len(result) < 2 {
			return err
		}

		var player Player
		if err := json.Unmarshal([]byte(result[1]), &player); err != nil {
			return err
		}

		players = append(players, player)
	}

	log.Printf("Matched players: %+v vs %+v\n", players[0], players[1])

	err = CreateMatch(store, players[0], players[1], ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateMatch(store redis.Store, player1 Player, player2 Player, ctx context.Context) error {
	matchid := uuid.NewString()

	var message []Message

	message = append(message, Message{
		MatchId: matchid,
		Status:  "matched",
	})

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := store.Publish(ctx, player1.UserID, data); err != nil {
		log.Printf("Failed to publish to %s: %v", player1.UserID, err)
		return err
	}

	if err := store.Publish(ctx, player2.UserID, data); err != nil {
		log.Printf("Failed to publish to %s: %v", player2.UserID, err)
		return err
	}

	return nil
}
