FROM golang:1.22-alpine
LABEL authors="bruno"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o /app/main .

EXPOSE 8080

CMD ["/app/main"]