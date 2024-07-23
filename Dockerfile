FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o bin/mono_pardo ./cmd/.

EXPOSE 8000

CMD ["./bin/mono_pardo"]