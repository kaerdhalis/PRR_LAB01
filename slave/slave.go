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
var artificialClockOffset time.Duration

var port int
var masterP2PAddr string
var masterMultiCastAddr string


func main() {
	argsWithoutProg := os.Args[1:]
	setup(argsWithoutProg)

	channelMult := make(chan network.MessageWithOrigin)

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
			msg := msgWithOrigin.Msg
			if msg.Msg == 0b00 {

				localTimeWhenSyncReceived = time.Now().Add( artificialClockOffset)
				id = msg.Id
			} else if msg.Msg == 0b01 {

				if id == msg.Id {
					masterTime = msg.Time

					timeDelta = masterTime.Sub(localTimeWhenSyncReceived)

					fmt.Println("actual time difference ", timeDelta.String())
				}
			}

		case delayMsgWithOrigin := <-channelP2P:
			delayMsg := delayMsgWithOrigin.Msg
			if idDelay == delayMsg.Id {

				currentNetworkDelay = delayMsg.Time.Sub(localTimeWhenLastDelayRequestSent)

				fmt.Println("transmission delay n°", idDelay, " =", currentNetworkDelay)

			}
		}

	}
}

func setup(args [] string) {

	if !(len(args)== 3 || len(args)==4){
		wrongArguments()
	}

	fmt.Println("Starting new Slave with:")

	if len(args) >= 3 {

		p, err := strconv.Atoi(args[0])
		if err!= nil {
			wrongArguments()
		}
		port = p
		fmt.Println("Port=", port)
		masterP2PAddr = args[1]
		masterMultiCastAddr = args[2]
		fmt.Println("Master P2P ip=", masterP2PAddr)
		fmt.Println("Master MultiCast ip =", masterMultiCastAddr)

		if len(args) == 4 {
			drift, err := strconv.Atoi(args[3])
			if err!= nil{
				wrongArguments()
			}

			artificialClockOffset= time.Duration(drift) *time.Millisecond

		}else{
			artificialClockOffset=0
		}
		fmt.Println("artificial clock Drift=", artificialClockOffset)
	}

	fmt.Printf("=========================")
}
/**
Prints the wrong format error and exits application
 */
func wrongArguments(){
	fmt.Printf("WRONG FORMAT Correct usage: go run slave.go <port> <masterP2PIp> <masterMultCastIp> OR " +
		"go run slave.go <port> <masterP2PIp> <masterMultCastIp> <artificial clock offset> ")
	os.Exit(1)
}

/**
Sends a Delay request to the server at a random interval between 5 and 15 seconds
 */
func delayRequest() {

	var buf bytes.Buffer
	var timeTilNextDelayRequest = rand.Intn(10) + 5
	for {

		localTimeWhenLastDelayRequestSent = time.Now()

		time.Sleep(time.Duration(timeTilNextDelayRequest) * time.Second)
		idDelay++
		buf.Reset()
		fmt.Println("Sending delay request n°", idDelay)
		err := gob.NewEncoder(&buf).Encode(network.Message{Id: idDelay, Time: time.Time{}, Msg: 0b10, OriginPort: port})
		if  err != nil {
			// handle error
		}

		network.ClientWriter(masterP2PAddr, buf)

		timeTilNextDelayRequest = rand.Intn(15)

	}
}
