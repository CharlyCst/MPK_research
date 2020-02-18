main:
	gosb build
	./test

test: test.o 
	go tool objdump test.o > x86_64
	go build

test.o: test.go
	go tool compile -S -N -l test.go > plan9