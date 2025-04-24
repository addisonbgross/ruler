FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod /app
RUN go mod download

COPY . /app

RUN go build -o ruler-node .

EXPOSE 8080

CMD ["./ruler-node"]