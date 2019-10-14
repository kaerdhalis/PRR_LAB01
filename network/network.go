package network

import (
	"bytes"
	"encoding/gob"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"time"
)

type Message struct{
	Id int
	Time time.Time
	Msg byte
}

const MulticastAddr = "224.0.1.1:6666"
const SrvAddrS = "127.0.0.1:6060"
const SrvAddrM  = "127.0.0.2:5050"

func ClientWriter(address string,buf bytes.Buffer) {

	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err  = buf.WriteTo(conn)
}

func ClientReaderMult(address string, channel chan Message) {

	// error testing suppressed to compact listing on slides
	conn, _ := net.ListenPacket("udp", address) // listen on port
	defer conn.Close()
	p := ipv4.NewPacketConn(conn) // convert to ipv4 packetConn
	addr, _ := net.ResolveUDPAddr("udp", address)
	p.JoinGroup(nil, addr) // darwin : interface en0
	decrypt(conn,channel)

}

func ClientReader(address string, channel chan Message) {
	// error testing suppressed to compact listing on slides
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	decrypt(conn,channel)

}

func decrypt(conn net.PacketConn ,channel chan Message){

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFrom(buf) // n, _, addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		var msg Message
		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&msg); err != nil {
			// handle error
		}
		channel <- msg
	}
}

