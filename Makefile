all: restore-deps test vet

vendor:
	glide up --strip-vendor
test:
	go test -v $$(glide novendor)
vet: 
	go vet $$(glide novendor)
clean:
	rm *.test
restore-deps:
	@command -v glide >/dev/null 2>&1 || { echo >&2 "Error: glide (https://github.com/Masterminds/glide) is not installed.  Please install.  Aborting."; exit 1; }
	rm -rf vendor/
	glide up --strip-vendor
