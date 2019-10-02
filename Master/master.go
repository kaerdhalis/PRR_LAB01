package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
)

const multicastAddr = "224.0.0.1:6666"

func main() {
	conn, err := net.ListenPacket("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, cliAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			s := s.Text() + " from " + cliAddr.String() + "\n"
			if _, err := conn.WriteTo([]byte(s), cliAddr); err != nil {
				log.Fatal(err)
			}
		}
	}
}
