package main


import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/ipv4"
	"io"
	"log"
	"net"
	"os"
)

const multicastAddr = "224.0.0.1:6666"

func main() {
	go clientReader()
	conn, err := net.Dial("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(conn, os.Stdin)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func clientReader() {
	// error testing suppressed to compact listing on slides
	conn, _ := net.ListenPacket("udp", multicastAddr) // listen on port
	defer conn.Close()
	p := ipv4.NewPacketConn(conn) // convert to ipv4 packetConn
	addr, _ := net.ResolveUDPAddr("udp", multicastAddr)
	p.JoinGroup(nil, addr) // darwin : interface en0
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf) // n, _, addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			fmt.Printf("%s from %v\n", s.Text(), addr)
		}
	}
}