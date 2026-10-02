package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	var rooms []*game
	go runLobby(&rooms)
	ln, err := net.Listen("tcp", ":9000") // Acepta todo lo que entre por el puerto 9000
	if err != nil {
		fmt.Println("Error to listen:", err)
		return
	}
	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Println("ERROR al conectar:", err)
			continue
		}

		go handleJoin(conn, &rooms)
	}
}

func handleJoin(conn net.Conn, rooms *[]*game) {
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	message := scanner.Text()
	parts := strings.Fields(message)
	if len(parts) != 2 {
		send(conn, "ERROR Bad request")
		conn.Close()
		return
	} else if parts[0] != "SALA" {
		send(conn, "ERROR Select a room")
		conn.Close()
		return
	}
	roomID, err := strconv.Atoi(parts[1])

	if err != nil {
		send(conn, fmt.Sprintf("ERROR %v", err))
		conn.Close()
		return
	}

	roomsMu.RLock() // Bloquea para evitar race condition
	if len(*rooms) < roomID+1 || 0 > roomID {
		roomsMu.RUnlock()
		send(conn, "ERROR Room don't found")
		conn.Close()
		return
	}
	room := (*rooms)[roomID]
	roomsMu.RUnlock()

	room.mu.Lock()
	// Se fija si hay espacio en la sala, si hay asigna el jugador
	if room.players[0].conn == nil {
		room.players[0].conn = conn
	} else if room.players[1].conn == nil {
		room.players[1].conn = conn

		// Estan los dos jugadores => Empieza la partida
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
	} else {
		room.mu.Unlock()
		send(conn, "ERROR Room is full")
		conn.Close()
		return
	}
	room.mu.Unlock()
	send(conn, "OK")
	//handleConnection(room, player)
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

func send(player net.Conn, message string) {
	fmt.Fprintln(player, message)
}
