package main

import (
	"log"
	"math/rand"
	"medicontrol/models"
	"medicontrol/services"
	"time"
)

// O map `medicamentosParaAdicionar` agora define os códigos ANVISA e suas categorias.
var medicamentosParaAdicionar = map[string]string{
	"1097401420043": "Analgésico",        // Paracetamol (Tylenol) 500 mg – Janssen-Cilag
	"1024701490043": "Analgésico",        // Dipirona Monoidratada (Novalgina) 500 mg/ml – Sanofi
	"1004307270013": "Antibiótico",       // Amoxicilina 500 mg – EMS
	"1781700780021": "Anti-hipertensivo", // Losartana Potássica 50mg - Genérico Medley (Exemplo adicionado)
	"1058302990029": "Anti-inflamatório", // Ibuprofeno 600mg - Genérico Medley (Exemplo adicionado)
}

func getKnownAnvisaCodes() []string {
	codes := make([]string, 0, len(medicamentosParaAdicionar))
	for code := range medicamentosParaAdicionar {
		codes = append(codes, code)
	}
	return codes
}

// getUniqueCategories extrai os nomes de categorias únicas do nosso mapa.
func getUniqueCategories() []string {
	catMap := make(map[string]bool)
	for _, cat := range medicamentosParaAdicionar {
		catMap[cat] = true
	}

	cats := make([]string, 0, len(catMap))
	for cat := range catMap {
		cats = append(cats, cat)
	}
	return cats
}

func main() {
	log.Println("Iniciando script para popular medicamentos...")

	rand.Seed(time.Now().UnixNano())

	if err := models.InitDB(); err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}

	// 1. Adicionar as categorias primeiro
	log.Println("Adicionando categorias...")
	categoriasUnicas := getUniqueCategories()
	for _, nomeCat := range categoriasUnicas {
		if _, err := models.AddCategoria(nomeCat); err != nil {
			log.Printf("Erro ao adicionar categoria '%s': %v. Pode já existir, o que é esperado.", nomeCat, err)
			// Não tratamos como fatal, pois a categoria pode já existir.
			// A função AddCategoria foi ajustada para não retornar erro se já existir.
		}
	}
	log.Println("Categorias adicionadas/verificadas.")

	knownCodes := getKnownAnvisaCodes()
	medicamentosAdicionados := 0
	medicamentosExistentes := 0
	errosAoBuscar := 0

	for _, codigoAnvisa := range knownCodes {
		log.Printf("Processando código ANVISA: %s", codigoAnvisa)

		existente := models.GetMedicamentoByCodigoANVISA(codigoAnvisa)
		if existente != nil {
			log.Printf("Medicamento com código ANVISA %s já existe no banco. ID: %s, Nome: %s", codigoAnvisa, existente.ID, existente.Nome)
			medicamentosExistentes++
			continue
		}

		dadosAnvisa, err := services.BuscarDadosAnvisa(codigoAnvisa)
		if err != nil {
			log.Printf("Erro ao buscar dados para o código ANVISA %s: %v", codigoAnvisa, err)
			errosAoBuscar++
			continue
		}

		// Buscar a categoria que definimos para este medicamento
		nomeCategoria := medicamentosParaAdicionar[codigoAnvisa]
		categoria, err := models.GetCategoriaByNome(nomeCategoria)
		if err != nil || categoria == nil {
			log.Printf("ERRO CRÍTICO: Categoria '%s' não encontrada no banco para o código %s. Isso não deveria acontecer.", nomeCategoria, codigoAnvisa)
			continue
		}

		// Criar e adicionar novo medicamento
		novoMedicamento := &models.Medicamento{
			Nome:         dadosAnvisa.Nome,
			Fabricante:   dadosAnvisa.Fabricante,
			CodigoANVISA: dadosAnvisa.Registro,
			Quantidade:   rand.Intn(200) + 50,
			CategoriaID:  categoria.ID, // Atribuir o ID da categoria
		}

		if err := models.AddMedicamento(novoMedicamento); err != nil {
			log.Printf("Erro ao adicionar medicamento %s (ANVISA: %s): %v", novoMedicamento.Nome, codigoAnvisa, err)
			continue
		}
		log.Printf("Medicamento '%s' (Categoria: %s) adicionado com sucesso com ID: %s", novoMedicamento.Nome, nomeCategoria, novoMedicamento.ID)
		medicamentosAdicionados++
	}

	log.Printf("\n--- Resumo da Importação ---")
	log.Printf("Total de códigos processados: %d", len(knownCodes))
	log.Printf("Medicamentos adicionados com sucesso: %d", medicamentosAdicionados)
	log.Printf("Medicamentos que já existiam: %d", medicamentosExistentes)
	log.Printf("Erros ao buscar dados na ANVISA (mock): %d", errosAoBuscar)
	log.Println("Script de popularização concluído.")
}
