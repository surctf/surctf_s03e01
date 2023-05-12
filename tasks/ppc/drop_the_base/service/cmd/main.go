package main

import (
	"bufio"
	"easyppc/internal"
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	HOST   = "0.0.0.0"
	PORT   = "9876"
	TYPE   = "tcp"
	LEVELS = 1000
	FLAG   = "surctf_drop_the_base_or_drop_the_bass"
)

func main() {
	log.Printf("Starting TCP server on %s:%s", HOST, PORT)
	listener, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Println("Error connecting:", err.Error())
			return
		}

		log.Printf("Client %s connected\n", c.RemoteAddr().String())

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for i := 0; i < LEVELS; i++ {
		conn.Write([]byte("[" + strconv.Itoa(i+1) + "] "))

		task := internal.NewTask()
		conn.Write([]byte(fmt.Sprintf("Base(%s): %s\n", task.Exp, task.Encoded)))

		conn.Write([]byte("Decoded: "))

		resp, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			log.Println("Client left")
			return
		}

		if string(resp[:len(resp)-1]) != task.Raw {
			conn.Write([]byte("Wrong!\n"))
			conn.Write([]byte(fmt.Sprintf("Decoded: %s\n", task.Raw)))
			return
		}
	}

	conn.Write([]byte("NICE CO...khm.khm...BASE\n"))
	conn.Write([]byte(FLAG + "\n"))
}
