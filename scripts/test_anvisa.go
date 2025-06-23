package main

import (
	"fmt"
	"log"
	"medicontrol/services"
	"time"
)

func main() {
	log.Println("Testando serviço da ANVISA...")

	// Teste 1: Buscar medicamento existente
	codigo := "1234567890123"
	log.Printf("Buscando medicamento com código %s...\n", codigo)
	dados, err := services.BuscarDadosAnvisa(codigo)
	if err != nil {
		log.Fatal("Erro ao buscar dados:", err)
	}
	fmt.Printf("\nDados encontrados:\nNome: %s\nFabricante: %s\nClasse: %s\n\n",
		dados.Nome, dados.Fabricante, dados.Classe)

	// Teste 2: Buscar medicamento inexistente
	codigoInexistente := "0000000000000"
	log.Printf("Tentando buscar medicamento inexistente %s...\n", codigoInexistente)
	_, err = services.BuscarDadosAnvisa(codigoInexistente)
	if err != nil {
		fmt.Printf("Erro esperado: %v\n\n", err)
	}

	// Teste 3: Testar cache
	log.Println("Testando cache - buscando mesmo medicamento novamente...")
	start := time.Now()
	dados, err = services.BuscarDadosAnvisa(codigo)
	if err != nil {
		log.Fatal("Erro ao buscar dados do cache:", err)
	}
	fmt.Printf("Tempo de resposta do cache: %v\n", time.Since(start))

	log.Println("\nTestes concluídos com sucesso!")
}
