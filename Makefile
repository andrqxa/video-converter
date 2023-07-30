BINARY_NAME=video-converter

## build: Build binary
build:
	@echo "Building..."
ifeq ($(OS),Windows_NT)
	go build -ldflags="-s -w" -o ${BINARY_NAME}.exe .
else
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} .
endif
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
ifeq ($(OS),Windows_NT)
	.\${BINARY_NAME}.exe &
else
	./${BINARY_NAME} &
endif
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
ifeq ($(OS),Windows_NT)
	@del ${BINARY_NAME}.exe
else
	@rm ${BINARY_NAME}
endif
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
ifeq ($(OS),Windows_NT)
	@taskkill /F /IM ${BINARY_NAME}.exe
else
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
endif
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	go test -v ./...