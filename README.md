## Token Service (sem banco de dados)

### Variaveis de ambiente

- `DEFENSE_HOST`
- `DEFENSE_PORT`
- `DEFENSE_SCHEME` (opcional, padrao `https`)
- `DEFENSE_INSECURE_SKIP_VERIFY` (opcional, `true`/`false`)
- `DEFENSE_USERNAME`
- `DEFENSE_PASSWORD`
- `HTTP_ADDR` (opcional, padrao `:8080`)
- `HTTP_API_KEY`
- `ENVIRONMENT` (opcional, padrao `DEVELOPMENT`)

### Como rodar

```
go run ./cmd/token-service
```

### Como buscar o token

- Requisicao HTTP: `GET /token` (ex: `http://localhost:8080/token`) com `X-API-Key`
- Console: o token e impresso no log quando e criado

Exemplo:
```
curl -H "X-API-Key: change-me" http://localhost:8080/token
```
