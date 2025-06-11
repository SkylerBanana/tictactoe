package que

import (
	"context"
	"tictactoe/go/redis"
)

type Player struct {
	UserName string `json:"userName"`
	UserID   string `json:"userID"`
}

func ProcessQue(store *redis.Store) {
	context.Background()

}
