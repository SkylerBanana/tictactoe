package que

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"tictactoe/go/redis"
	"time"
)

type Player struct {
	UserName string `json:"userName"`
	UserID   string `json:"userID"`
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
	return nil
}
