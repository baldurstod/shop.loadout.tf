.PHONY: build clean

BINARY_NAME=shop.loadout.tf

build:
	go env -w CGO_ENABLED=1
	go build -tags cse -ldflags="-X shop.loadout.tf/src/server/server.ReleaseMode=false" -o dist/${BINARY_NAME} ./src/server/

run: build
	dist/${BINARY_NAME}

prod:
	go env -w CGO_ENABLED=1
	@echo 'Bundling shop.loadout.tf'
	rollup -c --environment BUILD:production
	@echo 'Building go app'
	go build -tags cse -o dist/${BINARY_NAME} ./src/server/

clean:
	go clean
