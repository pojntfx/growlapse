#include <Wire.h>

void setup() {
  Wire.begin();
  Serial.begin(9600);
}

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

void loop() {
  Serial.println("Starting scan ...");

  int deviceCount = 0;

  for (byte address = 1; address < 127; address++) {
    Wire.beginTransmission(address);
    byte err = Wire.endTransmission();

    if (err == 0) {
      Serial.print("I2C device found at 0x");
      if (address < 16)
        Serial.print("0");
      Serial.print(address, HEX);
      Serial.println("  !");

      deviceCount++;
    } else if (err == 4) {
      Serial.print("Unknown error at address 0x");
      if (address < 16)
        Serial.print("0");
      Serial.println(address, HEX);
    }
  }

  if (deviceCount == 0)
    Serial.println("Scan done; no I2C devices found.");
  else
    Serial.println("Scan done.");

  delay(500);
}
