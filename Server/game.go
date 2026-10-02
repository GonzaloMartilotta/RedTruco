package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

var roomsMu sync.RWMutex

type game struct {
	players [2]player
	turn    int
	mu      sync.Mutex
}

type player struct {
	name   string
	conn   net.Conn
	hand   []string
	played []bool
	points int
}

func handlePlays(room *game, playerInt int, message string) {
	player := &room.players[playerInt]
	conn := player.conn

	if room.turn != playerInt {
		send(conn, "ERROR no es tu turno")
		return
	}
	cardInt, err := strconv.Atoi(message)
	if err != nil || len(message) != 1 || cardInt > 3 || cardInt <= 0 { // Filtra lo que no sea 1, 2 o 3
		send(conn, "Seleccione una carta valida")
		return
	} else if player.played[cardInt-1] == true {
		send(conn, "Carta ya jugada")
		return
	}
	card := room.players[playerInt].hand[cardInt-1]
	player.played[cardInt-1] = true
	fmt.Printf("Carta %s jugada por jugador %d\n", card, playerInt)
	send(conn, "OK")
}
