FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /bin/responder

FROM alpine:3.20
WORKDIR /app
COPY --from=build /bin/responder /usr/local/bin/responder
ENTRYPOINT ["/usr/local/bin/responder"]
