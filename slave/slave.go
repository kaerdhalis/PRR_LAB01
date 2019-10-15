package main

import (
	"../network"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"time"
)
var idDelay int
var localTimeWhenLastDelayRequestSent time.Time
var currentNetworkDelay time.Duration;
func main() {

	channelMult := make(chan network.Message)
	channelP2P := make(chan network.Message)
	go network.ClientReaderMult(network.MulticastAddr,channelMult)
	go network.ClientReader(network.SrvAddrS,channelP2P)

	go delayRequest()

	var tMaitre,ti time.Time
	var ecart time.Duration
	var id  int

	var buf bytes.Buffer
	for;;{

		idDelay = 0

		select {
		case msg := <-channelMult:
			if msg.Msg == 0b00{

				ti = time.Now()
				id = msg.Id
			}else if msg.Msg == 0b01{

				if id == msg.Id {
					tMaitre = msg.Time

					ecart = tMaitre.Sub(ti)

					fmt.Println("actual ecart ",ecart.String())

					buf.Reset()
					if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: idDelay, Time: time.Now(), Msg: 0b00}); err != nil {
						// handle error
					}
					//network.ClientWriter(network.SrvAddrM,buf)
					idDelay++
				}
			}

		case msgd := <-channelP2P:

			if idDelay == msgd.Id {


				var delay = msgd.Time.Sub(localTimeWhenLastDelayRequestSent)

				fmt.Println("transmission delay n°",idDelay," =", delay)


			}
		}
	


	}
}

func delayRequest()  {

	var buf bytes.Buffer
	var timeTilNextDelayRequest= rand.Intn(10)+5
	for {

		localTimeWhenLastDelayRequestSent= time.Now()

		time.Sleep(time.Duration(timeTilNextDelayRequest) * time.Second)
		idDelay++
		buf.Reset()



		fmt.Println ("Sending delay request n°", idDelay)


		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id:idDelay, Time: time.Time{}, Msg: 0b10}); err != nil {
			// handle error
		}


		network.ClientWriter(network.SrvAddrM,buf)

		timeTilNextDelayRequest= rand.Intn(15);

	}
}
