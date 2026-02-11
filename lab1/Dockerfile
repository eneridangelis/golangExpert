FROM golang:1.22.0 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

From scratch
WORKDIR /app
COPY --from=builder /app/cloudrun .
ENTRYPOINT ["./cloudrun"]