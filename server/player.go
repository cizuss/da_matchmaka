package main

import (
	"errors"
	"math"
	"sync"
	"time"
)

type Player struct {
	mux        sync.Mutex
	name       string
	skill      int
	timestamp  int64
	foundParty bool
	delta      int
	party      *Party
	inParty    bool
	inProcess  bool
}

func NewPlayer(name string, skill int) Player {
	return Player{name: name, skill: skill, timestamp: time.Now().Unix(), foundParty: false, delta: 2, party: nil, inParty: false, inProcess: false}
}

func (player *Player) findParty(parties []*Party) (*Party, error) {
	if player.inProcess {
		return nil, errors.New("player already searching for party...")
	}
	player.inProcess = true
	var goodParties []*Party
	for _, party := range parties {
		if isPartyGoodForPlayer(*player, party) {
			goodParties = append(goodParties, party)
		}
	}
	return findBestParty(goodParties), nil
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
