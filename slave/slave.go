package main

import (
	"../network"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var idDelay int
var localTimeWhenLastDelayRequestSent time.Time
var currentNetworkDelay time.Duration
var clockDrift int = 0
var networkDelay int = 0
var port int
var masterP2PAddr string
var masterMultiCastAddr string

func main() {
	argsWithoutProg := os.Args[1:]
	setup(argsWithoutProg)

	channelMult := make(chan network.MessageWithOrigin)
	fmt.Println(port)
	channelP2P := make(chan network.MessageWithOrigin)
	go network.ClientReaderMult(masterMultiCastAddr, channelMult)
	go network.ClientReaderPort(port, channelP2P)

	go delayRequest()

	var masterTime, localTimeWhenSyncReceived time.Time
	var timeDelta time.Duration
	var id int

	var buf bytes.Buffer
	for {

		select {
		case msgWithOrigin := <-channelMult:
			msg:= msgWithOrigin.Msg
			if msg.Msg == 0b00 {

				localTimeWhenSyncReceived = time.Now()
				id = msg.Id
			} else if msg.Msg == 0b01 {

				if id == msg.Id {
					masterTime = msg.Time

					timeDelta = masterTime.Sub(localTimeWhenSyncReceived)

					fmt.Println("actual time difference ", timeDelta.String())

					buf.Reset()
					if err := gob.NewEncoder(&buf).Encode(network.Message{Id: idDelay, Time: time.Now(), Msg: 0b00}); err != nil {
						// handle error
					}
					//network.ClientWriter(network.SrvAddrM,buf)

				}
			}

		case delayMsgWithOrigin := <-channelP2P:
			delayMsg:= delayMsgWithOrigin.Msg
			if idDelay == delayMsg.Id {

				currentNetworkDelay = delayMsg.Time.Sub(localTimeWhenLastDelayRequestSent)

				fmt.Println("transmission delay n°", idDelay, " =", currentNetworkDelay)

			}
		}

	}
}

func setup(args [] string) {

	switch len(args) {

	case 3:
		p, err := strconv.Atoi(args[0])
		if err == nil {
			port= p
			fmt.Println("Port=", port)
		}

		masterP2PAddr = args[1]
		masterMultiCastAddr = args[2]
		fmt.Println("Master P2P ip=", masterP2PAddr)
		fmt.Println("Master MultiCast ip =", masterMultiCastAddr)
		break


	case 5:
		p, err := strconv.Atoi(args[0])
		if err == nil {
			port=p
			fmt.Println("Port=", port)
		}

		masterP2PAddr = args[1]
		masterMultiCastAddr = args[2]
		fmt.Println("Master P2P ip=", masterP2PAddr)
		fmt.Println("Master MultiCast ip =", masterMultiCastAddr)
		clockDrift, err := strconv.Atoi(args[3])
		if err == nil {
			fmt.Println("clockDrift=", clockDrift)
		}
		networkDelay, err := strconv.Atoi(args[4])
		if err == nil {
			fmt.Println("networkDelay=", networkDelay)
		}
		break
	default:
		fmt.Printf("WRONG FORMAT Correct usage: go run slave.go <port> <masterP2PIp> <masterMultCastIp> OR go run slave.go <port> <masterP2PIp> <masterMultCastIp> <clockdrift> <networkDelay>")
		os.Exit(1)
	}


}

func delayRequest() {

	var buf bytes.Buffer
	var timeTilNextDelayRequest = rand.Intn(10) + 5
	for {

		localTimeWhenLastDelayRequestSent = time.Now()

		time.Sleep(time.Duration(timeTilNextDelayRequest) * time.Second)
		idDelay++
		buf.Reset()
		fmt.Println("Sending delay request n°", idDelay)

		if err := gob.NewEncoder(&buf).Encode(network.Message{Id: idDelay, Time: time.Time{}, Msg: 0b10,OriginPort:port}); err != nil {
			// handle error
		}

		network.ClientWriter(masterP2PAddr, buf)

		timeTilNextDelayRequest = rand.Intn(15)

	}
}
