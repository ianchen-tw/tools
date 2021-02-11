BIN:= bin/main

ALL: build


.PHONY: build
build:
	go build -o $(BIN) ./src/main.go

.PHONY: clean
clean:
	rm -f *.json *.yaml
	rm -rf bin

.PHONY: test
test:
	go test -v ./... --cover