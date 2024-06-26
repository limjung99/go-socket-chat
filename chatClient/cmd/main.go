package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var printChan chan string

func HandleSend(c net.Conn) {
	writer := bufio.NewWriter(c)
	for {
		fmt.Print("Enter message:")
		msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		writer.WriteString(msg)
		writer.Flush()
	}
}

func HandleReceive(c net.Conn) {
	reader := bufio.NewReader(c)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		printChan <- fmt.Sprintf("Received:%s", msg)
	}
}

func main() {
	dial, err := net.Dial("tcp", "127.0.0.1:8080")
	printChan = make(chan string)
	if err != nil {
		fmt.Println("Establish fail ...")
		panic(err)
	}
	fmt.Println("connection establish ...")
	go HandleSend(dial)
	go HandleReceive(dial)
	PrintConsole()
}

func PrintConsole() {
	for msg := range printChan {
		fmt.Println("\r\033[k")
		fmt.Println(msg)
		fmt.Println("Enter message:")
	}
}
