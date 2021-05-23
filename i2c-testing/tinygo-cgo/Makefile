all: build

DEVICE = arduino

build:
	mkdir -p out
	tinygo build -target arduino -o out/main.hex .

install: build
	tinygo flash -target arduino .

clean:
	rm -rf out