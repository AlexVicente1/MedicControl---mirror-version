#!/bin/bash

# Criar diretórios necessários
mkdir -p static/css static/js data

# Criar arquivos estáticos
echo "Criando arquivos estáticos..."

# Atualizar dependências
go mod tidy

echo "Configuração concluída!"
echo "Para iniciar o sistema, execute: go run main.go"
echo "Depois acesse: http://localhost:8080"
echo "Credenciais: admin / senha123" 