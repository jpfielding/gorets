
CONFIG_ENV ?= docker

all: restore-deps test

.PHONY: explorer
explorer: build-explorer build-explorer-copy-files build-explorer-container

build-explorer: build-explorer-svc build-explorer-client

build-explorer-svc:
	CGO_ENABLED=0 GOOS=linux go build -o bin/explorer/explorer cmd/explorer/*.go

build-explorer-client:
	cd web/explorer/client && npm install && CONFIG_ENV=$(CONFIG_ENV) npm run build

build-explorer-copy-files:
	cp build/explorer/Dockerfile bin/explorer
	cp cmd/explorer/config.json bin/explorer

build-explorer-container:
	docker build bin/explorer -t "gorets_explorer:latest"

test-explorer:
	CONFIG_ENV=test make build-explorer-svc build-explorer-client build-explorer-copy-files
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
