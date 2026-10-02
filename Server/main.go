package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	go runLobby()
	ln, err := net.Listen("tcp", ":9000") // Acepta todo lo que entre por el puerto 9000
	if err != nil {
		fmt.Println("Error to listen:", err)
		return
	}
	room := &game{} // Crea una nueva sala
	for i := range room.players {
		room.players[i].conn, err = ln.Accept()
		if err != nil {
			fmt.Println("Error to connect:", err)
			return
		}
	}
	deck := newDeck()
	shuffleDeck(deck)
	//fmt.Println(deck) // Ver el mazo mezclado
	for i := range room.players {
		send(room.players[i].conn, "INICIO 15")
		room.players[i].hand, deck = dealCards(deck) // Les asigna cartas a los jugadores
		room.players[i].played = []bool{false, false, false}
		send(room.players[i].conn, "MANO "+strings.Join(room.players[i].hand, " "))
		go handleConnection(room, i)
	}
	select {}
}

func handleConnection(room *game, player int) {
	conn := room.players[player].conn
	defer conn.Close()
	client := conn.RemoteAddr()
	fmt.Printf("Client: %s connected\n", client)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		handleMessage(room, player, scanner.Text())
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

func handleMessage(room *game, player int, message string) {
	conn := room.players[player].conn

	command := strings.Fields(message)

	if len(command) <= 0 {
		send(conn, "ERROR mensaje vacio")
		return
	}

	switch command[0] {
	case "JUGAR":
		if len(command) < 2 {
			send(conn, "Seleccione una carta valida")
			break
		}
		handlePlays(room, player, command[1])
	default:
		send(conn, "ERROR comando desconocido")
	}
}
