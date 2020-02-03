main: test.o 
	go tool objdump test.o

test.o: test.go
	go tool compile -S -N -l test.go
