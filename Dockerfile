FROM golang:1.20.7-alpine3.18 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app/api

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/api/main.go

FROM alpine:3.18

WORKDIR /app/bin

COPY --from=builder /app/api/bin .
COPY --from=builder /app/api/.env .

ENTRYPOINT [ "./api" ]