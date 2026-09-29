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
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error to connect:", err)
			continue
		}
		go handleConnection(conn)
	}
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
