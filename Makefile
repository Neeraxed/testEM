ENTRYPOINT=cmd/app/main.go
BINARY_NAME=testEM
BUILD_FOLDER=build

.PHONY: all run build clean docs start

all: clean build
run: clean build start

build:
	mkdir -p $(BUILD_FOLDER)
	cp -n .env.example $(BUILD_FOLDER)/.env
	go build -o $(BUILD_FOLDER)/$(BINARY_NAME) -v $(ENTRYPOINT)
clean:
	go clean
	rm -rf $(BUILD_FOLDER)/*
tools:
	go install github.com/swaggo/swag/cmd/swag@latest
docs:
	 swag init -d ./cmd/app --pd --pdl 2

start:
	./$(BUILD_FOLDER)/$(BINARY_NAME)
