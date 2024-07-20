# Rate Limiter

## Como rodar o projeto:

Use o Docker Compose para subir o Redis, rodando o comando abaixo na pasta raiz:

```sh
docker-compose up -d
```

Instale as dependências do projeto com o seguinte comando:

```sh
go mod tidy
```

Para executar a aplicação:

```sh
go run main.go
````

Para fazer uma requisição, utilizando o `curl`, existem as seguintes formas:

1. Requisição simples:
    ```sh
    curl http://localhost:8080
    ```

2. Requisição com Token:
    ```sh
    curl -H "API_KEY: abc123" http://localhost:8080
    ```

## Como rodar os testes:

Utilize o comando abaixo:

```sh
go test ./tests/...
```