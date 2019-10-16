package network

import (
	"bytes"
	"encoding/gob"
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
}

type MessageWithOrigin struct{
	Msg Message
	Ip string
}


const MulticastAddr = "224.0.1.1"

func ClientWriter(address string,fromPort int,buf bytes.Buffer) {
	var localAddr= new (net.UDPAddr)
	localAddr.Port= fromPort

	var ipPort=strings.Split(address,":")
	var remoteAddr = new (net.UDPAddr)
	remoteAddr.IP= net.ParseIP(ipPort[0])
	remoteAddr.Port,_= strconv.Atoi(ipPort[1])

	conn, err := net.DialUDP("udp",localAddr, remoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err  = buf.WriteTo(conn)
}
func ClientReaderMult(address string, channel chan MessageWithOrigin) {

	// error testing suppressed to compact listing on slides
	conn, _ := net.ListenPacket("udp", address) // listen on port
	defer conn.Close()
	p := ipv4.NewPacketConn(conn) // convert to ipv4 packetConn
	addr, _ := net.ResolveUDPAddr("udp", address)
	p.JoinGroup(nil, addr) // darwin : interface en0
	decrypt(conn,channel)

}

func ClientReaderOnPort(port int,channel chan MessageWithOrigin)  {
	var localAddr string
	localAddr ="127.0.0.1:"+ strconv.Itoa(port)
	clientReader(localAddr,channel)
}
func clientReader(address string, channel chan MessageWithOrigin) {
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

		result.Ip=ip.String()

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

