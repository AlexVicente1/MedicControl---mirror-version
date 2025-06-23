# Documentação da API - Medic Control

## Endpoints da API

### Autenticação

#### Login
```http
POST /api/login
Content-Type: application/json

{
    "username": string,
    "password": string
}
```

**Resposta de Sucesso:**
```json
{
    "token": "jwt_token_here"
}
```

### Medicamentos

#### Listar Medicamentos
```http
GET /api/medicamentos
Authorization: Bearer {token}
```

#### Obter Medicamento Específico
```http
GET /api/medicamentos/:id
Authorization: Bearer {token}
```

#### Criar Medicamento
```http
POST /api/medicamentos
Authorization: Bearer {token}
Content-Type: application/json

{
    "nome": string,
    "codigo_anvisa": string,
    "quantidade": number,
    "preco": number,
    "fabricante": string,
    "validade": string (ISO date),
    "categoria_id": number
}
```

#### Atualizar Medicamento
```http
PUT /api/medicamentos/:id
Authorization: Bearer {token}
Content-Type: application/json

{
    "nome": string,
    "codigo_anvisa": string,
    "quantidade": number,
    "preco": number,
    "fabricante": string,
    "validade": string (ISO date),
    "categoria_id": number
}
```

#### Deletar Medicamento
```http
DELETE /api/medicamentos/:id
Authorization: Bearer {token}
```

### Vendas

#### Registrar Venda
```http
POST /api/vendas
Authorization: Bearer {token}
Content-Type: application/json

{
    "itens": [
        {
            "medicamento_id": number,
            "quantidade": number
        }
    ]
}
```

#### Listar Vendas
```http
GET /api/vendas
Authorization: Bearer {token}
```

### Relatórios

#### Total de Vendas
```http
GET /api/relatorios/vendas
Authorization: Bearer {token}
```

#### Estoque Baixo
```http
GET /api/relatorios/baixo-estoque?limite=50
Authorization: Bearer {token}
```

### ANVISA

#### Consultar Dados ANVISA
```http
GET /api/anvisa/:codigo
Authorization: Bearer {token}
```

## Códigos de Status

- 200: Sucesso
- 201: Criado
- 400: Requisição inválida
- 401: Não autorizado
- 403: Proibido
- 404: Não encontrado
- 500: Erro interno do servidor

## Exemplos de Uso

### Exemplo de Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "senha123"}'
```

### Exemplo de Listagem de Medicamentos
```bash
curl -X GET http://localhost:8080/api/medicamentos \
  -H "Authorization: Bearer {seu_token}"
```

### Exemplo de Criação de Medicamento
```bash
curl -X POST http://localhost:8080/api/medicamentos \
  -H "Authorization: Bearer {seu_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Paracetamol 500mg",
    "codigo_anvisa": "1234567890123",
    "quantidade": 100,
    "preco": 5.99,
    "fabricante": "EMS",
    "validade": "2025-12-31T00:00:00Z",
    "categoria_id": 1
  }'
``` 