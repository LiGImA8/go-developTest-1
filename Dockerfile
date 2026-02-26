FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/user-service ./services/user/cmd
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/order-service ./services/order/cmd
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/log-service ./services/log/cmd

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates libzmq5 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /bin/user-service /usr/local/bin/user-service
COPY --from=builder /bin/order-service /usr/local/bin/order-service
COPY --from=builder /bin/log-service /usr/local/bin/log-service
WORKDIR /srv
CMD ["/usr/local/bin/user-service"]
