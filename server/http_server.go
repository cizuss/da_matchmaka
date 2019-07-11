package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type EnqueuePlayerRequest struct {
	Name  string
	Skill int
}

func (r EnqueuePlayerRequest) toPlayer() Player {
	return NewPlayer(r.Name, r.Skill)
}

func handleEnqueuePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req EnqueuePlayerRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			panic(err)
		}
		player := req.toPlayer()
		addPlayerToQueue(&player)
	}
}

func startHttpServer() {
	http.HandleFunc("/enqueue", handleEnqueuePlayer)
	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
