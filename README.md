# Rate Limiter API
Uma API em Go que implementa um sistema de rate limiting usando Redis para controle de requisições por IP e Token.

---

## 🚀 Requisitos
- Docker
- Docker Compose
- Go 1.23.8 (para desenvolvimento local)

---

## ⚙️ Configuração
1. Clone o repositório:
```bash
git clone https://github.com/AndreD23/goexpert-desafio-ratelimiter.git
cd goexpert-desafio-ratelimiter
```

2. Crie o arquivo de configuração `.env`:
```bash
cp .env.example .env
```

---

## 🐳 Executando com Docker Compose

1. Inicie os containers:
```bash
docker compose up -d --build
```

2. Verifique se os containers estão rodando:
```bash
docker compose ps
```

---

## 🔨 Desenvolvimento Local

1. Instale as dependências:
```bash
go mod download
```

2. Execute a aplicação:
```bash
go run main.go
```

---

## 🧪 Executando Testes

Execute os testes automatizados:
```bash
go test ./... -v
```

---

## 📡 Utilizando a API

### Endpoint Base
```
http://localhost:8080
```

### Exemplos de Requisições

1. Requisição sem token (limitado por IP):
```bash
curl http://localhost:8080/
```

2. Requisição com token (limitado por token):
```bash
curl -H "API_KEY: seu-token" http://localhost:8080/
```

## 🚫 Possíveis Respostas de Erro

### Limite Excedido (HTTP 429)
```json
{
    "error": "You have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

### Erro Interno (HTTP 500)
```json
{
    "error": "Error checking the limit"
}
```

### Erro de IP (HTTP 500)
```json
{
    "error": "Unable to parse IP address"
}
```

---

## ⚠️ Limites de Requisições

- **Por IP**: 10 requisições por segundo (configurável via REQUESTS_PER_IP)
- **Por Token**: 20 requisições por segundo (configurável via REQUESTS_PER_TOKEN)
- **Tempo de Bloqueio**: 30 segundos (configurável via BLOCK_DURATION)

---

## 🔍 Monitoramento Redis

Para verificar as chaves no Redis:
```bash
docker exec -it redis redis-cli
```

Comandos úteis:
```redis
KEYS *           # Lista todas as chaves
GET chave        # Obtém valor de uma chave
TTL chave        # Verifica tempo restante de uma chave
```

---

## 🛠️ Configurações Avançadas

O arquivo `docker-compose.yaml` pode ser modificado para ajustar:
- Portas expostas
- Configurações de rede
- Variáveis de ambiente
- Persistência de dados Redis

---

## 📝 Notas

- O sistema utiliza Redis para armazenamento dos contadores
- As configurações podem ser ajustadas através de variáveis de ambiente
- Em caso de falha do Redis, o sistema retornará erro 500
