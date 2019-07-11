package main

import (
	"fmt"
	"testing"
)

func TestRemoveParty(t *testing.T) {
	parties := []*Party{}
	for i := 0; i < 3; i++ {
		p := NewParty()
		addParty(&parties, &p)
	}
	assert(len(parties) == 3, t)
	removeParty(&parties, parties[0])
	assert(len(parties) == 2, t)
}

func TestRemovePlayerFromParty(t *testing.T) {
	party := NewParty()
	pl1 := randomPlayer()
	pl2 := randomPlayer()
	pl3 := randomPlayer()
	party.addPlayer(&pl1)
	party.addPlayer(&pl2)
	party.addPlayer(&pl3)
	assert(len(party.players) == 3, t)

	fmt.Println(pl2)
	fmt.Println(party.players)
	party.removePlayer(&pl2)
	fmt.Println(party.players)
	assert(len(party.players) == 2, t)
}

func TestRemovePlayerFromPartyInPartyArray(t *testing.T) {
	party := NewParty()
	pl1 := randomPlayer()
	pl2 := randomPlayer()
	party.addPlayer(&pl1)
	party.addPlayer(&pl2)

	parties := []*Party{}
	addParty(&parties, &party)

	assert(len(parties[0].players) == 2, t)
	parties[0].removePlayer(&pl1)
	assert(len(parties[0].players) == 1, t)
}

func assert(cond bool, t *testing.T) {
	if !cond {
		t.Errorf("Assertion failed")
	}
}
