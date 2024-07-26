# Rate Limiter

## Como rodar o projeto:

Use o Docker Compose para subir o Redis e a aplicação, rodando o comando abaixo na pasta raiz:

```sh
docker-compose up --build
```

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