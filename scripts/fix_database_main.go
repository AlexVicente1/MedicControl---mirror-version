package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Medicamento struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Fabricante   string `json:"fabricante"`
	Tipo         string `json:"tipo"`
	CodigoANVISA string `json:"codigo_anvisa"`
	Quantidade   int    `json:"quantidade"`
	Validade     string `json:"validade"`
	CriadoEm     string `json:"criado_em"`
}

type Movimentacao struct {
	ID            string `json:"id"`
	MedicamentoID string `json:"medicamento_id"`
	Tipo          string `json:"tipo"`
	Quantidade    int    `json:"quantidade"`
	Data          string `json:"data"`
	Observacao    string `json:"observacao"`
}

type Database struct {
	Medicamentos  []Medicamento  `json:"medicamentos"`
	Movimentacoes []Movimentacao `json:"movimentacoes,omitempty"`
}

func main() {
	// Ler o arquivo database.json
	data, err := ioutil.ReadFile("./data/database.json")
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo database.json: %v", err)
	}

	// Verificar se o arquivo está vazio
	if len(data) == 0 {
		log.Fatal("O arquivo database.json está vazio")
	}

	// Tentar decodificar para a estrutura correta
	var db Database
	err = json.Unmarshal(data, &db)
	if err != nil {
		log.Printf("Erro ao decodificar o arquivo JSON: %v. Tentando extrair movimentações...", err)
		
		// Se falhar, tentar extrair as movimentações manualmente
		var tempMap map[string]json.RawMessage
		err = json.Unmarshal(data, &tempMap)
		if err != nil {
			log.Fatalf("Erro ao analisar o JSON: %v", err)
		}

		// Extrair medicamentos
		if medicamentosData, ok := tempMap["medicamentos"]; ok {
			err = json.Unmarshal(medicamentosData, &db.Medicamentos)
			if err != nil {
				log.Fatalf("Erro ao decodificar medicamentos: %v", err)
			}
		}

		// Extrair movimentações do final do array de medicamentos
		for i := len(db.Medicamentos) - 1; i >= 0; i-- {
			med := db.Medicamentos[i]
			// Verificar se é uma movimentação (tem uma estrutura diferente de medicamento)
			// Vamos verificar se tem um campo que só existe em movimentações
			// Como não temos acesso direto, vamos verificar se tem um formato que se parece com uma movimentação
			if med.Nome == "" && med.Fabricante == "" && med.Tipo != "" && med.Quantidade > 0 {
				// É uma movimentação
				// Criar uma nova movimentação com os dados do objeto
				// Vamos converter o objeto para JSON e depois para a estrutura de Movimentacao
				medJSON, err := json.Marshal(med)
				if err != nil {
					log.Printf("Erro ao converter medicamento para JSON: %v", err)
					continue
				}

				var mov Movimentacao
				err = json.Unmarshal(medJSON, &mov)
				if err != nil {
					log.Printf("Erro ao converter JSON para movimentação: %v", err)
					continue
				}

				// Se o ID da movimentação estiver vazio, gerar um novo
				if mov.ID == "" {
					mov.ID = med.ID
				}
				db.Movimentacoes = append(db.Movimentacoes, mov)
				// Remover a movimentação do array de medicamentos
				db.Medicamentos = append(db.Medicamentos[:i], db.Medicamentos[i+1:]...)
			}
		}
	}

	// Verificar se encontrou alguma movimentação
	if len(db.Movimentacoes) == 0 {
		log.Println("Nenhuma movimentação encontrada para extrair")
	} else {
		log.Printf("Encontradas %d movimentações", len(db.Movimentacoes))
	}

	// Criar backup do arquivo original
	err = os.Rename("./data/database.json", "./data/database.json.bak")
	if err != nil {
		log.Fatalf("Erro ao criar backup do arquivo: %v", err)
	}
	log.Println("Backup criado: ./data/database.json.bak")

	// Salvar o arquivo corrigido
	saveDatabase(db)
}

func saveDatabase(db Database) {
	// Ordenar as movimentações por data
	for i := 0; i < len(db.Movimentacoes); i++ {
		for j := i + 1; j < len(db.Movimentacoes); j++ {
			if db.Movimentacoes[i].Data > db.Movimentacoes[j].Data {
				db.Movimentacoes[i], db.Movimentacoes[j] = db.Movimentacoes[j], db.Movimentacoes[i]
			}
		}
	}

	// Converter para JSON formatado
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		log.Fatalf("Erro ao converter para JSON: %v", err)
	}

	// Escrever no arquivo
	err = ioutil.WriteFile("./data/database.json", data, 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar o arquivo: %v", err)
	}

	log.Println("Arquivo database.json corrigido com sucesso!")
	log.Printf("Total de medicamentos: %d", len(db.Medicamentos))
	log.Printf("Total de movimentações: %d", len(db.Movimentacoes))
}
