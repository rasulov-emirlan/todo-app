build:
	go generate ./...
	go build -o bin/apiserver cmd/api/main.go

dev:
	go generate ./...
	go build -o bin/apiserver cmd/api/main.go
	./bin/apiserver -config .env -migrations true

compose:
	docker-compose build && docker-compose up