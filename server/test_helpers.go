package main

import (
	"math/rand"
)

func randomPlayer() Player {
	skill := rand.Intn(1000)
	name := randSeq(10)
	p := NewPlayer(name, skill)

	return p
}
