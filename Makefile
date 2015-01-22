all: test

test:
	go test -v client/*.go
	go test -v metadata/*.go

clean:
	rm *.test
