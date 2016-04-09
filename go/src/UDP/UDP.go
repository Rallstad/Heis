package UDP

import (
	. "../message"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	UDP_port = ":20014"
)

func UDP_sender(channel chan UDPMessage) {
	broadcast_addr := []string{"129.241.187.255", UDP_port}
	broadacst_udp, _ := net.ResolveUDPAddr("udp", strings.Join(broadcast_addr, ""))
	broadcast_connection, _ := net.DialUDP("udp", nil, broadacst_udp)
	defer broadcast_connection.Close()
	for {
		buf, err := json.Marshal(<-channel)
		if err == nil {
			broadcast_connection.Write(buf)
		}
	}
}

func UDP_listener(channel chan UDPMessage) {
	udp_receive_addr, err := net.ResolveUDPAddr("udp", UDP_port)
	if err != nil {
		fmt.Println(err)
	}
	udp_connection, err := net.ListenUDP("udp", udp_receive_addr)
	if err != nil {
		fmt.Println(err)
	}
	defer udp_connection.Close()

	buf := make([]byte, 2048)
	trimmed_buf := make([]byte, 1)
	var received_message UDPMessage
	for {
		n, _, _ := udp_connection.ReadFromUDP(buf)
		trimmed_buf = buf[:n]
		err := json.Unmarshal(trimmed_buf, &received_message)
		if err == nil {
			channel <- received_message
		}
	}
}
