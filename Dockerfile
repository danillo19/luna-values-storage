# syntax=docker/dockerfile:1

FROM golang:1.19.6-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o api

EXPOSE 8111

CMD ["sh", "./docker-entrypoint.sh"]
