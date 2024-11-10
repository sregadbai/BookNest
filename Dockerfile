ARG ENV=local

FROM golang:1.22 AS builder

ARG ENV

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY $ENV.env /app/$ENV.env

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o booknest-app

FROM alpine:latest
ARG ENV
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/booknest-app .
COPY --from=builder /app/$ENV.env .

EXPOSE 8080

CMD ["./booknest-app"]