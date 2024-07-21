# Use a imagem base do Go
FROM golang:1.20 AS builder

# Defina o diretório de trabalho
WORKDIR /app

# Copie o código-fonte para o contêiner
COPY . .

# Compile o aplicativo
RUN go build -o app .

# Use uma imagem base menor para o runtime
FROM alpine:latest

# Instale ca-certificates (necessário para conexões HTTPS)
RUN apk --no-cache add ca-certificates

# Copie o binário compilado do estágio anterior
COPY --from=builder /app/app /app/app

# Defina o ponto de entrada
ENTRYPOINT ["/app/app"]

# Exponha a porta que o serviço usa
EXPOSE 8080
