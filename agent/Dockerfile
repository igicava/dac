FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o agents agent/cmd/agent/main.go

CMD ["/app/agents"]
