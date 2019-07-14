package main

import (
	"sync"
	"time"
)

// Party is a struct for a group of players searching for a match
type Party struct {
	mux       sync.Mutex
	id        string
	players   []*Player
	avgSkill  int
	createdAt int64
}

//NewParty returns a new Party instance
func NewParty() Party {
	return Party{id: randSeq(10), players: []*Player{}, avgSkill: 0, createdAt: time.Now().Unix()}
}

func (party *Party) lock() {
	party.mux.Lock()
}

func (party *Party) unlock() {
	party.mux.Unlock()
}

func (party *Party) addPlayer(player *Player) {
	if party == nil {
		return
	}
	// avoid case where concurrent threads find multiple parties for the same player
	player.lock()
	if player.inParty || player.foundParty {
		player.unlock()
		return
	}
	player.party = party
	player.inParty = true
	player.unlock()
	party.lock()
	defer party.unlock()
	party.players = append(party.players, player)
	party.computeAvgSkill()
}

func (party *Party) removePlayer(player *Player) {
	result := make([]*Player, 0)
	party.lock()
	defer party.unlock()
	for _, p := range party.players {
		if p.name != player.name {
			result = append(result, p)
		}
	}
	party.players = result
	party.computeAvgSkill()
}

func (party *Party) isEmpty() bool {
	return len(party.players) == 0
}

func (party *Party) computeAvgSkill() {
	if party.isEmpty() {
		party.avgSkill = 0
		return
	}
	sum := 0
	for _, p := range party.players {
		sum += p.skill
	}
	party.avgSkill = sum / len(party.players)
}

func (party *Party) markAllPlayersAsHaveFoundParty() {
	party.lock()
	for _, p := range party.players {
		p.foundParty = true
	}
	party.unlock()
}
