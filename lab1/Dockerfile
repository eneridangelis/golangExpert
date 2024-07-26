FROM golang:1.22.0 
WORKDIR /app
COPY . .

RUN go build -o cloudrun

EXPOSE 8080

CMD ["/app/cloudrun"]