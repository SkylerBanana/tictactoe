package main

import (
	"net/http"
	"tictactoe/go/redis"
)

func main() {

	client := redis.InitRedis()
	store := redis.NewRedisInstance(client)
	http.HandleFunc("/game", startGame(store))

	http.ListenAndServe(":8085", nil)
}
