FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/telegram-gateway

FROM alpine:3.20
WORKDIR /app
COPY --from=build /bin/telegram-gateway /usr/local/bin/telegram-gateway
ENTRYPOINT ["/usr/local/bin/telegram-gateway"]
