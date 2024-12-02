FROM golang:1.23.1 AS builder

LABEL maintainer="a11199"
LABEL description="Base image fitme dev"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/fitme ./*.go
ENTRYPOINT ["/app/fitme"]

FROM base AS reload
## Install the air binary so we get live code-reloading when we save files
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# Run the air command in the directory where our code will live
WORKDIR /app
CMD ["air"]

FROM alpine:latest AS dev
WORKDIR /app
RUN apk add --no-cache bash
COPY --from=builder /app/fitme /usr/bin/fitme
EXPOSE 8000
EXPOSE 8001
CMD ["fitme", "start"]
