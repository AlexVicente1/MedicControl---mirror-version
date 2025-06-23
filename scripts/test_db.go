package main

import (
	"fmt"
	"log"
	"medicontrol/models"
)

func main() {
	log.Println("Iniciando teste do banco de dados...")

	// Inicializar o banco de dados
	if err := models.InitDB(); err != nil {
		log.Fatal("Erro ao inicializar banco de dados:", err)
	}

	// Teste 1: Adicionar um medicamento
	med1 := &models.Medicamento{
		Nome:         "Dipirona",
		Fabricante:   "Medley",
		CodigoANVISA: "1234567890",
		Quantidade:   100,
	}

	log.Println("Adicionando medicamento de teste...")
	if err := models.AddMedicamento(med1); err != nil {
		log.Fatal("Erro ao adicionar medicamento:", err)
	}

	// Teste 2: Listar medicamentos
	log.Println("\nListando todos os medicamentos:")
	meds := models.GetMedicamentos()
	for _, med := range meds {
		fmt.Printf("ID: %s\nNome: %s\nFabricante: %s\nQuantidade: %d\n\n",
			med.ID, med.Nome, med.Fabricante, med.Quantidade)
	}

	// Teste 3: Registrar uma movimentação de entrada
	mov1 := models.Movimentacao{
		MedicamentoID: med1.ID,
		Tipo:          "entrada",
		Quantidade:    50,
		Observacao:    "Entrada inicial de estoque",
	}

	log.Println("Registrando movimentação de entrada...")
	if err := models.RegistrarMovimentacao(mov1); err != nil {
		log.Fatal("Erro ao registrar movimentação:", err)
	}

	// Teste 4: Registrar uma movimentação de saída
	mov2 := models.Movimentacao{
		MedicamentoID: med1.ID,
		Tipo:          "saida",
		Quantidade:    30,
		Observacao:    "Venda teste",
	}

	log.Println("Registrando movimentação de saída...")
	if err := models.RegistrarMovimentacao(mov2); err != nil {
		log.Fatal("Erro ao registrar movimentação:", err)
	}

	// Teste 5: Verificar movimentações
	log.Println("\nListando todas as movimentações:")
	movs := models.GetMovimentacoes()
	for _, mov := range movs {
		fmt.Printf("ID: %s\nTipo: %s\nQuantidade: %d\nObservação: %s\n\n",
			mov.ID, mov.Tipo, mov.Quantidade, mov.Observacao)
	}

	// Teste 6: Verificar total de vendas
	log.Println("\nVerificando estatísticas de vendas:")
	vendas := models.GetTotalVendas()
	fmt.Printf("Total de vendas: %v\n", vendas)

	// Teste 7: Atualizar medicamento
	med1.Nome = "Dipirona Sódica"
	log.Println("\nAtualizando nome do medicamento...")
	if err := models.UpdateMedicamento(med1); err != nil {
		log.Fatal("Erro ao atualizar medicamento:", err)
	}

	// Verificar atualização
	medAtualizado := models.GetMedicamento(med1.ID)
	fmt.Printf("\nMedicamento atualizado:\nNome: %s\nQuantidade atual: %d\n",
		medAtualizado.Nome, medAtualizado.Quantidade)

	log.Println("\nTestes concluídos com sucesso!")
}
