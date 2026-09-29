package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	mensaje := []byte("Hola prueba por TCP\n")
	_, err = conn.Write(mensaje)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Mensaje enviado")
}
