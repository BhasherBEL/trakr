.PHONY: build run clean

all: clean build run

build:
	mkdir -p bin
	go build -o ./bin/trakr .

run:
	./bin/trakr

clean:
	go clean
	rm -rf bin