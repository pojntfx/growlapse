#include <avr/interrupt.h>
#include <util/twi.h>

#define I2C_ADDRESS 0x10

int main() {
  asm volatile("cli" ::: "memory");

  TWAR = I2C_ADDRESS << 1;
  TWCR |= 1 << TWIE;
  TWCR |= 1 << TWEA;
  TWCR |= 1 << TWINT;
  TWCR |= 1 << TWEN;

  asm volatile("sei" ::: "memory");

  while (1) {
  }
}