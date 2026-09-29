package main

import (
	"net"
)

type game struct {
	players [2]player
	turn    int
}

type player struct {
	name   string
	conn   net.Conn
	hand   []string
	points int
}
