package elev

import (
	//"math"
	. "../driver"
	. "../io"
)

const N_FLOORS = 4
const N_BUTTON_TYPES = 3
const MOTOR_SPEED = 2800

const (
	BUTTON_CALL_UP   = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND   = 2
)

var lamp_matrix = [N_FLOORS][N_BUTTON_TYPES]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_matrix = [N_FLOORS][N_BUTTON_TYPES]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Elev_init() {
	init_success := IO_init() != 0
	if init_success {
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTON_TYPES; b++ {
				lamp_matrix[f][b] = 0
				button_matrix[f][b] = 0
			}
		}

	}
}

func elev_set_button_lamp(button int, floor int, value int) {
	if value != 0 {
		IO_set_bit(lamp_matrix[floor][button])
	} else {
		IO_clear_bit(lamp_matrix[floor][button])
	}
}

func Set_motor_dir(dir int) {
	if dir == 0 {
		IO_write_analog(MOTOR, 0)
	} else if dir > 0 {
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	} else if dir < 0 {
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	}

}

func elev_set_floor_indicator(floor int) {
	switch floor {
	case 0:
		IO_clear_bit(LIGHT_FLOOR_IND1)
		IO_clear_bit(LIGHT_FLOOR_IND2)
	case 1:
		IO_clear_bit(LIGHT_FLOOR_IND1)
		IO_set_bit(LIGHT_FLOOR_IND2)
	case 2:
		IO_set_bit(LIGHT_FLOOR_IND1)
		IO_clear_bit(LIGHT_FLOOR_IND2)
	case 3:
		IO_set_bit(LIGHT_FLOOR_IND1)
		IO_set_bit(LIGHT_FLOOR_IND2)
	}

}

func elev_set_door_open_lamp(value int) {
	if value != 0 {
		IO_set_bit(LIGHT_DOOR_OPEN)
	} else {
		IO_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func elev_set_stop_lamp(value int) {
	if value != 0 {
		IO_set_bit(LIGHT_STOP)
	} else {
		IO_clear_bit(LIGHT_STOP)
	}
}

func elev_get_button_signal(button int, floor int) int {
	if IO_read_bit(button_matrix[floor][button]) {
		return 1
	} else {
		return 0
	}
}

func elev_get_floor_sensor_signal() int {
	if IO_read_bit(SENSOR_FLOOR1) {
		return 0
	} else if IO_read_bit(SENSOR_FLOOR2) {
		return 1
	} else if IO_read_bit(SENSOR_FLOOR3) {
		return 2
	} else if IO_read_bit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}
