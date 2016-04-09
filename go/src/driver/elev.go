package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
#include "elev.h"
*/
import "C"

import (
	"fmt"
)

type Elev_dir int
type button_type int

const (
	BUTTON_UP     = 0
	BUTTON_DOWN   = 1
	BUTTON_INSIDE = 2
)

const (
	DOWN      Elev_dir = -1
	STOPMOTOR Elev_dir = 0
	UP        Elev_dir = 1
)

const N_FLOOR = 4
const N_BUTTON_TYPES = 3

func Elev_init() {
	C.elev_init()
	for Elev_get_floor_sensor_signal() == -1 {
		Elev_set_motor_direction(DOWN)

	}
	Elev_set_motor_direction(STOPMOTOR)
	fmt.Println("--------------------------------------------------")
	fmt.Println("Initialized at floor ", Elev_get_floor_sensor_signal())
}

func Elev_set_motor_direction(dir Elev_dir) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dir))
}

func Elev_set_button_lamp(button int, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Elev_set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Elev_set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Elev_set_stop_lamp(value int) {
	C.elev_set_stop_lamp(C.int(value))

}

func Elev_get_button_signal(button int, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Elev_get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Elev_get_stop_signal() int {
	return int(C.elev_get_stop_signal())
}
