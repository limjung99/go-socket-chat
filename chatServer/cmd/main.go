package main

import (
	"bufio"
	"fmt"
	"net"
)

var (
	conList []net.Conn
	msgChan chan string
)

func main() {
	fmt.Println("             ____  _____ ______     _______ ____              \n _____ _____/ ___|| ____|  _ \\ \\   / / ____|  _ \\ _____ _____ \n|_____|_____\\___ \\|  _| | |_) \\ \\ / /|  _| | |_) |_____|_____|\n|_____|_____|___) | |___|  _ < \\ V / | |___|  _ <|_____|_____|\n            |____/|_____|_| \\_\\ \\_/  |_____|_| \\_\\            ")
	conChan := make(chan net.Conn)
	conList = make([]net.Conn, 0, 5)
	msgChan = make(chan string)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Server listening on port 8080...")

	go ListenAndHandle(conChan)
	go SendMsgToCon()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Establish fail ...")
			continue
		}
		conChan <- conn
		conList = append(conList, conn)
	}
}

func SendMsgToCon() {
	// send msg to every con
	for msg := range msgChan {
		for _, con := range conList {
			writer := bufio.NewWriter(con)
			writer.WriteString(msg)
			writer.Flush()
		}
	}
}

func ListenAndHandle(conChan chan net.Conn) {
	for con := range conChan {
		go func(c net.Conn) {
			defer c.Close()
			fmt.Println("new connection establish ...")
			reader := bufio.NewReader(c)
			for {
				msg, err := reader.ReadString('\n')
				if err != nil {
					if err.Error() == "EOF" {
						fmt.Println("Client closed connection ... ")
						return
					}
					fmt.Println("error:", err)
					return
				}
				msgChan <- msg
			}
		}(con)
	}
}
