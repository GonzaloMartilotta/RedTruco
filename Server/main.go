package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":9000") // Acepta todo lo que entre por el puerto 9000
	if err != nil {
		fmt.Println("Error to listen:", err)
		return
	}
	var players [2]net.Conn

	for i := range players {
		players[i], err = ln.Accept()
		if err != nil {
			fmt.Println("Error to connect:", err)
			return
		}

	}
	var hand []string
	deck := newDeck()
	shuffleDeck(deck)
	//fmt.Println(deck) // Ver el mazo mezclado
	for i := range players {
		send(players[i], "INICIO 15")
		hand, deck = dealCards(deck)
		send(players[i], strings.Join(hand, " "))
		go handleConnection(players[i])
	}
	select {}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	client := conn.RemoteAddr()
	fmt.Printf("Client: %s connected\n", client)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		handleMessage(conn, scanner.Text())
	}

	err := scanner.Err()
	if err != nil {
		fmt.Printf("Client: %s disconnected with error: %v\n", client, err)
	} else {
		fmt.Printf("Client: %s disconnected\n", client)
	}

}

func send(player net.Conn, message string) {
	fmt.Fprintln(player, message)
}

func handleMessage(conn net.Conn, message string) {
	command := strings.Fields(message)

	if len(command) <= 0 {
		fmt.Fprint(conn, "ERROR mensaje vacio\n")
		return
	}

	switch command[0] {
	case "JUGAR", "TRUCO", "RETRUCO", "VALE4", "ENVIDO", "QUIERO", "NOQUIERO":
		fmt.Fprint(conn, "OK\n")
	default:
		fmt.Fprint(conn, "ERROR comando desconocido\n")
	}

}
