#include <Wire.h>

#define SENSOR_ADDRESS 0x20

void setup() {
  Wire.begin();
  Serial.begin(9600);
}

void writeI2CRegister8bit(int addr, int value) {
  Wire.beginTransmission(addr);
  Wire.write(value);
  Wire.endTransmission();
}

unsigned int readI2CRegister16bit(int addr, int reg) {
  Wire.beginTransmission(addr);
  Wire.write(reg);
  Wire.endTransmission();
  delay(1100);
  Wire.requestFrom(addr, 2);
  unsigned int t = Wire.read() << 8;
  t = t | Wire.read();
  return t;
}

void loop() {
  Serial.print("Capacitance: ");
  Serial.println(readI2CRegister16bit(SENSOR_ADDRESS,
                                      0)); // Read from capacitance register

  writeI2CRegister8bit(0x32, 3); // Request light measurement
  delay(9000);                   // Wait until light measurement is done
  Serial.print("Light: ");
  Serial.println(
      readI2CRegister16bit(SENSOR_ADDRESS, 4)); // Read from light register
}
