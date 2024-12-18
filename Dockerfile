FROM golang:1.23.1 AS builder

LABEL maintainer="a11199"
LABEL description="Base image fitme dev"

WORKDIR /app

ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN echo "Contents of /app after COPY:" && ls -al /app && sleep 1
RUN echo "Listing contents of /app after COPY:" && ls -al /app && sleep 1
RUN echo "Contents of /app/config:" && ls -al /app/config || echo "/app/config does not exist" && sleep 1
RUN ls -al /app/internal || echo "No /app/internal found" && sleep 1
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/fitme -ldflags="-w -s" ./*.go
ENTRYPOINT ["/app/fitme"]

FROM alpine:latest AS dev
WORKDIR /app

RUN apk add --no-cache bash
COPY --from=builder /app/fitme /usr/bin/fitme
COPY --from=builder /app/config ./config

RUN echo "Contents of /app after COPY:" && ls -al /app && sleep 1
RUN echo "Listing contents of /app after COPY:" && ls -al /app && sleep 1
RUN echo "Contents of /app/config:" && ls -al /app/config || echo "/app/config does not exist" && sleep 1
RUN ls -al /app/internal || echo "No /app/internal found" && sleep 1

EXPOSE 8000
EXPOSE 8001
CMD ["fitme", "start"]
