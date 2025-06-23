package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// DadosAnvisaAPI representa a estrutura da resposta da API externa.
// Os nomes dos campos correspondem ao JSON retornado pela API API_FLASK_MEDICAMENTOS.
type DadosAnvisaAPI struct {
	Produto           string `json:"PRODUTO"`
	Laboratorio       string `json:"LABORATÓRIO"`
	Registro          string `json:"REGISTRO"` // Este é o código de registro
	Apresentacao      string `json:"APRESENTAÇÃO"`
	ClasseTerapeutica string `json:"CLASSE TERAPÊUTICA"`
	TipoProduto       string `json:"TIPO DE PRODUTO (STATUS DO PRODUTO)"`
	Substancia        string `json:"SUBSTÂNCIA"`
	SituacaoRegistro  string `json:"SITUACAO_REGISTRO"` // Assumindo que pode haver um campo assim, ou precisaremos inferir
	// Adicione outros campos conforme necessário
}

// DadosAnvisa representa os dados que nosso sistema utiliza internamente.
// Mantemos esta struct para compatibilidade com o resto do sistema.
type DadosAnvisa struct {
	Nome       string
	Fabricante string
	Registro   string // Código de registro ANVISA
	Classe     string // Classe Terapêutica
	// Status (ATIVO/INATIVO) - a API externa pode não fornecer isso diretamente
}

const apiBaseURL = "https://apiflaskmedicamentos.herokuapp.com/medicamentos"

// BuscarDadosAnvisa busca informações de um medicamento pelo código de registro usando a API pública.
func BuscarDadosAnvisa(codigoRegistro string) (*DadosAnvisa, error) {
	if codigoRegistro == "" {
		return nil, errors.New("código de registro não pode ser vazio")
	}

	// A API permite buscar por ?registro=CODIGO
	url := fmt.Sprintf("%s?registro=%s", apiBaseURL, codigoRegistro)
	log.Printf("Consultando API externa: %s", url)

	httpClient := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição HTTP: %v", err)
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}
	// Algumas APIs podem requerer um User-Agent
	req.Header.Set("User-Agent", "MediControlApp/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Erro ao fazer requisição para API ANVISA: %v", err)
		return nil, fmt.Errorf("erro na comunicação com a API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("API ANVISA retornou status não OK: %d. Corpo: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("API ANVISA retornou status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler corpo da resposta da API ANVISA: %v", err)
		return nil, fmt.Errorf("erro ao ler resposta da API: %w", err)
	}

	// A API retorna uma lista, mesmo que consultando por um registro único.
	var resultadosAPI []DadosAnvisaAPI
	if err := json.Unmarshal(body, &resultadosAPI); err != nil {
		// Tentar unmarshal como objeto único se a lista falhar (algumas APIs se comportam assim)
		var resultadoUnicoAPI DadosAnvisaAPI
		if errSingle := json.Unmarshal(body, &resultadoUnicoAPI); errSingle == nil {
			if resultadoUnicoAPI.Registro == codigoRegistro || (resultadoUnicoAPI.Registro == "" && resultadoUnicoAPI.Produto != "") { // Considerar válido se houver dados
				resultadosAPI = append(resultadosAPI, resultadoUnicoAPI)
			}
		} else {
			log.Printf("Erro ao decodificar JSON da API ANVISA: %v. Corpo: %s", err, string(body))
			return nil, fmt.Errorf("erro ao decodificar resposta da API: %w. Detalhe single: %w", err, errSingle)
		}
	}

	if len(resultadosAPI) == 0 {
		log.Printf("Nenhum medicamento encontrado na API para o código: %s", codigoRegistro)
		return nil, errors.New("medicamento não encontrado na API da ANVISA")
	}

	// Assumimos que o primeiro resultado é o correto se a busca for por código de registro
	// Em um cenário ideal, iteraríamos e confirmaríamos o código de registro.
	apiMed := resultadosAPI[0]

	// Mapear para nossa struct interna DadosAnvisa
	nossoMed := &DadosAnvisa{
		Nome:       apiMed.Produto,
		Fabricante: apiMed.Laboratorio,
		Registro:   apiMed.Registro, // Garantir que este é o código ANVISA
		Classe:     apiMed.ClasseTerapeutica,
	}

	// Se o campo Registro da API estiver vazio, mas recebemos o produto, preenchemos com o código buscado.
	// Isso é uma heurística caso a API não retorne o campo 'REGISTRO' de forma consistente.
	if nossoMed.Registro == "" && nossoMed.Nome != "" {
		nossoMed.Registro = codigoRegistro
	}

	log.Printf("Dados recuperados da API para %s: %+v", codigoRegistro, nossoMed)
	return nossoMed, nil
}

// AdicionarMedicamentoMock não é mais necessário com a API real.
// func AdicionarMedicamentoMock(codigo string, dados DadosAnvisa) {
// 	// Esta função seria para o mock, que não estamos mais usando primariamente.
// }

// Funções de cache podem ser reavaliadas depois, por enquanto vamos remover para focar na API.
// const cacheDuration = 24 * time.Hour
// func getCachePath() (string, error) { ... }
// func buscarDoCache(codigo string) (*DadosAnvisa, error) { ... }

// func salvarNoCache(codigo string, dados *DadosAnvisa) error { ... }
