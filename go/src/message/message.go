package message

import "../orders"

//Message ID
const (
	Ping            = 1
	Elev_move       = 2
	New_order       = 3
	Elev_delete     = 4
	Elev_add        = 5
	Calc_cost       = 6
	Cost_calculated = 7
)

type UDPMessage struct {
	MessageId int
	Source    int
	Target    int
	Order     orders.External_order
}

/*func CalculateChecksum(Msg *UDPMessage) int { // not a very good crc, just for testing
	c := Msg.MessageId%7 + Msg.OrderQueue[0]%7 + Msg.ElevatorStateUpdate[0]%7
	return c
}*/
