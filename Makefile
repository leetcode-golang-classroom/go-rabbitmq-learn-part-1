.PHONY=build

build-send:
	@go build -o bin/send send/send.go
build-receive:
	@go build -o bin/receive receive/receive.go

run-send: build-send
	@./bin/send
run-receive: build-receive
	@./bin/receive
test:
	@go test -v -cover ./test/...