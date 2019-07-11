package main

import (
	"math"
	"sync"
	"time"
)

type Player struct {
	name       string
	skill      int
	timestamp  int64
	foundParty bool
	delta      int
	party      *Party
	mux        sync.Mutex
}

func NewPlayer(name string, skill int) Player {
	return Player{name: name, skill: skill, timestamp: time.Now().Unix(), foundParty: false, delta: 2, party: nil}
}

func (player Player) findParty(parties []*Party) *Party {
	var goodParties []*Party
	for _, party := range parties {
		if isPartyGoodForPlayer(player, party) {
			goodParties = append(goodParties, party)
		}
	}
	return findBestParty(goodParties)
}

func findBestParty(parties []*Party) *Party {
	if len(parties) == 0 {
		return nil
	}

	var bestParty *Party = parties[0]
	var maxLen int = len(bestParty.players)
	for _, p := range parties {
		if len(p.players) > maxLen {
			maxLen = len(p.players)
			bestParty = p
		}
	}
	var minCrAt int64 = math.MaxInt64
	for _, p := range parties {
		if len(p.players) == maxLen {
			if p.createdAt < minCrAt {
				minCrAt = p.createdAt
				bestParty = p
			}
		}
	}
	return bestParty
}

func isPartyGoodForPlayer(player Player, party *Party) bool {
	ps := player.skill
	as := party.avgSkill
	d := player.delta

	return ps >= (as-d) && ps <= (as+d)
}
