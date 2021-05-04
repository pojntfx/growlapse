all: build

DEVICE = attiny44
COMPILER = avr-gcc -DF_CPU="1000000" -mmcu=$(DEVICE) -Os
PROGRAMMER = avrdude -c usbasp -P usb -p $(DEVICE)

build:
	mkdir -p out
	$(COMPILER) main.c usiTwiSlave.c -o out/main.elf
	avr-objcopy -j .text -j .data -O ihex out/main.elf out/main
	rm out/main.elf

install: build
	$(PROGRAMMER) -U hfuse:w:0xDE:m -U lfuse:w:0xE2:m -B 20.96
	$(PROGRAMMER) -U flash:w:out/main:i

clean:
	rm -rf out