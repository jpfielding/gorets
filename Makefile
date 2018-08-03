all: restore-deps test

.PHONY: explorer
explorer:
	CGO_ENABLED=0 GOOS=linux go build -o bin/explorer/explorer cmds/explorer/*.go
	cd explorer/client && npm install && CONFIG_ENV=docker npm run build
	cp docker/explorer/Dockerfile bin/explorer
	cp cmds/explorer/config.json bin/explorer
	docker build bin/explorer -t "gorets_explorer:latest"

test-explorer:
		CGO_ENABLED=0 GOOS=linux go build -o bin/explorer/explorer cmds/explorer/*.go
		cd explorer/client && npm install && CONFIG_ENV=test npm run build
		cp docker/explorer/Dockerfile bin/explorer
		docker build bin/explorer -t "gorets_explorer_test:latest"

vendor:
	glide up
test:
	go test -v $$(glide novendor)
frontend-test:
	make test-explorer
	docker-compose up -d
	-cd explorer/client/ && npm run test
	docker-compose down
vet:
	go vet $$(glide novendor)
clean:
	rm -rf bin *.test
restore-deps:
	@command -v glide >/dev/null 2>&1 || { echo >&2 "Error: glide (https://github.com/Masterminds/glide) is not installed.  Please install.  Aborting."; exit 1; }
	rm -rf vendor/ glide.lock
	glide up
