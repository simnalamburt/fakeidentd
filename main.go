package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"

	b64 "encoding/base64"
)

var re = regexp.MustCompile(`^ *(\d{1,5}) *, *(\d{1,5}) *$`)

func main() {
	// Start TCP server
	l, err := net.Listen("tcp", ":113")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Print("Started fakeidentd on port 113")

	for {
		// Wait for a new client
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		// Handle client in separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Parse lines
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		// Drop connection without response on invalid reqeust. This behavior is allowed by RFC 1413.
		if match == nil {
			log.Printf("Invalid request: \"%s\"", line)
			return
		}

		// Parse port numbers
		serverport, err := parsePort(match[1])
		if err != nil {
			io.WriteString(conn, fmt.Sprintf("%v, %v : ERROR : INVALID-PORT\r\n", match[1], match[2]))
			log.Print("Invalid port:", err)
			continue
		}
		clientport, err := parsePort(match[2])
		if err != nil {
			io.WriteString(conn, fmt.Sprintf("%v, %v : ERROR : INVALID-PORT\r\n", match[1], match[2]))
			log.Print("Invalid port:", err)
			continue
		}

		// Response with fake identity
		//
		// name = [serverport clientport] XOR [4B 00 DD DB]
		name := b64.RawURLEncoding.EncodeToString([]byte{
			byte(serverport>>8) ^ 0x4B,
			byte(serverport&0xff) ^ 0x00,
			byte(clientport>>8) ^ 0xDD,
			byte(clientport&0xff) ^ 0xDB,
		})
		output := fmt.Sprintf("%v, %v : USERID : UNIX : %v\r\n", serverport, clientport, name)
		if _, err := io.WriteString(conn, output); err != nil {
			log.Print(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
	}
}

func parsePort(input string) (uint16, error) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	if i < 1 || i > 65535 {
		return 0, errors.New("Invalid port number")
	}
	return uint16(i), nil
}
