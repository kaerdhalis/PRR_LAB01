package main

import (
	"Laboratoire01/network"
	"bytes"
	"encoding/gob"
	"time"
)



func main() {

	go network.ClientReader(network.SrvAddr)

	synchronization()

	
}

func synchronization()  {

	var buf bytes.Buffer

	for id:= 0;;{
		id +=1




		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: time.Time{}, Msg: 0b10}); err != nil {
			// handle error
		}
		tMaitre := time.Now()
		network.ClientWriter(network.MulticastAddr,buf)
		buf.Reset()
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: tMaitre, Msg: 0b11}); err != nil {
			// handle error
		}
		network.ClientWriter(network.MulticastAddr,buf)
		time.Sleep(2 * time.Second)
	}
}


