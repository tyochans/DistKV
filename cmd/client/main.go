package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err!=nil {
		fmt.Println("Server Not Connected")
	}
	defer conn.Close()

	fmt.Println("Server Connected on :9000")

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Println(">")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if(input == "exit") {

			fmt.Println("Closing Application")
			return
		}

		// didnt understand what this does
		fmt.Fprintln(conn, input)

		response, _ := serverReader.ReadString('\n')

		fmt.Println(response)
	}


}