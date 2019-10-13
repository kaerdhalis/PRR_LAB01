package main

import (
	"Laboratoire01/network"
	"bytes"
	"encoding/gob"
	"time"
)



func main(){

	channel := make(chan network.Message)

	go network.ClientReader(network.SrvAddrM, channel)

	 synchronization()

	
}

func synchronization()  {

	var buf bytes.Buffer

	for id:= 0;;{

		id +=1
		buf.Reset()
		//fmt.Print(time.Time{}.String())
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: time.Time{}, Msg: 0b00}); err != nil {
			// handle error
		}

		tMaitre := time.Now()
		network.ClientWriter(network.MulticastAddr,buf)
		buf.Reset()
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: tMaitre, Msg: 0b01}); err != nil {
			// handle error
		}
		network.ClientWriter(network.MulticastAddr,buf)
		time.Sleep(2 * time.Second)
	}
}


