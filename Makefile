run:
	go build -o bin/http-server ./cmd/ && bin/http-server

test:
	go test ./...