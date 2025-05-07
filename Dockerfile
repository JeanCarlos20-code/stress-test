FROM golang:1.22.4-alpine

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o load-tester cmd/main.go

ENTRYPOINT ["./load-tester"]