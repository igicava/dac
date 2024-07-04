FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o harmonist orchestrator/cmd/calc-service/main.go

EXPOSE 8080 8081

CMD ["/app/harmonist"]
