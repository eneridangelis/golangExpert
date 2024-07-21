# Rate Limiter

## Como rodar o projeto:

Build a imagem `docker` do projeto:

```sh
docker run -it --entrypoint /bin/sh loadtest
```

Rode o projeto dentro do container passando as flags:

```sh
docker run loadtest --url=https://httpbin.org/get --requests=1000 --concurrency=10
```

O relatório será gerado no próprio terminal.
A url utilizada para teste e setada como default no código é a `https://httpbin.org/get`, pois o google impede acessos automatizados e acaba redirecionando as requisições para página de erro.