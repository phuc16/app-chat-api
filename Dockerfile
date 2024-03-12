FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-chat-api

COPY ./wait /wait

EXPOSE 8080

# Run
CMD chmod +x /wait && /wait && /app-chat-api