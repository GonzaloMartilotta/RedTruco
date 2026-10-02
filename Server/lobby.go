package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func runLobby() {
	ln, err := net.Listen("tcp", ":8080") // Acepta todo lo que entre por el puerto 8080
	if err != nil {
		fmt.Println("Error to listen:", err)
		return
	}
	var rooms []*game
	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Println("ERROR al conectar:", err)
			continue
		}
		go handleHTTP(conn, &rooms)
	}
}

func handleHTTP(conn net.Conn, rooms *[]*game) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	request := scanner.Text()
	parts := strings.Fields(request)
	if len(parts) != 3 {
		sendHTTP(conn, "400 Bad Request", "")
		return
	}

	switch parts[0] {
	case "GET":
		if parts[1] == "/salas" {
			if len(*rooms) > 0 {
				sendHTTP(conn, "200 OK", "Hay salas")
			} else {
				sendHTTP(conn, "200 OK", "No hay salas")
			}
		} else {
			sendHTTP(conn, "404 Not Found", "")
		}
	case "POST":
		room := new(game)
		*rooms = append(*rooms, room)
		sendHTTP(conn, "201 Created", fmt.Sprintf("Room %d created", len(*rooms)-1))
	default:
		sendHTTP(conn, "405 Bad Request", "")
	}
}

func sendHTTP(client net.Conn, status string, body string) {
	fmt.Fprintf(client, "HTTP/1.1 %s\r\n"+
		"Content-Type: text/plain\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n"+
		"\r\n"+
		"%s\n", status, len(body)+1, body)
}
