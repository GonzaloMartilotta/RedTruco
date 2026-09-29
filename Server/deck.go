package main

import (
	"fmt"
	"math/rand"
)

func newDeck() []string {
	var deck []string
	var number = []int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12}
	var palo = []string{"B", "O", "E", "C"}
	for i := range palo {
		for j := range number {
			deck = append(deck, fmt.Sprintf("%d%s", number[j], palo[i]))
		}
	}
	return deck
}

func shuffleDeck(deck []string) {
	rand.Shuffle(len(deck), func(i int, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

func dealCards(deck []string) ([]string, []string) {
	hand := deck[:3]
	deck = deck[3:]
	return hand, deck
}
