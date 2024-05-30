# Rate Limiter em Go

Este projeto implementa um sistema de rate limiting em Go, permitindo controlar o tráfego de requisições para um serviço web com base em endereço IP ou token de acesso. O rate limiter pode ser configurado para limitar o número de requisições por segundo e aplicar bloqueios temporários quando os limites são excedidos.

## Funcionalidades

- Limitação de requisições por endereço IP.
- Limitação de requisições por token de acesso (API_KEY).
- Prioridade para limitação por token de acesso sobre a limitação por IP.
- Configuração do número máximo de requisições permitidas por segundo.
- Opção de bloqueio temporário de IP ou token após exceder o limite de requisições.
- Armazenamento de informações de limitação em um banco de dados Redis.
- Middleware para fácil integração com um servidor web usando o framework Gin.

## Como Funciona

O rate limiter utiliza o Redis para rastrear o número de requisições feitas por um IP ou token de acesso em um determinado intervalo de tempo. Se o número de requisições exceder o limite configurado, o IP ou token será bloqueado por um período de tempo especificado.

### Estrutura do Projeto
```plaintext
rate-limiter/
├── cmd/
│   └── main.go           # Ponto de entrada da aplicação
├── pkg/
│   ├── limiter/
│   │   ├── limiter.go    # Implementação do rate limiter
│   │   └── limiter_test.go # Testes para o rate limiter
│   └── middleware/
│       ├── middleware.go # Middleware para integração com o Gin
│       └── middleware_test.go # Testes para o middleware
├── .env                  # Variáveis de ambiente
├── .gitignore
├── docker-compose.yml    # Configuração do Docker Compose
├── Dockerfile            # Dockerfile para a aplicação Go
└── go.mod                # Gerenciamento de dependências Go
```

## Configuração

### Variáveis de Ambiente

As configurações do rate limiter são feitas através de variáveis de ambiente definidas em um arquivo .env na raiz do projeto:

```env
REDIS_ADDR=redis:6379             # Endereço do Redis
LIMIT_PER_SECOND_IP=5             # Limite de requisições por segundo por IP
LIMIT_PER_SECOND_TOKEN=10         # Limite de requisições por segundo por token de acesso
BLOCK_TIME_SECONDS=300            # Tempo de bloqueio em segundos após exceder o limite de requisições
````

### Docker Compose

O projeto inclui um arquivo docker-compose.yml para facilitar o setup do ambiente de desenvolvimento:
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      - redis

  redis:
    image: "redis:alpine"
    ports:
      - "6379:6379"
```

## Como Rodar

### Pré-requisitos
- Docker
- Docker Compose

### Passos
1.	Clone o repositório:
```bash
git clone https://github.com/seu-usuario/rate-limiter-go.git
cd rate-limiter-go
```
2. Crie um arquivo .env na raiz do projeto e configure as variáveis de ambiente conforme descrito acima.
3. Execute o Docker Compose para iniciar a aplicação e o Redis:
```bash
docker-compose up --build
```
4. Acesse a aplicação em http://localhost:8080.

## Testes
Os testes unitários para o rate limiter e o middleware podem ser executados com o comando:
```bash
go test ./...
```

## Exemplos de Uso
### Limitação por IP

Faça 10 requisições com o mesmo token dentro de um segundo. A décima primeira requisição deve ser bloqueada:
```bash
for i in {1..5}; do curl -i http://localhost:8080/; done
```

Caso o comando anterior não funcione corraretaente, tente o seguinte:
```bash
for i in $(seq 1 5); do
  curl -i http://localhost:8080/
done
```

### Limitação por Token
Faça 20 requisições com o mesmo token dentro de um segundo. A vigésima primeira requisição deve ser bloqueada:
```bash
for i in {1..10}; do curl -i -H "API_KEY: abc123" http://localhost:8080/; done
curl -i -H "API_KEY: abc123" http://localhost:8080/
```

## Sobre
### Este projeto foi desenvolvido para a pós-graduação Full Cycle Go Expert.