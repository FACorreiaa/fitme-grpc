FROM golang:1.23.1 AS builder

LABEL maintainer="a11199"
LABEL description="Base image fitme dev"

WORKDIR /app

ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/fitme -ldflags="-w -s" ./*.go
ENTRYPOINT ["/app/fitme"]

FROM alpine:latest AS dev
WORKDIR /app
RUN apk add --no-cache bash
COPY --from=builder /app/fitme /usr/bin/fitme
EXPOSE 8000
EXPOSE 8001
CMD ["fitme", "start"]
