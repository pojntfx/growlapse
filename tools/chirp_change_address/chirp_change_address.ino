#include <Wire.h>

#define OLD_ADDRESS 0x20
#define NEW_ADDRESS 0x31

void writeI2CRegister8bit(int addr, int reg, int value) {
  Wire.beginTransmission(addr);
  Wire.write(reg);
  Wire.write(value);
  Wire.endTransmission();
}

void writeI2CRegister8bit(int addr, int value) {
  Wire.beginTransmission(addr);
  Wire.write(value);
  Wire.endTransmission();
}

void setup() {
  Wire.begin();
  Serial.begin(9600);

  Serial.println("Starting address change ...");

  writeI2CRegister8bit(OLD_ADDRESS, 6); // Reset the Chirp!
  delay(1000);                          // Wait until reboot
  writeI2CRegister8bit(
      OLD_ADDRESS, 1,
      NEW_ADDRESS); // Change from default address 0x20 to new address 0x32
  writeI2CRegister8bit(OLD_ADDRESS, 6); // Reset the Chirp!
  delay(1000);                          // Wait until reboot

  Serial.println("Address change done! Press the button on the Chirp! to use "
                 "the new address.");
}

void loop() {}
