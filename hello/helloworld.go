package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Welcome to the playground!")

	conn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		return
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	scanner := bufio.NewScanner(bufio.NewReader(conn))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
