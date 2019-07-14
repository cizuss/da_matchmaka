package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func matchWorker(jobs chan *Player, parties *[]*Party, partiesChan chan *Party) {
	for player := range jobs {
		if player.foundParty || player.inProcess {
			continue
		}
		// try to find a party for the player
		party, err := player.findParty(*parties)
		if err != nil {
			fmt.Println("Player already searching for party by another thread...")
			continue
		}
		if party == nil {
			p := NewParty()
			party = &p
			addParty(parties, party)
		}
		if player.party != party {
			party.addPlayer(player)
		}
		if len(party.players) == 8 {
			party.markAllPlayersAsHaveFoundParty()
			fmt.Println("Formed the following party : ", party)
			partiesChan <- party
			removeParty(parties, party)
		} else {
			addTimerToReenqueuePlayer(player, jobs)
		}
	}
}

func addTimerToReenqueuePlayer(p *Player, jobs chan *Player) {
	// spawn a timer that ticks after 3 seconds
	timer := time.NewTimer(3 * time.Second)
	go func() {
		<-timer.C
		// if player didn't find party, increase delta and push him back in the queue
		p.lock()
		if !p.foundParty {
			p.party.removePlayer(p)
			p.delta = p.delta * 2
			p.party = nil
			p.inParty = false
			p.inProcess = false
			p.unlock()
			jobs <- p
		} else {
			p.unlock()
		}
	}()
}

func addParty(parties *[]*Party, party *Party) {
	*parties = append(*parties, party)
}

func removeParty(parties *[]*Party, party *Party) {
	result := make([]*Party, 0)
	for _, p := range *parties {
		if p != party {
			result = append(result, p)
		}
	}
	*parties = result
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numWorkers := 4
	parties := make([]*Party, 0)
	var jobsChannel = make(chan *Player, 100)
	partiesChan := make(chan *Party, 100)
	connectionsMap := map[string]*websocket.Conn{}
	connections := map[*websocket.Conn]bool{}
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	for w := 0; w < numWorkers; w++ {
		go matchWorker(jobsChannel, &parties, partiesChan)
	}

	httpServer := MHttpServer{jobsChannel, partiesChan, connectionsMap, connections, upgrader}
	httpServer.start()
}
