all: build

DEVICE = atmega328p
PORT = /dev/ttyACM0
COMPILER = avr-gcc -mmcu=$(DEVICE)
COMPILER_CLANG = clang --target=avr -I/usr/avr/include -mmcu=$(DEVICE)
PROGRAMMER = avrdude -c arduino -p $(DEVICE) -P $(PORT)

build:
	mkdir -p out
	$(COMPILER) main.c -o out/main.elf
	avr-objcopy -j .text -j .data -O ihex out/main.elf out/main
	rm out/main.elf

install: build
	$(PROGRAMMER) -U flash:w:out/main:i

build-clang:
	mkdir -p out
	$(COMPILER_CLANG) main.c -o out/main.elf
	avr-objcopy -j .text -j .data -O ihex out/main.elf out/main
	rm out/main.elf

install-clang: build-clang
	$(PROGRAMMER) -U flash:w:out/main:i

clean:
	rm -rf out