FROM golang

WORKDIR /app

COPY go.mod /app

COPY go.sum /app

RUN go mod download

COPY . .

RUN go build .\main.go

EXPOSE 4000:4000

CMD ["./main"]