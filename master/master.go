package main

import (
	"Laboratoire01/network"
	"bytes"
	"encoding/gob"
	"time"
)



func main() {

	go network.ClientReader(network.SrvAddr)

	for id:= 0;;{
		id +=1


		var buf bytes.Buffer

		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: time.Now(), Msg: "SYNC"}); err != nil {
			// handle error
		}
		tMaitre := time.Now()
		network.ClientWriter(network.MulticastAddr,buf)
		buf.Reset()
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: tMaitre, Msg: "FOLLOW_UP"}); err != nil {
			// handle error
		}
		network.ClientWriter(network.MulticastAddr,buf)
		time.Sleep(2 * time.Second)
	}
	
}

