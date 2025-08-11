FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/message-processor

FROM alpine:3.20
WORKDIR /app
COPY --from=build /bin/message-processor /usr/local/bin/message-processor
ENTRYPOINT ["/usr/local/bin/message-processor"]
