#include "i2c.h"

// Based on https://github.com/thegouger/avr-i2c-slave/blob/master/I2CSlave.c

void i2c_init(uint8_t address) {
  cli();
  // load address into TWI address register
  TWAR = address << 1;
  // set the TWCR to enable address matching and enable TWI, clear TWINT, enable
  // TWI interrupt
  TWCR = (1 << TWIE) | (1 << TWEA) | (1 << TWINT) | (1 << TWEN);
  sei();
  
}

void i2c_stop() {
  // clear acknowledge and enable bits
  cli();
  TWCR = 0;
  TWAR = 0;
  sei();
}
