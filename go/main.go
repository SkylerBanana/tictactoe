package main

import (
	"net/http"
	"tictactoe/go/auth"

	"tictactoe/go/redis"
)

func main() {

	client := redis.InitRedis()
	store := redis.NewRedisInstance(client)

	http.HandleFunc("/game", startGame(store))

	http.HandleFunc("/checkuser", auth.IsLoggedIn())
	http.HandleFunc("/login", auth.Login())

	http.ListenAndServe(":8085", nil)
}
