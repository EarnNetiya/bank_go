FROM golang:1.24.1-alpine

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o bankapp .

EXPOSE 8080

CMD ["./bankapp"]
