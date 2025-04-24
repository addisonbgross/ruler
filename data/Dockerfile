FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod /app
RUN go mod download

COPY . /app

RUN go build -o ruler-data .

RUN chmod +x /app/ruler-data

EXPOSE 8081

CMD ["./ruler-data"]