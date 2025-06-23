# Guia de Instalação - Medic Control

## Pré-requisitos

### Sistema Operacional
- Windows 10 ou superior
- Linux (Ubuntu 20.04 ou superior)
- macOS 10.15 ou superior

### Software Necessário
- Go 1.16 ou superior
- Git
- Banco de dados SQL (MySQL/PostgreSQL)
- Editor de código (VS Code recomendado)

## Instalação

### 1. Clonar o Repositório
```bash
git clone https://github.com/seu-usuario/medicontrol.git
cd medicontrol
```

### 2. Instalar Dependências
```bash
go mod download
```

### 3. Configurar Banco de Dados

#### MySQL
```sql
CREATE DATABASE medicontrol;
CREATE USER 'medicontrol_user'@'localhost' IDENTIFIED BY 'sua_senha';
GRANT ALL PRIVILEGES ON medicontrol.* TO 'medicontrol_user'@'localhost';
FLUSH PRIVILEGES;
```

#### PostgreSQL
```sql
CREATE DATABASE medicontrol;
CREATE USER medicontrol_user WITH PASSWORD 'sua_senha';
GRANT ALL PRIVILEGES ON DATABASE medicontrol TO medicontrol_user;
```

### 4. Configurar Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=medicontrol_user
DB_PASSWORD=sua_senha
DB_NAME=medicontrol
JWT_SECRET=seu_segredo_super_secreto
```

### 5. Executar Scripts SQL

```bash
# Para MySQL
mysql -u medicontrol_user -p medicontrol < sql/schema.sql

# Para PostgreSQL
psql -U medicontrol_user -d medicontrol -f sql/schema.sql
```

### 6. Compilar e Executar

```bash
# Compilar
go build -o medicontrol

# Executar
./medicontrol
```

## Verificação da Instalação

1. Acesse `http://localhost:8080` no navegador
2. Faça login com as credenciais padrão:
   - Usuário: admin
   - Senha: senha123

## Solução de Problemas

### Erro de Conexão com Banco de Dados
- Verifique se o banco de dados está rodando
- Confirme as credenciais no arquivo `.env`
- Verifique se o banco de dados foi criado corretamente

### Erro de Compilação
- Verifique se o Go está instalado corretamente
- Execute `go mod tidy` para atualizar dependências
- Verifique se todas as dependências foram baixadas

### Erro de Execução
- Verifique se a porta 8080 está disponível
- Confirme se todas as variáveis de ambiente estão configuradas
- Verifique os logs do sistema

## Atualização

Para atualizar o sistema:

```bash
git pull origin main
go mod download
go build -o medicontrol
```

## Backup

### Backup do Banco de Dados

#### MySQL
```bash
mysqldump -u medicontrol_user -p medicontrol > backup.sql
```

#### PostgreSQL
```bash
pg_dump -U medicontrol_user medicontrol > backup.sql
```

### Restauração

#### MySQL
```bash
mysql -u medicontrol_user -p medicontrol < backup.sql
```

#### PostgreSQL
```bash
psql -U medicontrol_user -d medicontrol -f backup.sql
```

## Segurança

### Recomendações
1. Altere a senha padrão do administrador
2. Configure um JWT_SECRET forte
3. Use HTTPS em produção
4. Mantenha o sistema atualizado
5. Faça backups regulares

### Firewall
- Abra apenas a porta 8080 (ou a porta configurada)
- Configure regras de firewall adequadas
- Use um proxy reverso em produção 
