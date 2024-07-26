## Deploy Cloud Run

Endereço para acessar o serviço deployado:
```sh
https://cloudrun-goexpert-c32dlgo2vq-uc.a.run.app/temperatura/{cep}
```

## Como rodar o projeto:

Use o Docker Compose para subir o Redis e a aplicação, rodando o comando abaixo na pasta raiz:

```sh
docker-compose up --build
```

Para fazer uma requisição, utilizando o `curl`:

```sh
curl http://localhost:8080/temperatura/{cep}
```

## Como rodar os testes:

Utilize o comando abaixo:

```sh
go test ./tests/...
```