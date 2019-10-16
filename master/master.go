package main

import (
	"../network"
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
)



func main(){

	channel := make(chan network.MessageWithOrigin)

	go network.ClientReader(network.SrvAddrMaster, channel)

	go  synchronization()

	var buf bytes.Buffer
	for {
		select {
		case msgWithOrigin := <-channel:
			{
				msg:=msgWithOrigin.Msg
				buf.Reset()
				slaveAddr:= msgWithOrigin.Ip+":"+strconv.Itoa(msg.OriginPort)
				fmt.Println("Received DELAY_REQUEST from ",slaveAddr)
				err := gob.NewEncoder(&buf).Encode(network.Message{Id: msg.Id, Time: time.Now(), Msg: 0b00});
				if err != nil {
					// handle error
				}

				network.ClientWriter(slaveAddr,buf)
			}
		}
	}
}

func synchronization()  {

		var buf bytes.Buffer

		for id:= 1;;{

		id +=1
		buf.Reset()
		//fmt.Print(time.Time{}.String())
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: time.Time{}, Msg: 0b00}); err != nil {
			// handle error
		}

		masterTime := time.Now()
		network.ClientWriter(network.MulticastAddr,buf)
		buf.Reset()
		if err := gob.NewEncoder(&buf).Encode(network.Message{ Id: id, Time: masterTime, Msg: 0b01}); err != nil {
			// handle error
		}
		network.ClientWriter(network.MulticastAddr,buf)
		time.Sleep(2 * time.Second)
	}
}


