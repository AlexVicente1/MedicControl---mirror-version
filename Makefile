# Makefile para o projeto Medicontrol

# Nome do executável de saída
BINARY_NAME=medicontrol

# ==============================================================================
# Comandos de Build
# ==============================================================================

## build: Compila o projeto para desenvolvimento.
build:
	@echo "Compilando o projeto (desenvolvimento)..."
	go build -o ./${BINARY_NAME} .

## build-secure: Compila o projeto com ofuscação para produção.
build-secure:
	@echo "Compilando e ofuscando o projeto (produção)..."
	garble build -o ./${BINARY_NAME} .

## clean: Remove os binários compilados.
clean:
	@echo "Limpando arquivos de build..."
	rm -f ./${BINARY_NAME} ./${BINARY_NAME}.exe

## test: Roda os testes do projeto.
test:
	@echo "Rodando testes..."
	go test ./...

.PHONY: build build-secure clean test 
