package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// BulaInfo representa as informações da bula de um medicamento
type BulaInfo struct {
	Nome              string `json:"nome"`
	Indicacoes        string `json:"indicacoes"`
	Contraindicacoes  string `json:"contraindicacoes"`
	Posologia         string `json:"posologia"`
	EfeitosColaterais string `json:"efeitosColaterais"`
	Laboratorio       string `json:"laboratorio"`
	Registro          string `json:"registro"`
}

// AnvisaResponse representa a resposta da API da Anvisa
type AnvisaResponse struct {
	Content []struct {
		NumeroRegistro        string `json:"numeroRegistro"`
		NomeProduto           string `json:"nomeProduto"`
		Empresa               string `json:"empresa"`
		Processo              string `json:"processo"`
		SituacaoRegistro      string `json:"situacaoRegistro"`
		VencimentoRegistro    string `json:"vencimentoRegistro"`
		ClasseTerapeutica     string `json:"classeTerapeutica"`
		PrincipioAtivo        string `json:"principioAtivo"`
		MedicamentoReferencia string `json:"medicamentoReferencia"`
	} `json:"content"`
}

// GetBulaHandler retorna as informações da bula de um medicamento
func GetBulaHandler(w http.ResponseWriter, r *http.Request) {
	// Configura os headers
	w.Header().Set("Content-Type", "application/json")

	// Extrai o nome do medicamento da URL
	nomeMedicamento := strings.TrimPrefix(r.URL.Path, "/api/bula/")
	nomeMedicamento = strings.ReplaceAll(nomeMedicamento, "-", " ")

	// Busca as informações do medicamento na API da Anvisa
	bula, err := buscarInformacoesAnvisa(nomeMedicamento)
	if err != nil {
		// Se não encontrar na API da Anvisa, tenta buscar no mock
		mockBula, mockErr := buscarInformacoesMock(nomeMedicamento)
		if mockErr != nil {
			http.Error(w, "Bula não encontrada", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(mockBula)
		return
	}

	// Retorna as informações em JSON
	json.NewEncoder(w).Encode(bula)
}

// buscarInformacoesAnvisa busca informações do medicamento na API da Anvisa
func buscarInformacoesAnvisa(nome string) (*BulaInfo, error) {
	// URL da API da Anvisa (Consulta de Medicamentos)
	baseURL := "https://consultas.anvisa.gov.br/api/consulta/medicamento/produtos/"

	// Codifica o nome do medicamento para a URL
	query := url.QueryEscape(nome)

	// Faz a requisição para a API
	resp, err := http.Get(baseURL + "?nome=" + query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Se não encontrou o medicamento
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("medicamento não encontrado na base da Anvisa")
	}

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decodifica a resposta
	var anvisaResp AnvisaResponse
	if err := json.Unmarshal(body, &anvisaResp); err != nil {
		return nil, err
	}

	// Se não encontrou nenhum medicamento
	if len(anvisaResp.Content) == 0 {
		return nil, fmt.Errorf("medicamento não encontrado")
	}

	// Pega o primeiro resultado
	med := anvisaResp.Content[0]

	// Monta as informações da bula
	bula := &BulaInfo{
		Nome:              med.NomeProduto,
		Laboratorio:       med.Empresa,
		Registro:          med.NumeroRegistro,
		Indicacoes:        fmt.Sprintf("Classe Terapêutica: %s\nPrincípio Ativo: %s", med.ClasseTerapeutica, med.PrincipioAtivo),
		Contraindicacoes:  "Consulte a bula ou um profissional de saúde para informações sobre contraindicações.",
		Posologia:         "A posologia deve ser definida pelo médico de acordo com a condição do paciente.",
		EfeitosColaterais: "Consulte a bula ou um profissional de saúde para informações sobre possíveis efeitos colaterais.",
	}

	return bula, nil
}

// buscarInformacoesMock é um backup para quando a API da Anvisa não retorna resultados
func buscarInformacoesMock(nome string) (*BulaInfo, error) {
	// Converte o nome para minúsculas para comparação
	nomeLower := strings.ToLower(nome)

	// Mock de alguns medicamentos comuns
	medicamentos := map[string]*BulaInfo{
		"dipirona": {
			Nome:              "Dipirona",
			Indicacoes:        "Indicado para dor e febre.",
			Contraindicacoes:  "Alergia à dipirona ou outros analgésicos. Problemas de medula óssea.",
			Posologia:         "Adultos e adolescentes acima de 15 anos: 500-1000mg até 4 vezes ao dia.",
			EfeitosColaterais: "Reações alérgicas, problemas gastrointestinais, dor de cabeça.",
			Laboratorio:       "Genérico",
		},
		"paracetamol": {
			Nome:              "Paracetamol",
			Indicacoes:        "Dores leves a moderadas e febre.",
			Contraindicacoes:  "Doença hepática grave, alergia ao paracetamol.",
			Posologia:         "Adultos: 500-1000mg a cada 4-6 horas, não excedendo 4g por dia.",
			EfeitosColaterais: "Raramente causa reações alérgicas. Em doses altas pode causar dano hepático.",
			Laboratorio:       "Genérico",
		},
		"omeprazol": {
			Nome:              "Omeprazol",
			Indicacoes:        "Úlceras gástricas e duodenais, refluxo gastroesofágico.",
			Contraindicacoes:  "Hipersensibilidade ao omeprazol.",
			Posologia:         "20-40mg uma vez ao dia, preferencialmente pela manhã.",
			EfeitosColaterais: "Dor de cabeça, diarreia, constipação, dor abdominal.",
			Laboratorio:       "Genérico",
		},
	}

	// Procura o medicamento no mock
	for nomeMock, info := range medicamentos {
		if strings.Contains(nomeLower, nomeMock) {
			return info, nil
		}
	}

	// Se não encontrou, retorna erro
	return nil, fmt.Errorf("bula não encontrada para: %s", nome)
}
