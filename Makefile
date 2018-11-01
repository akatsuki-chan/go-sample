EXE = sample
EXE_LINUX = $(EXE)
EXE_OSX = $(EXE)-osx

BIN_LINUX = build/$(EXE_LINUX)
BIN_OSX = build/$(EXE_OSX)

IMAGE_NAME = go-sample:latest

ENV_MODE = dev

prod: ENV_MODE=prod
prod: docker-build

docker-run: $(BIN_LINUX)
	docker run --rm -e ENV_MODE=$(ENV_MODE) -t $(IMAGE_NAME)

docker-build: $(BIN_LINUX)
	docker build --build-arg ENV_MODE=$(ENV_MODE) -t $(IMAGE_NAME) .

$(BIN_OSX): vendor
	GOOS=darwin GOARCH=amd64 go build -o $@

$(BIN_LINUX): vendor
	GOOS=linux GOARCH=amd64 go build -o $@

vendor:
	dep ensure

clean:
	-rm -rf build
	-rm -rf vendor
	-docker rmi $(IMAGE_NAME)