build:
	go build -o poc-aws-websocket-gateway ./cmd/server/main.go

test:
	go test ./...

clean:
	go clean
	rm -f scale-websocket

run-server:
	go run ./cmd/main.go

get-sign:
	go test -v ./services/authsvc