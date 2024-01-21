FROM golang:1.21-alpine as builder

RUN apk update

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/server/server.go

FROM alpine:latest

RUN apk update

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/config /app/config

EXPOSE 8080

RUN adduser -D -g '' appuser
USER appuser

CMD ["./app"]