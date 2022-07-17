FROM golang:alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY ./ ./

RUN go build -o apiserver cmd/api/main.go

FROM alpine:latest

COPY --from=0 /app/apiserver /apiserver

EXPOSE 8080

COPY .env .env

ENTRYPOINT ["/apiserver"]