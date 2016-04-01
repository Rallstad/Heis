#include <comedilib.h>

#include "io.h"



static comedi_t *it_g = NULL;



//in port 4
#define PORT_4_SUBDEVICE        3
#define PORT_4_CHANNEL_OFFSET   16
#define PORT_4_DIRECTION        COMEDI_INPUT
#define OBSTRUCTION             (0x300+23)
#define STOP                    (0x300+22)
#define BUTTON_COMMAND1         (0x300+21)
#define BUTTON_COMMAND2         (0x300+20)
#define BUTTON_COMMAND3         (0x300+19)
#define BUTTON_COMMAND4         (0x300+18)
#define BUTTON_UP1              (0x300+17)
#define BUTTON_UP2              (0x300+16)

//in port 1
#define PORT_1_SUBDEVICE        2
#define PORT_1_CHANNEL_OFFSET   0
#define PORT_1_DIRECTION        COMEDI_INPUT
#define BUTTON_DOWN2            (0x200+0)
#define BUTTON_UP3              (0x200+1)
#define BUTTON_DOWN3            (0x200+2)
#define BUTTON_DOWN4            (0x200+3)
#define SENSOR_FLOOR1           (0x200+4)
#define SENSOR_FLOOR2           (0x200+5)
#define SENSOR_FLOOR3           (0x200+6)
#define SENSOR_FLOOR4           (0x200+7)

//out port 3
#define PORT_3_SUBDEVICE        3
#define PORT_3_CHANNEL_OFFSET   8
#define PORT_3_DIRECTION        COMEDI_OUTPUT
#define MOTORDIR                (0x300+15)
#define LIGHT_STOP              (0x300+14)
#define LIGHT_COMMAND1          (0x300+13)
#define LIGHT_COMMAND2          (0x300+12)
#define LIGHT_COMMAND3          (0x300+11)
#define LIGHT_COMMAND4          (0x300+10)
#define LIGHT_UP1               (0x300+9)
#define LIGHT_UP2               (0x300+8)

//out port 2
#define PORT_2_SUBDEVICE        3
#define PORT_2_CHANNEL_OFFSET   0
#define PORT_2_DIRECTION        COMEDI_OUTPUT
#define LIGHT_DOWN2             (0x300+7)
#define LIGHT_UP3               (0x300+6)
#define LIGHT_DOWN3             (0x300+5)
#define LIGHT_DOWN4             (0x300+4)
#define LIGHT_DOOR_OPEN         (0x300+3)
#define LIGHT_FLOOR_IND2        (0x300+1)
#define LIGHT_FLOOR_IND1        (0x300+0)

//out port 0
#define MOTOR                   (0x100+0)

//non-existing ports (for alignment)
#define BUTTON_DOWN1            -1
#define BUTTON_UP4              -1
#define LIGHT_DOWN1             -1
#define LIGHT_UP4               -1











int io_init(void) {

    it_g = comedi_open("/dev/comedi0");

    if (it_g == NULL) {
        return 0;
    }

    int status = 0;
    for (int i = 0; i < 8; i++) {
        status |= comedi_dio_config(it_g, PORT_1_SUBDEVICE, i + PORT_1_CHANNEL_OFFSET, PORT_1_DIRECTION);
        status |= comedi_dio_config(it_g, PORT_2_SUBDEVICE, i + PORT_2_CHANNEL_OFFSET, PORT_2_DIRECTION);
        status |= comedi_dio_config(it_g, PORT_3_SUBDEVICE, i + PORT_3_CHANNEL_OFFSET, PORT_3_DIRECTION);
        status |= comedi_dio_config(it_g, PORT_4_SUBDEVICE, i + PORT_4_CHANNEL_OFFSET, PORT_4_DIRECTION);
    }

    return (status == 0);
}


void io_set_bit(int channel) {
    comedi_dio_write(it_g, channel >> 8, channel & 0xff, 1);
}


void io_clear_bit(int channel) {
    comedi_dio_write(it_g, channel >> 8, channel & 0xff, 0);
}


void io_write_analog(int channel, int value) {
    comedi_data_write(it_g, channel >> 8, channel & 0xff, 0, AREF_GROUND, value);
}


int io_read_bit(int channel) {
    unsigned int data = 0;
    comedi_dio_read(it_g, channel >> 8, channel & 0xff, &data);
    return (int)data;
}


int io_read_analog(int channel) {
    lsampl_t data = 0;
    comedi_data_read(it_g, channel >> 8, channel & 0xff, 0, AREF_GROUND, &data);
    return (int)data;
}

