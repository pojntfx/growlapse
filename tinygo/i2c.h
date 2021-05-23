#ifndef I2C_H
#define I2C_H

#include <avr/interrupt.h>
#include <stdint.h>
#include <util/twi.h>

void i2c_init(uint8_t addr);

void i2c_stop();

#endif