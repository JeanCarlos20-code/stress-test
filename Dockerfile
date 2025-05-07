FROM golang:1.22.4-alpine

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o stress-test cmd/main.go

ENTRYPOINT ["./stress-test"]