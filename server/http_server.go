package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type EnqueuePlayerRequest struct {
	Name  string
	Skill int
}

type CreateWebsocketConnectionRequest struct {
	Name string
}

type MHttpServer struct {
	jobsChan       chan *Player
	partiesChan    chan *Party
	connectionsMap map[string]*websocket.Conn
	connections    map[*websocket.Conn]bool
	upgrader       websocket.Upgrader
}

func (r EnqueuePlayerRequest) toPlayer() Player {
	return NewPlayer(r.Name, r.Skill)
}

func (server MHttpServer) handleConnectionMessage(w http.ResponseWriter, r *http.Request) {
	c, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer c.Close()
	server.connections[c] = true
	for {
		var msg CreateWebsocketConnectionRequest
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Fatal(err)
			delete(server.connections, c)
			continue
		}
		fmt.Println("Adding connection for player ", msg.Name)
		server.connectionsMap[msg.Name] = c
	}
}

func (server MHttpServer) handlePartyMessages() {
	for {
		party := <-server.partiesChan
		for _, player := range party.players {
			playerName := player.name
			conn := server.connectionsMap[playerName]
			if conn != nil {
				msg := fmt.Sprintf("found party : %v", party)
				err := conn.WriteJSON(msg)
				if err != nil {
					fmt.Println(err)
					conn.Close()
					delete(server.connectionsMap, playerName)
				}
			}
		}
	}
}

func (server MHttpServer) handleEnqueuePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req EnqueuePlayerRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			fmt.Println("Error in handling enqueue player ", err)
			return
		}
		player := req.toPlayer()
		server.addPlayerToQueue(&player)
	}
}

func (server MHttpServer) addPlayerToQueue(p *Player) {
	server.jobsChan <- p
}

func (server MHttpServer) start() {
	http.HandleFunc("/enqueue", server.handleEnqueuePlayer)
	http.HandleFunc("/ws", server.handleConnectionMessage)
	// start a goroutine that handles found parties and sends messages back through the websocket connections
	go server.handlePartyMessages()
	// start the http server
	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
