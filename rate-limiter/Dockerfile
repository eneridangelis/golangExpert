From golang:1.21.5

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o rate-limiter

EXPOSE 8080

CMD ["/app/rate-limiter"]