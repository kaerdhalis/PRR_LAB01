package main

import (
	"Laboratoire01/network"
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

func main() {

	channelMult := make(chan network.Message)
	channelP2P := make(chan network.Message)
	go network.ClientReaderMult(network.MulticastAddr,channelMult)
	go network.ClientReader(network.SrvAddrS,channelP2P)


	var tMaitre,ti,tes time.Time
	var ecart time.Duration
	var id ,idDelay int

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
					network.ClientWriter(network.SrvAddrM,buf)
					idDelay++
				}
			}

		case msgd := <-channelP2P:

			if idDelay == msgd.Id {

				var delay = msgd.Time.Sub(tes)/2

				fmt.Println(delay)

			}
		}
	


	}
}
