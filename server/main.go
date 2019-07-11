// In this example we'll look at how to implement
// a _worker pool_ using goroutines and channels.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

var jobsChannel = make(chan *Player, 100)

func matchWorker(jobs chan *Player, parties *[]*Party) {
	for player := range jobs {
		// try to find a party for the player
		//fmt.Printf("Trying to find party for %v\n", player)
		party := player.findParty(*parties)
		if party == nil {
			p := NewParty()
			party = &p
			addParty(parties, party)
		}
		if player.party != party {
			party.addPlayer(player)
		}
		if len(party.players) == 8 {
			handleFoundParty(*party)
			for _, p := range party.players {
				p.mux.Lock()
				p.party = nil
				p.foundParty = true
				p.mux.Unlock()
			}
			removeParty(parties, party)
		} else {
			timer := time.NewTimer(3 * time.Second)
			go func(p *Player) {
				<-timer.C
				if p.foundParty {
					return
				}
				p.delta = player.delta + 2
				p.party.removePlayer(player)
				jobs <- p
			}(player)
		}
	}
}

func addParty(parties *[]*Party, party *Party) {
	*parties = append(*parties, party)
}

func handleFoundParty(party Party) {
	fmt.Println("Found party for the following players:")
	for _, p := range party.players {
		fmt.Printf("%v ", p)
	}
	now := time.Now().Unix()
	var sum int64 = 0
	for _, p := range party.players {
		sum += now - p.timestamp
		fmt.Printf("Player %s found a party in %d seconds\n", p.name, now-p.timestamp)
	}
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

func addPlayerToQueue(p *Player) {
	//fmt.Println("Adding player ", p, " to queue")
	jobsChannel <- p
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numWorkers := 4
	parties := make([]*Party, 0)

	for w := 0; w < numWorkers; w++ {
		go matchWorker(jobsChannel, &parties)
	}

	startHttpServer()
}

func NewParty() Party {
	return Party{players: []*Player{}, avgSkill: 0, createdAt: time.Now().Unix()}
}
