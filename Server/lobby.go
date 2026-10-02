package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type request struct {
	method  string
	path    string
	version string
	headers map[string]string
}

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
	requestLine := scanner.Text()
	parts := strings.Fields(requestLine)

	if len(parts) != 3 {
		sendHTTP(conn, "400 Bad Request", "")
		return
	}

	req := request{method: parts[0], path: parts[1], version: parts[2], headers: make(map[string]string)}

	for scanner.Scan() {
		requestLine := scanner.Text()
		if requestLine == "" {
			break
		}
		headerParts := strings.SplitN(requestLine, ":", 2) // Separa en 2 partes los headers
		key := strings.TrimSpace(headerParts[0])
		value := strings.TrimSpace(headerParts[1])
		req.headers[key] = value
	}
	switch req.method {
	case "GET":
		if req.path == "/salas" {
			if len(*rooms) > 0 {
				sendHTTP(conn, "200 OK", "Hay salas")
			} else {
				sendHTTP(conn, "200 OK", "No hay salas")
			}
		} else {
			sendHTTP(conn, "404 Not Found", "")
		}
	case "POST":
		if req.path == "/salas" {
			room := new(game)
			*rooms = append(*rooms, room)
			sendHTTP(conn, "201 Created", fmt.Sprintf("Room %d created", len(*rooms)-1))
		} else {
			sendHTTP(conn, "404 Not Found", "")
		}
	default:
		sendHTTP(conn, "405 Method Not Allowed", "")
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
