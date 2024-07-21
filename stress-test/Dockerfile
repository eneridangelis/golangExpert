# Use a imagem base do Go
FROM golang:1.21.5

# Crie e defina o diretório de trabalho
WORKDIR /app

# Copie os arquivos do projeto
COPY . .

# Compile o binário
RUN go build -o loadtest

# Defina o comando padrão
ENTRYPOINT ["/app/loadtest"]
