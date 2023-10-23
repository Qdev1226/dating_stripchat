FROM golang:1.14-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
# COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/build/api /app/cmd

CMD ["/app/build/api"]