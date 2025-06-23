package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// MedicamentoImportado representa um medicamento do arquivo de importação
type MedicamentoImportado struct {
	ID                int     `json:"id"`
	Nome              string  `json:"nome"`
	CodigoAnvisa      string  `json:"codigo_anvisa"`
	QuantidadeEstoque int     `json:"quantidade_estoque"`
	PrecoVenda        float64 `json:"preco_venda"`
	Fabricante        string  `json:"fabricante"`
	DataValidade      string  `json:"data_validade"`
	DataEntrada       string  `json:"data_entrada"`
	Bula              string  `json:"bula"`
}

// ImportarMedicamentos carrega medicamentos do arquivo JSON e os adiciona ao banco de dados
func ImportarMedicamentos() error {
	// Verificar se já existem medicamentos no banco para evitar duplicação
	medicamentosExistentes, err := GetMedicamentos()
	if err == nil && len(medicamentosExistentes) > 0 {
		log.Println("O banco de dados já contém medicamentos. A importação será ignorada.")
		return nil // Não retorna erro, simplesmente não faz nada
	}

	// Caminho para o arquivo de dados
	dataDir := "data"
	filename := "medicamentos_500_com_bula.json"
	filePath := filepath.Join(dataDir, filename)

	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de medicamentos: %v", err)
	}
	defer file.Close()

	// Ler o conteúdo do arquivo
	var medicamentos []MedicamentoImportado
	if err := json.NewDecoder(file).Decode(&medicamentos); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	// Importar cada medicamento
	for _, medData := range medicamentos {
		// Criar novo medicamento
		med := &Medicamento{
			ID:           uuid.New().String(),
			Nome:         medData.Nome,
			Fabricante:   medData.Fabricante,
			CodigoANVISA: medData.CodigoAnvisa,
			Quantidade:   medData.QuantidadeEstoque,
			Preco:        medData.PrecoVenda,
			Validade:     medData.DataValidade,
			CriadoEm:     time.Now(),
		}

		// Adicionar ao banco de dados
		if err := AddMedicamento(med); err != nil {
			log.Printf("Erro ao adicionar medicamento %s: %v", med.Nome, err)
			continue
		}
		log.Printf("Medicamento %s adicionado com sucesso", med.Nome)
	}

	return nil
}

// CorrigirPrecosMedicamentos atualiza os preços dos medicamentos existentes no banco com base no arquivo JSON.
func CorrigirPrecosMedicamentos() error {
	log.Println("Iniciando a correção de preços dos medicamentos existentes...")

	// 1. Carregar os dados do JSON
	dataDir := "data"
	filename := "medicamentos_500_com_bula.json"
	filePath := filepath.Join(dataDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de medicamentos para correção: %v", err)
	}
	defer file.Close()

	var medicamentosImportados []MedicamentoImportado
	if err := json.NewDecoder(file).Decode(&medicamentosImportados); err != nil {
		return fmt.Errorf("erro ao decodificar JSON para correção: %v", err)
	}

	// Mapear os preços por CodigoAnvisa para fácil acesso
	precosPorCodigo := make(map[string]float64)
	for _, med := range medicamentosImportados {
		precosPorCodigo[med.CodigoAnvisa] = med.PrecoVenda
	}

	// 2. Buscar todos os medicamentos do banco de dados
	medicamentosDB, err := GetMedicamentos()
	if err != nil {
		return fmt.Errorf("erro ao buscar medicamentos do banco para correção: %v", err)
	}

	// 3. Atualizar cada medicamento no banco de dados
	queryUpdatePreco := "UPDATE medicamentos SET preco = ? WHERE CodigoANVISA = ?"
	for _, medDB := range medicamentosDB {
		if preco, ok := precosPorCodigo[medDB.CodigoANVISA]; ok {
			// Apenas atualiza se o preço atual for 0 ou nulo, para segurança
			if medDB.Preco == 0 {
				_, err := sqlDB.Exec(queryUpdatePreco, preco, medDB.CodigoANVISA)
				if err != nil {
					log.Printf("Erro ao atualizar o preço para o medicamento com código ANVISA %s: %v", medDB.CodigoANVISA, err)
					// Continua para o próximo mesmo se um falhar
				} else {
					log.Printf("Preço do medicamento %s (ANVISA: %s) corrigido para R$ %.2f", medDB.Nome, medDB.CodigoANVISA, preco)
				}
			}
		}
	}

	log.Println("Correção de preços concluída.")
	return nil
}
