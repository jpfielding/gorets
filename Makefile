all: test

test:
	go test -v

clean:
	rm *.test
