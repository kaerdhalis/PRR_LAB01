package main

import (
	"../network"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"
)

var portP2P int
var portMulticast int
var multicastAddr string
var artificialNetworkDelay time.Duration
func main(){
	argsWithoutProg := os.Args[1:]
	setup(argsWithoutProg)

	channel := make(chan network.MessageWithOrigin)

	go network.ClientReaderOnPort(portP2P, channel)

	go  synchronization()

	var buf bytes.Buffer
	for {
		select {
		case msgWithOrigin := <-channel:
			{
				msg:=msgWithOrigin.Msg
				buf.Reset()

				fmt.Println("Received DELAY_REQUEST from",msgWithOrigin.Ip)
				err := gob.NewEncoder(&buf).Encode(network.Message{Id: msg.Id, Time: time.Now(), Msg: 0b00})
				if err != nil {
					// handle error
				}
				time.Sleep(artificialNetworkDelay) //simulates network delay of artificialNetworkDelay milliseconds
				network.ClientWriter(msgWithOrigin.Ip,portP2P,buf)
			}
		}
	}
}

func setup(args [] string) {

	if !(len(args)== 2 || len(args)==3){
		wrongArguments()
	}

	fmt.Println("Starting Master with:")

	if len(args) >= 2 {

		p, err := strconv.Atoi(args[0])
		if err!= nil {
			wrongArguments()
		}
		portP2P = p
		fmt.Println("P2P Port=", portP2P)

		p, err = strconv.Atoi(args[1])
		if err!= nil {
			wrongArguments()
		}
		portMulticast=p

		fmt.Println("Multicast port=", portMulticast)
		multicastAddr=network.MulticastAddr+":"+strconv.Itoa(portMulticast)
		fmt.Println("Multicast addr=", multicastAddr)
		if len(args) == 3{
			drift, err := strconv.Atoi(args[2])
			if err!= nil{
				wrongArguments()
			}

			artificialNetworkDelay= time.Duration(drift) *time.Millisecond

		}else{
			artificialNetworkDelay=0
		}
		fmt.Println("Artificial network delay=", artificialNetworkDelay)

	}

	fmt.Println("=========================")
}

func wrongArguments(){
	fmt.Printf("WRONG FORMAT Correct usage: go run master.go <P2P port> <multicast port>  OR " +
		"go run master.go <P2P port> <multicast port> <artificial network delay> <artificial clock offset> ")
	os.Exit(1)
}
func synchronization()  {

		var buf bytes.Buffer

		for id:= 1;;{

		id +=1
		buf.Reset()
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: time.Time{}, Msg: 0b00}); err != nil {
			// handle error
		}
		time.Sleep(artificialNetworkDelay) //simulates network delay of artificialNetworkDelay milliseconds
		network.ClientWriter(multicastAddr,portP2P,buf)
		masterTime := time.Now()
		buf.Reset()
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: masterTime, Msg: 0b01}); err != nil {
			// handle error
		}
		time.Sleep(artificialNetworkDelay) //simulates network delay of artificialNetworkDelay milliseconds
		network.ClientWriter(multicastAddr,portP2P,buf)
		time.Sleep(2 * time.Second)
	}
}


