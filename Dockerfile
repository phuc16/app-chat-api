FROM golang:1.21-alpine AS builder

RUN apk update && apk add git

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v --trimpath -o app -ldflags="-X 'main.buildVersion=$(git rev-parse HEAD)' -X 'main.buildDate=$(date)'" main.go

FROM ubuntu:20.04

WORKDIR /app
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait ./wait
COPY ./config.yaml ./config.yaml
COPY --from=builder /app/app ./app
RUN chmod +x ./wait
RUN chmod +x ./app
CMD ./wait && ./app
