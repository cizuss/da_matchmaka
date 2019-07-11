package main

import (
	"sync"
)

type Party struct {
	players   []*Player
	avgSkill  int
	createdAt int64
	mux       sync.Mutex
}

func (party *Party) addPlayer(player *Player) {
	party.mux.Lock()
	defer party.mux.Unlock()
	for _, p := range party.players {
		if p.name == player.name {
			return
		}
	}
	party.players = append(party.players, player)
	party.computeAvgSkill()
	player.party = party
}

func (party *Party) removePlayer(player *Player) {
	if party == nil {
		return
	}
	party.mux.Lock()
	defer party.mux.Unlock()
	idx := -1
	for i, p := range party.players {
		if p.name == player.name {
			idx = i
			break
		}
	}
	if idx != -1 {
		party.players[idx] = party.players[len(party.players)-1]
		party.players = party.players[:len(party.players)-1]
	}
	party.computeAvgSkill()
	player.party = nil
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
