# Medic Control
## Sistema de Gestão para Farmácias

![Logo Medic Control](logo/logo.png)

---

## Documentação

A documentação completa do projeto está disponível nos seguintes arquivos:

- [Documentação da API](docs/api.md) - Detalhes sobre os endpoints e integrações
- [Guia de Instalação](docs/setup.md) - Instruções detalhadas de instalação e configuração

---

## Sumário
1. [Visão Geral](#1-visão-geral)
2. [Arquitetura do Sistema](#2-arquitetura-do-sistema)
3. [Funcionalidades Principais](#3-funcionalidades-principais)
4. [API Endpoints](#4-api-endpoints)
5. [Modelos de Dados](#5-modelos-de-dados)
6. [Configuração e Instalação](#6-configuração-e-instalação)
7. [Segurança](#7-segurança)
8. [Manutenção e Suporte](#8-manutenção-e-suporte)
9. [Roadmap Futuro](#9-roadmap-futuro)
10. [Contribuição](#10-contribuição)
11. [Licença](#11-licença)
12. [Contato](#12-contato)

---

## 1. Visão Geral

O Medic Control é um sistema de gestão para farmácias desenvolvido em Go, utilizando o framework Gin para a API REST e uma interface web moderna. O sistema visa automatizar e otimizar os processos de gestão de farmácias, desde o controle de estoque até as vendas.

### 1.1 Objetivos
- Automatizar processos de gestão farmacêutica
- Otimizar controle de estoque
- Melhorar a experiência do usuário
- Garantir conformidade com regulamentações

### 1.2 Público-Alvo
- Farmacêuticos
- Proprietários de farmácias
- Gerentes de farmácias
- Atendentes de farmácia

---

## 2. Arquitetura do Sistema

### 2.1 Tecnologias Utilizadas
| Componente | Tecnologia |
|------------|------------|
| Backend | Go (Golang) |
| Framework Web | Gin |
| Banco de Dados | SQL |
| Autenticação | JWT |
| Frontend | HTML, CSS, JavaScript |
| CORS | Suporte a requisições cross-origin |

### 2.2 Estrutura de Diretórios
```
medicontrol/
├── auth/           # Autenticação e autorização
├── config/         # Configurações do sistema
├── data/           # Dados de exemplo e importação
├── db/            # Configurações do banco de dados
├── handlers/      # Manipuladores de requisições HTTP
├── middleware/    # Middlewares do Gin
├── models/        # Modelos de dados
├── services/      # Lógica de negócios
├── sql/           # Scripts SQL
├── sqlutils/      # Utilitários SQL
├── static/        # Arquivos estáticos (frontend)
└── logo/          # Imagens e logos
```

---

## 3. Funcionalidades Principais

### 3.1 Gestão de Medicamentos
- ✅ Cadastro completo de medicamentos
- ✅ Controle de estoque
- ✅ Validação de dados da ANVISA
- ✅ Categorização
- ✅ Controle de validade

### 3.2 Sistema de Vendas
- ✅ Registro de vendas
- ✅ Múltiplos itens por venda
- ✅ Controle de estoque automático
- ✅ Histórico de vendas

### 3.3 Relatórios
- ✅ Relatório de vendas
- ✅ Relatório de estoque baixo
- ✅ Movimentações de estoque

### 3.4 Segurança
- ✅ Autenticação JWT
- ✅ Middleware de autorização
- ✅ CORS configurado
- ✅ Senhas criptografadas

---

## 4. API Endpoints

### 4.1 Autenticação
```http
POST /api/login
Content-Type: application/json

{
    "username": string,
    "password": string
}
```

### 4.2 Medicamentos
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | /api/medicamentos | Listar todos |
| GET | /api/medicamentos/:id | Obter um específico |
| POST | /api/medicamentos | Criar novo |
| PUT | /api/medicamentos/:id | Atualizar |
| DELETE | /api/medicamentos/:id | Deletar |

### 4.3 Vendas
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | /api/vendas | Registrar venda |
| GET | /api/vendas | Listar vendas |

### 4.4 Relatórios
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | /api/relatorios/vendas | Total de vendas |
| GET | /api/relatorios/baixo-estoque | Estoque baixo |

---

## 5. Modelos de Dados

### 5.1 Medicamento
```go
type Medicamento struct {
    ID            int       `json:"id"`
    Nome          string    `json:"nome"`
    CodigoANVISA  string    `json:"codigo_anvisa"`
    Quantidade    int       `json:"quantidade"`
    Preco         float64   `json:"preco"`
    Fabricante    string    `json:"fabricante"`
    Validade      time.Time `json:"validade"`
    CriadoEm      time.Time `json:"criado_em"`
    Categoria     Categoria `json:"categoria"`
}
```

### 5.2 Venda
```go
type Venda struct {
    ID     int       `json:"id"`
    Data   time.Time `json:"data"`
    UserID int       `json:"user_id"`
}

type VendaItem struct {
    ID            int     `json:"id"`
    VendaID       int     `json:"venda_id"`
    MedicamentoID int     `json:"medicamento_id"`
    Quantidade    int     `json:"quantidade"`
    PrecoUnitario float64 `json:"preco_unitario"`
}
```

---

## 6. Configuração e Instalação

### 6.1 Pré-requisitos
- Go 1.16 ou superior
- Banco de dados SQL
- Git

### 6.2 Instalação
1. Clone o repositório:
```bash
git clone [URL_DO_REPOSITORIO]
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure o banco de dados:
- Crie um banco de dados
- Execute os scripts SQL em `sql/`

4. Configure as variáveis de ambiente:
```bash
export DB_HOST=localhost
export DB_USER=seu_usuario
export DB_PASSWORD=sua_senha
export DB_NAME=medicontrol
```

5. Execute o servidor:
```bash
go run main.go
```

---

## 7. Segurança

### 7.1 Autenticação
- Tokens JWT com expiração de 24 horas
- Senhas criptografadas com bcrypt
- Middleware de autenticação para rotas protegidas

### 7.2 CORS
```go
cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
})
```

---

## 8. Manutenção e Suporte

### 8.1 Logs
- Logs de sistema em `logs/`
- Logs de erro e acesso
- Monitoramento de performance

### 8.2 Backup
- Backup automático do banco de dados
- Backup de configurações
- Procedimentos de recuperação

---

## 9. Roadmap Futuro

### Sprint 1
- Sistema de alerta de vencimento
- Gestão de clientes
- Dashboard básico

### Sprint 2
- Relatórios financeiros
- Controle de lotes
- Sistema de backup

### Sprint 3
- API para e-commerce
- Sistema de notificações
- Programa de fidelidade

---

## 10. Contribuição
Para contribuir com o projeto:
1. Fork o repositório
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Crie um Pull Request

---

## 11. Licença
[Inserir informações sobre a licença do projeto]

---

## 12. Contato
[Inserir informações de contato para suporte e desenvolvimento]

---

*Documentação gerada em: [DATA]*
*Versão: 1.0.0* 