package message

import "../orders"
import . "../driver"

const (
	Ping              = 1
	Elev_move         = 2
	New_order         = 3
	Order_assigned    = 4
	Order_completed   = 5
	Elev_delete       = 6
	Elev_add          = 7
	Elev_state_update = 8
	Elev_dead         = 9
)

type UDPMessage struct {
	MessageId int
	Source    int
	Target    int
	Floor     int
	Dir       Elev_dir
	Order     orders.External_order
}
