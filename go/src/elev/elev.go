package elev

import(
	"math"
	"../io"
);

const N_FLOORS = 4;
const N_BUTTON_TYPES = 3;
const MOTOR_SPEED = 2800;



const(
	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2
)

var lamp_matrix = [N_FLOORS][N_BUTTON_TYPES]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_matrix = [N_FLOORS][N_BUTTONS]int{
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
};

func elev_init(){
	 init_success := io_init()
	 if init_success{
	 	for ( int f = 0; f < N_FLOORS; f++){
	 		for (int b = 0; b < N_BUTTONS; b++){
	 			lamp_matrix[f]][b] = 0;
	 			button_matrix[f][b] = 0;
	 		}
	 	}

	 }
}

func elev_set_button_lamp(button int, floor int, value int){
	if value {
		io_set_bit(lamp_matrix[floor][button])
	}
	else {
		io_clear_bit(lamp_matrix[floor][button])
	}
}

func set_motor_dir(dir int){
	if dir = 0 {
		io_write_analog(MOTOR,0)
	}
	else if dir > 0 {
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR,MOTOR_SPEED)
	}
	else if dir < 0 {
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR,MOTOR_SPEED)
	}

}



func elev_set_button_lamp(button int, floor int, value int){
	if value{
		io_set_bit(lamp_matrix[floor][button])

	}
	else{
		io_clear_bit(lamp_matrix[floor][button])
	}
}



func elev_set_floor_indicator(floor int){
	if (floor & 0x02){
		io_set_bit(LIGHT_FLOOR_IND1)
	}
	else{
		io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if (floor & 0x01){
		io_set_bit(LIGHT_FLOOR_IND2)
	}
	else{
		io_clear_bit(LIGHT_FLOOR_IND2)
	}

}

func elev_set_door_open_lamp(value int){
	if value{
		io_set_bit(LIGHT_DOOR_OPEN)
	}
	else{
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func elev_set_stop_lamp(value int){
	if value{
		io_set_bit(LIGHT_STOP)
	}
	else{
		io_clear_bit(LIGHT_STOP)
	}
}


func elev_get_button_signal(button int, floor int) int{
	if io_read_bit(button_matrix[floor][button]){
		return 1
	}
	else{
		return 0
	}
}


func elev_get_floor_sensor_signal() int{
	if io_read_bit(SENSOR_FLOOR1){
		return 0
	}
	else if io_read_bit(SENSOR_FLOOR2){
		return 1
	}
	else if io_read_bit(SENSOR_FLOOR3){
		return 2
	}
	else if io_read_bit(SENSOR_FLOOR4){
		return 3
	}
	else{
		return -1
	}
}

