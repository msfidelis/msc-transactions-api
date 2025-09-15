# MSC Transactions API

> Aplicação Transacional utilizada para os laboratórios da minha pesquisa de mestrado. 

Uma API de transações financeiras construída em Go usando Fiber framework, PostgreSQL e Redis. Esta API permite gerenciar transações de débito e crédito para clientes, mantendo controle de saldo e limites.

## 🏗️ Arquitetura

- **Framework**: Fiber (Go)
- **Banco de Dados**: PostgreSQL 16
- **Cache**: Redis 7.2.4
- **ORM**: Bun
- **Containerização**: Docker Compose
- **Monitoramento**: Prometheus metrics


## 🐳 Ambiente Local

```bash
# Clone o repositório
git clone https://github.com/msfidelis/msc-transactions-api
cd msc-transactions-api

# Execute com Docker Compose
docker-compose up -d

# A API estará disponível em http://localhost:8081
```

## 🔧 Configuração

### Variáveis de Ambiente

| Variável | Descrição | Valor Padrão |
|----------|-----------|--------------|
| `DATABASE_HOST` | Host do PostgreSQL | `postgres` |
| `DATABASE_PORT` | Porta do PostgreSQL | `5432` |
| `DATABASE_DB` | Nome do banco | `transactions` |
| `DATABASE_USER` | Usuário do banco | `fidelissauro` |
| `DATABASE_PASSWORD` | Senha do banco | `doutorequemtemdoutorado` |
| `CACHE_HOST` | Host do Redis | `redis` |
| `CACHE_PORT` | Porta do Redis | `6379` |
| `ENV` | Ambiente (shadow para rollback) | - |

## 📚 API Endpoints

### 🏥 Health Check

```http
GET /healthcheck
```

**Resposta:**
```json
{
  "status": "ok"
}
```

### 💰 Criar Transação

```http
POST /transactions
Content-Type: application/json
id_client: {uuid_do_cliente}
```

**Payload:**
```json
{
  "amount": 100,
  "type": "c",
  "description": "Teste"
}
```

**Tipos de Transação:**
- `c`: Crédito (soma ao saldo)
- `d`: Débito (subtrai do saldo)

**Resposta de Sucesso (200):**
```json
{
  "limits": 100000,
  "balance": 150
}
```

**Possíveis Erros:**
- `400`: Payload inválido
- `422`: Limite insuficiente
- `500`: Erro interno do servidor

### 📄 Consultar Extrato

```http
GET /statements
id_client: {uuid_do_cliente}
```

**Resposta:**
```json
{
  "balance": {
    "total": 150,
    "date_statement": "2025-09-15T23:18:53.817444464Z",
    "limits": 100000
  },
  "last_transactions": [
    {
      "id": 1,
      "id_client": "da2ae01e-100a-49a6-af7e-6ec18ee53c1b",
      "amount": 100,
      "type": "c",
      "description": "Depósito",
      "date": "2025-09-15T23:18:53.817444464Z"
    }
  ]
}
```

### 🔍 Detalhe da Transação

```http
GET /statements/{id_transaction}
id_client: {uuid_do_cliente}
```

**Resposta:**
```json
{
  "id": 1,
  "id_client": "da2ae01e-100a-49a6-af7e-6ec18ee53c1b",
  "amount": 100,
  "type": "c",
  "description": "Depósito",
  "date": "2025-09-15T23:18:53.817444464Z"
}
```

## 📊 Métricas

As métricas do Prometheus estão disponíveis em:

```http
GET /metrics
```

## 🗄️ Estrutura do Banco

### Tabela: `clients`
```sql
CREATE TABLE clients (
    id_client UUID PRIMARY KEY,
    balance INTEGER NOT NULL DEFAULT 0,
    limits INTEGER NOT NULL
);
```

### Tabela: `transactions`
```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    id_client UUID REFERENCES clients(id_client),
    amount INTEGER NOT NULL,
    type CHAR(1) NOT NULL CHECK (type IN ('c', 'd')),
    description TEXT NOT NULL,
    date TEXT NOT NULL
);
```

## 🧪 Exemplos de Uso

### Criar um Depósito

```bash
curl -X POST http://localhost:8081/transactions \
  -H "Content-Type: application/json" \
  -H "id_client: da2ae01e-100a-49a6-af7e-6ec18ee53c1b" \
  -d '{
    "amount": 500,
    "type": "c",
    "description": "Depósito inicial"
  }'
```

### Fazer um Saque

```bash
curl -X POST http://localhost:8081/transactions \
  -H "Content-Type: application/json" \
  -H "id_client: da2ae01e-100a-49a6-af7e-6ec18ee53c1b" \
  -d '{
    "amount": 200,
    "type": "d",
    "description": "Saque no caixa"
  }'
```

### Consultar Extrato

```bash
curl -X GET http://localhost:8081/statements \
  -H "id_client: da2ae01e-100a-49a6-af7e-6ec18ee53c1b"
```

### Consultar Transação Específica

```bash
curl -X GET http://localhost:8081/statements/1 \
  -H "id_client: da2ae01e-100a-49a6-af7e-6ec18ee53c1b"
```
