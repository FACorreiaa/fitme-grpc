# Build stage
FROM golang:1.23.1 AS base

LABEL maintainer="a11199"
LABEL description="Base image fitme dev"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0

RUN go build -o /fitme ./*.go

# Final stage
FROM busybox

COPY --from=base /fitme /usr/bin/fitme
CMD ["/usr/bin/fitme"]
