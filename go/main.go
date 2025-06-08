package main

import "net/http"

func main() {
	http.HandleFunc("/game", startGame())

	http.ListenAndServe(":8085", nil)
}
