package network

import (
	. "../UDP"
	. "../message"
	"fmt"
	"net"
	"time"
)

var elev_timer map[int]*time.Timer

func broadcast_ip(id int, channel chan UDPMessage) {
	id_message := UDPMessage{Source: id, MessageId: Ping}
	for {
		channel <- id_message
		time.Sleep(100 * time.Millisecond)
	}
}

func Delete_elev(id int, to_SM chan UDPMessage) {
	delete(elev_timer, id)
	to_SM <- UDPMessage{Source: id, MessageId: Elev_delete}

}

func Network_manager(from_SM chan UDPMessage, to_SM chan UDPMessage) {
	addr, _ := net.InterfaceAddrs()
	self_id := int(addr[1].String()[12]-'0')*100 + int(addr[1].String()[13]-'0')*10 + int(addr[1].String()[14]-'0')
	UDP_send := make(chan UDPMessage, 100)
	UDP_receive := make(chan UDPMessage, 100)
	go broadcast_ip(self_id, UDP_send)
	go UDP_sender(UDP_send)
	go UDP_listener(UDP_receive)

	elev_timer = make(map[int]*time.Timer)

	for {
		select {
		case message := <-UDP_receive:
			if message.MessageId == Ping {
				_, elev_present := elev_timer[message.Source]
				if message.Source != self_id {
					if elev_present {
						elev_timer[message.Source].Reset(time.Second)
					} else {
						elev_timer[message.Source] = time.AfterFunc(time.Second, func() { Delete_elev(message.Source, to_SM) })
						to_SM <- UDPMessage{MessageId: Elev_add, Source: message.Source}
					}
				}
			}
			to_SM <- message

		case message := <-from_SM:
			fmt.Println("Network manager received from SM")
			message.Source = self_id
			UDP_send <- message
		}
	}
}
