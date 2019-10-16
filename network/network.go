package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Message struct{
	Id int
	Time time.Time
	Msg byte
	OriginPort int
}

type MessageWithOrigin struct{
	Msg Message
	Ip string
}


const MulticastAddr = "224.0.1.1:6666"
const SrvAddrSlave = "127.0.0.1:6060"
const SrvAddrMaster  = "127.0.0.1:5010"

func ClientWriter(address string,buf bytes.Buffer) {

	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err  = buf.WriteTo(conn)
}

func ClientReaderMult(address string, channel chan MessageWithOrigin) {
	fmt.Println()
	// error testing suppressed to compact listing on slides
	conn, _ := net.ListenPacket("udp", address) // listen on port
	defer conn.Close()
	p := ipv4.NewPacketConn(conn) // convert to ipv4 packetConn
	addr, _ := net.ResolveUDPAddr("udp", address)
	p.JoinGroup(nil, addr) // darwin : interface en0
	decrypt(conn,channel)

}

func ClientReaderPort(port int,channel chan MessageWithOrigin)  {
	var localAddr string
	localAddr ="127.0.0.1:"+ strconv.Itoa(port)
	ClientReader(localAddr,channel)

}
func ClientReader(address string, channel chan MessageWithOrigin) {
	// error testing suppressed to compact listing on slides

	conn, err := net.ListenPacket("udp", address)



	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	decrypt(conn,channel)

}

func decrypt(conn net.PacketConn ,channel chan MessageWithOrigin){


	buf := make([]byte, 1024)
	for {
		var result MessageWithOrigin
		var ip net.Addr
		n, ip, err := conn.ReadFrom(buf) // n,addr, err := p.ReadFrom(buf)

		cleanedIp:=strings.Split(ip.String(),":")[0]
		result.Ip=cleanedIp

		if err != nil {
			log.Fatal(err)
		}

		var msg Message
		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&msg); err != nil {
			// handle error
		}



		result.Msg= msg


		channel <- result
	}
}

