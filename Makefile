.PHONY: bin
bin:
	mkdir -p bin

.PHONY: test
test:
	go test ./...

.PHONY: build
build: test bin
	go fmt ./...
	go vet ./...
	golint ./cmd ./gcode
	go build -o bin/gcode-analyzer ./cmd/gcode-analyzer/*

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
