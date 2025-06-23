package bula

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
	bula, err := BuscarInformacoesAnvisa(nomeMedicamento)
	if err != nil {
		// Se não encontrar na API da Anvisa, tenta buscar no mock
		mockBula, mockErr := BuscarInformacoesMock(nomeMedicamento)
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

// BuscarInformacoesAnvisa busca informações do medicamento na API da Anvisa
func BuscarInformacoesAnvisa(nome string) (*BulaInfo, error) {
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

// BuscarInformacoesMock é um backup para quando a API da Anvisa não retorna resultados
func BuscarInformacoesMock(nome string) (*BulaInfo, error) {
	// Converte o nome para minúsculas para comparação
	nomeLower := strings.ToLower(nome)

	// Mock de medicamentos
	medicamentos := map[string]*BulaInfo{
		"dipirona": {
			Nome:              "Dipirona 500mg",
			Indicacoes:        "Indicado para dor e febre. Eficaz no tratamento de dores de cabeça, dores musculares e estados febris.",
			Contraindicacoes:  "Alergia à dipirona ou outros analgésicos. Problemas de medula óssea. Pacientes com deficiência de G6PD.",
			Posologia:         "Adultos e adolescentes acima de 15 anos: 500-1000mg até 4 vezes ao dia, não excedendo 4g/dia.",
			EfeitosColaterais: "Reações alérgicas, problemas gastrointestinais, dor de cabeça, tontura, sonolência.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.000-0",
		},
		"omeprazol": {
			Nome:              "Omeprazol 20mg",
			Indicacoes:        "Úlceras gástricas e duodenais, refluxo gastroesofágico, síndrome de Zollinger-Ellison.",
			Contraindicacoes:  "Hipersensibilidade ao omeprazol ou outros inibidores da bomba de prótons.",
			Posologia:         "20-40mg uma vez ao dia, preferencialmente pela manhã, por 4-8 semanas.",
			EfeitosColaterais: "Dor de cabeça, diarreia, constipação, dor abdominal, náusea, gases.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.000-1",
		},
		"amoxicilina": {
			Nome:              "Amoxicilina 500mg",
			Indicacoes:        "Antibiótico para tratamento de infecções bacterianas do trato respiratório, geniturinário, pele e tecidos moles.",
			Contraindicacoes:  "Alergia a penicilinas ou cefalosporinas. Histórico de reações alérgicas graves.",
			Posologia:         "500mg a cada 8 horas ou conforme prescrição médica, por 7-10 dias.",
			EfeitosColaterais: "Diarreia, náusea, vômito, reações alérgicas cutâneas, candidíase.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.000-2",
		},
		"sertralina": {
			Nome:              "Sertralina 50mg",
			Indicacoes:        "Tratamento de depressão, transtorno obsessivo-compulsivo, pânico e ansiedade social.",
			Contraindicacoes:  "Uso de IMAO nos últimos 14 dias. Gravidez e lactação requerem avaliação médica.",
			Posologia:         "Iniciar com 50mg/dia, podendo ser ajustada conforme necessidade até 200mg/dia.",
			EfeitosColaterais: "Insônia, tontura, sonolência, náusea, diarreia, boca seca, disfunção sexual.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.000-3",
		},
		"alprazolam": {
			Nome:              "Alprazolam 1mg",
			Indicacoes:        "Tratamento de estados de ansiedade, incluindo ansiedade associada à depressão e transtorno do pânico.",
			Contraindicacoes:  "Glaucoma de ângulo estreito, miastenia gravis, insuficiência respiratória grave.",
			Posologia:         "0,5 a 1mg, 3 vezes ao dia. A dose pode ser aumentada gradualmente conforme necessidade.",
			EfeitosColaterais: "Sonolência, fadiga, memória prejudicada, tontura, depressão, dor de cabeça.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.000-4",
		},
		"paracetamol": {
			Nome:              "Paracetamol 750mg",
			Indicacoes:        "Analgésico e antitérmico indicado para dores leves a moderadas e febre.",
			Contraindicacoes:  "Doença hepática grave, hipersensibilidade ao paracetamol.",
			Posologia:         "750mg a cada 6 horas, não excedendo 4000mg por dia.",
			EfeitosColaterais: "Raramente causa reações alérgicas. Em doses altas pode afetar o fígado.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.000-5",
		},
		"metformina": {
			Nome:              "Metformina 850mg",
			Indicacoes:        "Controle da diabetes tipo 2, especialmente em pacientes com sobrepeso.",
			Contraindicacoes:  "Insuficiência renal ou hepática grave, acidose metabólica aguda.",
			Posologia:         "850mg 2 a 3 vezes ao dia, com as refeições.",
			EfeitosColaterais: "Desconforto gastrointestinal, náusea, diarreia, gosto metálico.",
			Laboratorio:       "Prati-Donaduzzi",
			Registro:          "MS 1.0000.0000.000-6",
		},
		"ibuprofeno": {
			Nome:              "Ibuprofeno 600mg",
			Indicacoes:        "Anti-inflamatório indicado para dor, febre e inflamação. Eficaz em dores musculares, artrite e cólicas.",
			Contraindicacoes:  "Úlcera gástrica ativa, insuficiência cardíaca grave, último trimestre de gravidez.",
			Posologia:         "600mg a cada 6-8 horas, não excedendo 2400mg por dia.",
			EfeitosColaterais: "Dor estomacal, náusea, dor de cabeça, tontura, retenção de líquidos.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.000-7",
		},
		"atenolol": {
			Nome:              "Atenolol 25mg",
			Indicacoes:        "Beta-bloqueador para tratamento da hipertensão e arritmias cardíacas.",
			Contraindicacoes:  "Bradicardia sinusal, bloqueio cardíaco, insuficiência cardíaca descompensada.",
			Posologia:         "25-100mg uma vez ao dia.",
			EfeitosColaterais: "Fadiga, extremidades frias, bradicardia, tontura, distúrbios do sono.",
			Laboratorio:       "Biosintética",
			Registro:          "MS 1.0000.0000.000-8",
		},
		"nimesulida": {
			Nome:              "Nimesulida 100mg",
			Indicacoes:        "Anti-inflamatório não esteroidal para dor e inflamação.",
			Contraindicacoes:  "Úlcera péptica ativa, insuficiência hepática grave, crianças menores de 12 anos.",
			Posologia:         "100mg duas vezes ao dia, após as refeições.",
			EfeitosColaterais: "Distúrbios gastrointestinais, dor de cabeça, tontura.",
			Laboratorio:       "Cimed",
			Registro:          "MS 1.0000.0000.000-9",
		},
		"levotiroxina": {
			Nome:              "Levotiroxina 50mcg",
			Indicacoes:        "Tratamento do hipotireoidismo e supressão do TSH em algumas condições.",
			Contraindicacoes:  "Hipertireoidismo não tratado, infarto agudo do miocárdio recente.",
			Posologia:         "Dose individualizada, geralmente tomada em jejum.",
			EfeitosColaterais: "Taquicardia, tremores, insônia, perda de peso se dose excessiva.",
			Laboratorio:       "Merck",
			Registro:          "MS 1.0000.0000.001-0",
		},
		"fluoxetina": {
			Nome:              "Fluoxetina 20mg",
			Indicacoes:        "Antidepressivo para tratamento de depressão, TOC, bulimia nervosa.",
			Contraindicacoes:  "Uso de IMAO nos últimos 14 dias, hipersensibilidade à fluoxetina.",
			Posologia:         "20mg uma vez ao dia, pela manhã.",
			EfeitosColaterais: "Náusea, insônia, ansiedade, diminuição do apetite.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.001-1",
		},
		"sinvastatina": {
			Nome:              "Sinvastatina 20mg",
			Indicacoes:        "Redução do colesterol e prevenção de doenças cardiovasculares.",
			Contraindicacoes:  "Doença hepática ativa, gravidez, amamentação.",
			Posologia:         "20-40mg uma vez ao dia, à noite.",
			EfeitosColaterais: "Dores musculares, alterações hepáticas, dor de cabeça.",
			Laboratorio:       "Sandoz",
			Registro:          "MS 1.0000.0000.001-2",
		},
		"cefalexina": {
			Nome:              "Cefalexina 500mg",
			Indicacoes:        "Antibiótico para infecções bacterianas diversas.",
			Contraindicacoes:  "Alergia a cefalosporinas ou penicilinas.",
			Posologia:         "500mg a cada 6-8 horas por 7-10 dias.",
			EfeitosColaterais: "Diarreia, náusea, reações alérgicas, candidíase.",
			Laboratorio:       "Teuto",
			Registro:          "MS 1.0000.0000.001-3",
		},
		"escitalopram": {
			Nome:              "Escitalopram 10mg",
			Indicacoes:        "Tratamento de depressão e transtornos de ansiedade.",
			Contraindicacoes:  "Uso concomitante com IMAO, hipersensibilidade ao escitalopram.",
			Posologia:         "10mg uma vez ao dia, podendo ser aumentada para 20mg.",
			EfeitosColaterais: "Náusea, insônia, sudorese, disfunção sexual.",
			Laboratorio:       "Eurofarma",
			Registro:          "MS 1.0000.0000.001-4",
		},
		"pantoprazol": {
			Nome:              "Pantoprazol 40mg",
			Indicacoes:        "Tratamento de úlceras e refluxo gastroesofágico.",
			Contraindicacoes:  "Hipersensibilidade ao pantoprazol.",
			Posologia:         "40mg uma vez ao dia, antes do café da manhã.",
			EfeitosColaterais: "Dor de cabeça, diarreia, náusea, tonturas.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.001-5",
		},
		"dexametasona": {
			Nome:              "Dexametasona 4mg",
			Indicacoes:        "Anti-inflamatório e imunossupressor para várias condições.",
			Contraindicacoes:  "Infecções sistêmicas sem tratamento adequado.",
			Posologia:         "Dose varia conforme a condição, seguir prescrição médica.",
			EfeitosColaterais: "Retenção de líquidos, aumento da pressão arterial, alterações metabólicas.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.001-6",
		},
		"azitromicina": {
			Nome:              "Azitromicina 500mg",
			Indicacoes:        "Antibiótico para infecções bacterianas do trato respiratório e outras.",
			Contraindicacoes:  "Hipersensibilidade à azitromicina ou outros macrolídeos.",
			Posologia:         "500mg uma vez ao dia por 3-5 dias.",
			EfeitosColaterais: "Diarreia, náusea, dor abdominal, alterações do paladar.",
			Laboratorio:       "Prati-Donaduzzi",
			Registro:          "MS 1.0000.0000.001-7",
		},
		"losartana": {
			Nome:              "Losartana 50mg",
			Indicacoes:        "Tratamento da hipertensão arterial e proteção renal em diabéticos.",
			Contraindicacoes:  "Gravidez, hipersensibilidade à losartana.",
			Posologia:         "50-100mg uma vez ao dia.",
			EfeitosColaterais: "Tontura, hipotensão, alterações nos níveis de potássio.",
			Laboratorio:       "Eurofarma",
			Registro:          "MS 1.0000.0000.001-8",
		},
		"enalapril": {
			Nome:              "Enalapril 10mg",
			Indicacoes:        "Tratamento da hipertensão arterial e insuficiência cardíaca.",
			Contraindicacoes:  "Histórico de angioedema, gravidez.",
			Posologia:         "10-40mg por dia, em uma ou duas tomadas.",
			EfeitosColaterais: "Tosse seca, tontura, alterações no paladar.",
			Laboratorio:       "Biosintética",
			Registro:          "MS 1.0000.0000.001-9",
		},
		"clonazepam": {
			Nome:              "Clonazepam 2mg",
			Indicacoes:        "Tratamento de crises epilépticas e transtornos de ansiedade.",
			Contraindicacoes:  "Glaucoma agudo, insuficiência respiratória grave, gravidez.",
			Posologia:         "Dose inicial de 0,5mg, podendo ser aumentada conforme necessidade.",
			EfeitosColaterais: "Sonolência, tontura, fadiga, alterações de coordenação.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.002-0",
		},
		"rosuvastatina": {
			Nome:              "Rosuvastatina 10mg",
			Indicacoes:        "Redução do colesterol e prevenção cardiovascular.",
			Contraindicacoes:  "Doença hepática ativa, gravidez, amamentação.",
			Posologia:         "5-40mg uma vez ao dia.",
			EfeitosColaterais: "Dores musculares, dor de cabeça, alterações hepáticas.",
			Laboratorio:       "Sandoz",
			Registro:          "MS 1.0000.0000.002-1",
		},
		"venlafaxina": {
			Nome:              "Venlafaxina 75mg",
			Indicacoes:        "Tratamento da depressão e ansiedade generalizada.",
			Contraindicacoes:  "Uso de IMAO nos últimos 14 dias, hipertensão não controlada.",
			Posologia:         "75mg por dia, podendo ser aumentada gradualmente.",
			EfeitosColaterais: "Náusea, insônia, tontura, sudorese, hipertensão.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.002-2",
		},
		"pregabalina": {
			Nome:              "Pregabalina 75mg",
			Indicacoes:        "Tratamento da dor neuropática e fibromialgia.",
			Contraindicacoes:  "Hipersensibilidade à pregabalina.",
			Posologia:         "75mg duas vezes ao dia, podendo ser aumentada.",
			EfeitosColaterais: "Tontura, sonolência, ganho de peso, visão turva.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.002-3",
		},
		"duloxetina": {
			Nome:              "Duloxetina 30mg",
			Indicacoes:        "Depressão, ansiedade, dor neuropática diabética.",
			Contraindicacoes:  "Uso de IMAO, glaucoma de ângulo fechado.",
			Posologia:         "30-60mg uma vez ao dia.",
			EfeitosColaterais: "Náusea, boca seca, insônia, fadiga.",
			Laboratorio:       "Eurofarma",
			Registro:          "MS 1.0000.0000.002-4",
		},
		"carvedilol": {
			Nome:              "Carvedilol 6.25mg",
			Indicacoes:        "Tratamento da insuficiência cardíaca e hipertensão.",
			Contraindicacoes:  "Asma brônquica, bloqueio cardíaco avançado.",
			Posologia:         "6.25mg duas vezes ao dia, podendo ser aumentada.",
			EfeitosColaterais: "Tontura, fadiga, bradicardia, hipotensão.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.002-5",
		},
		"anlodipino": {
			Nome:              "Anlodipino 5mg",
			Indicacoes:        "Tratamento da hipertensão e angina.",
			Contraindicacoes:  "Hipersensibilidade ao anlodipino.",
			Posologia:         "5-10mg uma vez ao dia.",
			EfeitosColaterais: "Edema, dor de cabeça, rubor facial.",
			Laboratorio:       "Biosintética",
			Registro:          "MS 1.0000.0000.002-6",
		},
		"citalopram": {
			Nome:              "Citalopram 20mg",
			Indicacoes:        "Tratamento da depressão e transtorno do pânico.",
			Contraindicacoes:  "Uso concomitante com IMAO.",
			Posologia:         "20mg uma vez ao dia, podendo ser aumentada.",
			EfeitosColaterais: "Náusea, boca seca, sonolência, sudorese.",
			Laboratorio:       "Germed",
			Registro:          "MS 1.0000.0000.002-7",
		},
		"metoprolol": {
			Nome:              "Metoprolol 50mg",
			Indicacoes:        "Hipertensão, angina, arritmias.",
			Contraindicacoes:  "Bloqueio cardíaco, insuficiência cardíaca não tratada.",
			Posologia:         "50-200mg por dia, em doses divididas.",
			EfeitosColaterais: "Fadiga, bradicardia, tontura, depressão.",
			Laboratorio:       "Astrazeneca",
			Registro:          "MS 1.0000.0000.002-8",
		},
		"gabapentina": {
			Nome:              "Gabapentina 300mg",
			Indicacoes:        "Epilepsia e dor neuropática.",
			Contraindicacoes:  "Hipersensibilidade à gabapentina.",
			Posologia:         "300-1200mg três vezes ao dia.",
			EfeitosColaterais: "Sonolência, tontura, ataxia, fadiga.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.002-9",
		},
		"bromazepam": {
			Nome:              "Bromazepam 3mg",
			Indicacoes:        "Ansiedade, tensão e outras queixas somáticas.",
			Contraindicacoes:  "Miastenia gravis, glaucoma agudo.",
			Posologia:         "1,5-3mg duas ou três vezes ao dia.",
			EfeitosColaterais: "Sonolência, confusão, fraqueza muscular.",
			Laboratorio:       "Roche",
			Registro:          "MS 1.0000.0000.003-0",
		},
		"domperidona": {
			Nome:              "Domperidona 10mg",
			Indicacoes:        "Náuseas, vômitos e refluxo gastroesofágico.",
			Contraindicacoes:  "Prolactinoma, hemorragia gastrointestinal.",
			Posologia:         "10mg até três vezes ao dia.",
			EfeitosColaterais: "Boca seca, dor de cabeça, alterações menstruais.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.003-1",
		},
		"ranitidina": {
			Nome:              "Ranitidina 150mg",
			Indicacoes:        "Úlcera péptica, refluxo gastroesofágico.",
			Contraindicacoes:  "Hipersensibilidade à ranitidina.",
			Posologia:         "150mg duas vezes ao dia ou 300mg à noite.",
			EfeitosColaterais: "Dor de cabeça, constipação, diarreia.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.003-2",
		},
		"doxazosina": {
			Nome:              "Doxazosina 2mg",
			Indicacoes:        "Hipertensão e hiperplasia prostática benigna.",
			Contraindicacoes:  "Hipotensão ortostática.",
			Posologia:         "1-8mg uma vez ao dia.",
			EfeitosColaterais: "Tontura, fadiga, edema, hipotensão.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.003-3",
		},
		"finasterida": {
			Nome:              "Finasterida 1mg",
			Indicacoes:        "Tratamento da calvície masculina.",
			Contraindicacoes:  "Mulheres, crianças e adolescentes.",
			Posologia:         "1mg uma vez ao dia.",
			EfeitosColaterais: "Diminuição da libido, disfunção erétil.",
			Laboratorio:       "Merck",
			Registro:          "MS 1.0000.0000.003-4",
		},
		"ciclobenzaprina": {
			Nome:              "Ciclobenzaprina 5mg",
			Indicacoes:        "Relaxante muscular para espasmos musculares.",
			Contraindicacoes:  "Uso de IMAO, arritmias cardíacas.",
			Posologia:         "5-10mg três vezes ao dia.",
			EfeitosColaterais: "Sonolência, boca seca, tontura.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.003-5",
		},
		"cetoprofeno": {
			Nome:              "Cetoprofeno 100mg",
			Indicacoes:        "Dor e inflamação em condições reumáticas.",
			Contraindicacoes:  "Úlcera péptica, asma grave.",
			Posologia:         "100mg duas vezes ao dia.",
			EfeitosColaterais: "Dor estomacal, náusea, dor de cabeça.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.003-6",
		},
		"hidroclorotiazida": {
			Nome:              "Hidroclorotiazida 25mg",
			Indicacoes:        "Hipertensão e edema.",
			Contraindicacoes:  "Anúria, hipersensibilidade a sulfonamidas.",
			Posologia:         "25-50mg uma vez ao dia.",
			EfeitosColaterais: "Desequilíbrio eletrolítico, hipotensão.",
			Laboratorio:       "Prati-Donaduzzi",
			Registro:          "MS 1.0000.0000.003-7",
		},
		"clopidogrel": {
			Nome:              "Clopidogrel 75mg",
			Indicacoes:        "Prevenção de eventos aterotrombóticos.",
			Contraindicacoes:  "Sangramento ativo, insuficiência hepática grave.",
			Posologia:         "75mg uma vez ao dia.",
			EfeitosColaterais: "Sangramento, dor abdominal, diarreia.",
			Laboratorio:       "Sandoz",
			Registro:          "MS 1.0000.0000.003-8",
		},
		"glimepirida": {
			Nome:              "Glimepirida 2mg",
			Indicacoes:        "Diabetes mellitus tipo 2.",
			Contraindicacoes:  "Diabetes tipo 1, cetoacidose diabética.",
			Posologia:         "1-4mg uma vez ao dia antes do café.",
			EfeitosColaterais: "Hipoglicemia, ganho de peso, náusea.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.003-9",
		},
		"insulina": {
			Nome:              "Insulina NPH 100UI/mL",
			Indicacoes:        "Diabetes mellitus tipo 1 e 2.",
			Contraindicacoes:  "Hipoglicemia, hipersensibilidade à insulina NPH.",
			Posologia:         "Dose individualizada conforme necessidade do paciente.",
			EfeitosColaterais: "Hipoglicemia, ganho de peso, reações no local da aplicação.",
			Laboratorio:       "Novo Nordisk",
			Registro:          "MS 1.0000.0000.004-0",
		},
		"gliclazida": {
			Nome:              "Gliclazida 30mg MR",
			Indicacoes:        "Diabetes mellitus tipo 2.",
			Contraindicacoes:  "Diabetes tipo 1, cetoacidose, insuficiência renal grave.",
			Posologia:         "30-120mg uma vez ao dia no café da manhã.",
			EfeitosColaterais: "Hipoglicemia, distúrbios gastrointestinais.",
			Laboratorio:       "Servier",
			Registro:          "MS 1.0000.0000.004-1",
		},
		"montelucaste": {
			Nome:              "Montelucaste 10mg",
			Indicacoes:        "Asma e rinite alérgica.",
			Contraindicacoes:  "Hipersensibilidade ao montelucaste.",
			Posologia:         "10mg uma vez ao dia à noite.",
			EfeitosColaterais: "Dor de cabeça, dor abdominal, sede.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.004-2",
		},
		"budesonida": {
			Nome:              "Budesonida Spray Nasal 32mcg",
			Indicacoes:        "Rinite alérgica e não alérgica.",
			Contraindicacoes:  "Infecções nasais não tratadas, tuberculose nasal.",
			Posologia:         "1-2 jatos em cada narina uma vez ao dia.",
			EfeitosColaterais: "Irritação nasal, sangramento nasal, dor de cabeça.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.004-3",
		},
		"salbutamol": {
			Nome:              "Salbutamol Spray 100mcg",
			Indicacoes:        "Asma, bronquite e outras doenças respiratórias.",
			Contraindicacoes:  "Hipersensibilidade ao salbutamol.",
			Posologia:         "1-2 jatos até 4 vezes ao dia quando necessário.",
			EfeitosColaterais: "Tremor, taquicardia, nervosismo.",
			Laboratorio:       "GSK",
			Registro:          "MS 1.0000.0000.004-4",
		},
		"risperidona": {
			Nome:              "Risperidona 2mg",
			Indicacoes:        "Esquizofrenia e transtorno bipolar.",
			Contraindicacoes:  "Hipersensibilidade à risperidona.",
			Posologia:         "2-8mg por dia, divididos em duas doses.",
			EfeitosColaterais: "Sonolência, aumento de peso, alterações hormonais.",
			Laboratorio:       "Eurofarma",
			Registro:          "MS 1.0000.0000.004-5",
		},
		"quetiapina": {
			Nome:              "Quetiapina 25mg",
			Indicacoes:        "Esquizofrenia, transtorno bipolar, depressão maior.",
			Contraindicacoes:  "Hipersensibilidade à quetiapina.",
			Posologia:         "25-800mg por dia, divididos em duas doses.",
			EfeitosColaterais: "Sonolência, tontura, boca seca, ganho de peso.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.004-6",
		},
		"topiramato": {
			Nome:              "Topiramato 50mg",
			Indicacoes:        "Epilepsia e prevenção de enxaqueca.",
			Contraindicacoes:  "Gravidez, amamentação.",
			Posologia:         "25-200mg duas vezes ao dia.",
			EfeitosColaterais: "Perda de peso, parestesia, dificuldade de concentração.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.004-7",
		},
		"lamotrigina": {
			Nome:              "Lamotrigina 100mg",
			Indicacoes:        "Epilepsia e transtorno bipolar.",
			Contraindicacoes:  "Hipersensibilidade à lamotrigina.",
			Posologia:         "25-200mg por dia, em uma ou duas doses.",
			EfeitosColaterais: "Erupção cutânea, dor de cabeça, tontura, diplopia.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.004-8",
		},
		"oxcarbazepina": {
			Nome:              "Oxcarbazepina 300mg",
			Indicacoes:        "Epilepsia e neuralgia do trigêmeo.",
			Contraindicacoes:  "Hipersensibilidade à oxcarbazepina.",
			Posologia:         "300-1200mg duas vezes ao dia.",
			EfeitosColaterais: "Sonolência, tontura, náusea, visão dupla.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.004-9",
		},
		"levodopa": {
			Nome:              "Levodopa + Carbidopa 250/25mg",
			Indicacoes:        "Doença de Parkinson.",
			Contraindicacoes:  "Glaucoma de ângulo fechado, melanoma.",
			Posologia:         "3-4 comprimidos por dia, divididos em doses.",
			EfeitosColaterais: "Náusea, hipotensão postural, movimentos involuntários.",
			Laboratorio:       "Cristália",
			Registro:          "MS 1.0000.0000.005-0",
		},
		"memantina": {
			Nome:              "Memantina 10mg",
			Indicacoes:        "Doença de Alzheimer moderada a grave.",
			Contraindicacoes:  "Hipersensibilidade à memantina.",
			Posologia:         "5-20mg por dia.",
			EfeitosColaterais: "Tontura, dor de cabeça, constipação, sonolência.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.005-1",
		},
		"donepezila": {
			Nome:              "Donepezila 5mg",
			Indicacoes:        "Doença de Alzheimer leve a moderada.",
			Contraindicacoes:  "Hipersensibilidade à donepezila.",
			Posologia:         "5-10mg uma vez ao dia.",
			EfeitosColaterais: "Náusea, diarreia, insônia, cãibras musculares.",
			Laboratorio:       "Cristália",
			Registro:          "MS 1.0000.0000.005-2",
		},
		"rivaroxabana": {
			Nome:              "Rivaroxabana 20mg",
			Indicacoes:        "Prevenção de trombose e AVC.",
			Contraindicacoes:  "Sangramento ativo clinicamente significativo.",
			Posologia:         "20mg uma vez ao dia com alimentos.",
			EfeitosColaterais: "Sangramento, náusea, anemia.",
			Laboratorio:       "Bayer",
			Registro:          "MS 1.0000.0000.005-3",
		},
		"apixabana": {
			Nome:              "Apixabana 5mg",
			Indicacoes:        "Prevenção de trombose e AVC.",
			Contraindicacoes:  "Sangramento ativo clinicamente significativo.",
			Posologia:         "5mg duas vezes ao dia.",
			EfeitosColaterais: "Sangramento, contusão, anemia.",
			Laboratorio:       "Pfizer",
			Registro:          "MS 1.0000.0000.005-4",
		},
		"varfarina": {
			Nome:              "Varfarina 5mg",
			Indicacoes:        "Anticoagulação em diversas condições.",
			Contraindicacoes:  "Sangramento ativo, gravidez.",
			Posologia:         "Dose ajustada individualmente pelo INR.",
			EfeitosColaterais: "Sangramento, necrose cutânea, alopecia.",
			Laboratorio:       "Farmoquímica",
			Registro:          "MS 1.0000.0000.005-5",
		},
		"cilostazol": {
			Nome:              "Cilostazol 100mg",
			Indicacoes:        "Claudicação intermitente.",
			Contraindicacoes:  "Insuficiência cardíaca, sangramento ativo.",
			Posologia:         "100mg duas vezes ao dia.",
			EfeitosColaterais: "Dor de cabeça, diarreia, palpitações.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.005-6",
		},
		"diosmina": {
			Nome:              "Diosmina + Hesperidina 450/50mg",
			Indicacoes:        "Insuficiência venosa crônica.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "2 comprimidos ao dia com as refeições.",
			EfeitosColaterais: "Distúrbios gastrointestinais leves.",
			Laboratorio:       "Servier",
			Registro:          "MS 1.0000.0000.005-7",
		},
		"ginkgo": {
			Nome:              "Ginkgo Biloba 80mg",
			Indicacoes:        "Distúrbios circulatórios cerebrais.",
			Contraindicacoes:  "Hipersensibilidade ao Ginkgo biloba.",
			Posologia:         "80-120mg duas a três vezes ao dia.",
			EfeitosColaterais: "Dor de cabeça, distúrbios gastrointestinais.",
			Laboratorio:       "Herbarium",
			Registro:          "MS 1.0000.0000.005-8",
		},
		"vitamina_d": {
			Nome:              "Vitamina D3 7000UI",
			Indicacoes:        "Prevenção e tratamento da deficiência de vitamina D.",
			Contraindicacoes:  "Hipervitaminose D, hipercalcemia.",
			Posologia:         "7000UI uma vez por semana.",
			EfeitosColaterais: "Hipercalcemia, náusea, vômito (em doses altas).",
			Laboratorio:       "Mantecorp",
			Registro:          "MS 1.0000.0000.005-9",
		},
		"complexo_b": {
			Nome:              "Complexo B",
			Indicacoes:        "Suplementação vitamínica do complexo B.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "1 comprimido ao dia.",
			EfeitosColaterais: "Urina com coloração amarela intensa, náusea leve.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.006-0",
		},
		"vitamina_c": {
			Nome:              "Vitamina C 1g",
			Indicacoes:        "Suplementação de vitamina C.",
			Contraindicacoes:  "Cálculos renais de oxalato, hemocromatose.",
			Posologia:         "1 comprimido ao dia.",
			EfeitosColaterais: "Diarreia em altas doses, acidez gástrica.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.006-1",
		},
		"centrum": {
			Nome:              "Centrum",
			Indicacoes:        "Suplementação de vitaminas e minerais.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "1 comprimido ao dia com refeição.",
			EfeitosColaterais: "Náusea, desconforto gástrico leve.",
			Laboratorio:       "Pfizer",
			Registro:          "MS 1.0000.0000.006-2",
		},
		"buscopan": {
			Nome:              "Buscopan Composto",
			Indicacoes:        "Cólicas e dores abdominais.",
			Contraindicacoes:  "Glaucoma, hipertrofia prostática.",
			Posologia:         "1 comprimido 3 a 4 vezes ao dia.",
			EfeitosColaterais: "Boca seca, visão turva, taquicardia.",
			Laboratorio:       "Boehringer",
			Registro:          "MS 1.0000.0000.006-3",
		},
		"dramin": {
			Nome:              "Dramin B6",
			Indicacoes:        "Náuseas, vômitos, vertigens.",
			Contraindicacoes:  "Glaucoma, epilepsia.",
			Posologia:         "1 comprimido até 3 vezes ao dia.",
			EfeitosColaterais: "Sonolência, boca seca.",
			Laboratorio:       "Takeda",
			Registro:          "MS 1.0000.0000.006-4",
		},
		"luftal": {
			Nome:              "Luftal",
			Indicacoes:        "Gases e flatulência.",
			Contraindicacoes:  "Hipersensibilidade à simeticona.",
			Posologia:         "40 gotas após as refeições.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "Reckitt",
			Registro:          "MS 1.0000.0000.006-5",
		},
		"dorflex": {
			Nome:              "Dorflex",
			Indicacoes:        "Dores musculares e tensionais.",
			Contraindicacoes:  "Glaucoma, problemas hepáticos graves.",
			Posologia:         "1 comprimido até 4 vezes ao dia.",
			EfeitosColaterais: "Sonolência, boca seca, tontura.",
			Laboratorio:       "Sanofi",
			Registro:          "MS 1.0000.0000.006-6",
		},
		"neosaldina": {
			Nome:              "Neosaldina",
			Indicacoes:        "Dores de cabeça e enxaqueca.",
			Contraindicacoes:  "Glaucoma, hipertensão grave.",
			Posologia:         "1 comprimido até 4 vezes ao dia.",
			EfeitosColaterais: "Sonolência, boca seca, taquicardia.",
			Laboratorio:       "Takeda",
			Registro:          "MS 1.0000.0000.006-7",
		},
		"enterogermina": {
			Nome:              "Enterogermina 5ml",
			Indicacoes:        "Desequilíbrio da flora intestinal.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "1 frasco 2 a 3 vezes ao dia.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "Sanofi",
			Registro:          "MS 1.0000.0000.006-8",
		},
		"lacto_purga": {
			Nome:              "Lacto-Purga",
			Indicacoes:        "Constipação intestinal.",
			Contraindicacoes:  "Obstrução intestinal, dor abdominal aguda.",
			Posologia:         "1-2 comprimidos à noite.",
			EfeitosColaterais: "Cólicas abdominais, diarreia.",
			Laboratorio:       "Cosmed",
			Registro:          "MS 1.0000.0000.006-9",
		},
		"bepantol": {
			Nome:              "Bepantol Derma",
			Indicacoes:        "Hidratação e regeneração da pele.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "Aplicar 2 a 3 vezes ao dia.",
			EfeitosColaterais: "Raramente causa reações alérgicas.",
			Laboratorio:       "Bayer",
			Registro:          "MS 1.0000.0000.007-0",
		},
		"episol": {
			Nome:              "Episol FPS 60",
			Indicacoes:        "Proteção solar UVA/UVB.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "Aplicar 15 minutos antes da exposição solar.",
			EfeitosColaterais: "Pode causar irritação em peles sensíveis.",
			Laboratorio:       "Mantecorp",
			Registro:          "MS 1.0000.0000.007-1",
		},
		"aas": {
			Nome:              "AAS 100mg",
			Indicacoes:        "Prevenção de eventos cardiovasculares.",
			Contraindicacoes:  "Úlcera gástrica ativa, hemofilia.",
			Posologia:         "100mg uma vez ao dia.",
			EfeitosColaterais: "Irritação gástrica, sangramento aumentado.",
			Laboratorio:       "Sanofi",
			Registro:          "MS 1.0000.0000.007-2",
		},
		"allegra": {
			Nome:              "Allegra 180mg",
			Indicacoes:        "Rinite alérgica, urticária.",
			Contraindicacoes:  "Hipersensibilidade à fexofenadina.",
			Posologia:         "1 comprimido uma vez ao dia.",
			EfeitosColaterais: "Dor de cabeça, sonolência, náusea.",
			Laboratorio:       "Sanofi",
			Registro:          "MS 1.0000.0000.007-3",
		},
		"loratadina": {
			Nome:              "Loratadina 10mg",
			Indicacoes:        "Rinite alérgica, urticária.",
			Contraindicacoes:  "Hipersensibilidade à loratadina.",
			Posologia:         "10mg uma vez ao dia.",
			EfeitosColaterais: "Sonolência, boca seca, fadiga.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.007-4",
		},
		"desloratadina": {
			Nome:              "Desloratadina 5mg",
			Indicacoes:        "Rinite alérgica, urticária.",
			Contraindicacoes:  "Hipersensibilidade à desloratadina.",
			Posologia:         "5mg uma vez ao dia.",
			EfeitosColaterais: "Fadiga, boca seca, dor de cabeça.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.007-5",
		},
		"polaramine": {
			Nome:              "Polaramine",
			Indicacoes:        "Alergias, prurido, rinite.",
			Contraindicacoes:  "Glaucoma, retenção urinária.",
			Posologia:         "1 comprimido a cada 4-6 horas.",
			EfeitosColaterais: "Sonolência, boca seca, visão turva.",
			Laboratorio:       "Mantecorp",
			Registro:          "MS 1.0000.0000.007-6",
		},
		"aciclovir": {
			Nome:              "Aciclovir 200mg",
			Indicacoes:        "Infecções por herpes simplex.",
			Contraindicacoes:  "Hipersensibilidade ao aciclovir.",
			Posologia:         "200mg 5 vezes ao dia por 5 dias.",
			EfeitosColaterais: "Dor de cabeça, náusea, diarreia.",
			Laboratorio:       "Medley",
			Registro:          "MS 1.0000.0000.007-7",
		},
		"tamiflu": {
			Nome:              "Tamiflu 75mg",
			Indicacoes:        "Tratamento e prevenção da influenza.",
			Contraindicacoes:  "Hipersensibilidade ao oseltamivir.",
			Posologia:         "75mg duas vezes ao dia por 5 dias.",
			EfeitosColaterais: "Náusea, vômito, dor de cabeça.",
			Laboratorio:       "Roche",
			Registro:          "MS 1.0000.0000.007-8",
		},
		"fluconazol": {
			Nome:              "Fluconazol 150mg",
			Indicacoes:        "Infecções fúngicas.",
			Contraindicacoes:  "Hipersensibilidade ao fluconazol.",
			Posologia:         "150mg dose única ou conforme prescrição.",
			EfeitosColaterais: "Dor de cabeça, náusea, dor abdominal.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.007-9",
		},
		"miconazol": {
			Nome:              "Miconazol Creme",
			Indicacoes:        "Infecções fúngicas cutâneas.",
			Contraindicacoes:  "Hipersensibilidade ao miconazol.",
			Posologia:         "Aplicar 2 vezes ao dia por 2-4 semanas.",
			EfeitosColaterais: "Irritação local, ardência.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.008-0",
		},
		"cetoconazol": {
			Nome:              "Cetoconazol Shampoo",
			Indicacoes:        "Caspa, dermatite seborreica.",
			Contraindicacoes:  "Hipersensibilidade ao cetoconazol.",
			Posologia:         "Aplicar 2 vezes por semana por 2-4 semanas.",
			EfeitosColaterais: "Irritação local, ressecamento.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.008-1",
		},
		"benegrip": {
			Nome:              "Benegrip Multi",
			Indicacoes:        "Sintomas de gripe e resfriado.",
			Contraindicacoes:  "Glaucoma, hipertensão grave.",
			Posologia:         "1 comprimido a cada 6 horas.",
			EfeitosColaterais: "Sonolência, boca seca, tontura.",
			Laboratorio:       "Hypera",
			Registro:          "MS 1.0000.0000.008-2",
		},
		"cimegripe": {
			Nome:              "Cimegripe",
			Indicacoes:        "Sintomas de gripe e resfriado.",
			Contraindicacoes:  "Glaucoma, problemas cardíacos graves.",
			Posologia:         "1 comprimido a cada 6 horas.",
			EfeitosColaterais: "Sonolência, tontura, boca seca.",
			Laboratorio:       "Cimed",
			Registro:          "MS 1.0000.0000.008-3",
		},
		"vick": {
			Nome:              "Vick VapoRub",
			Indicacoes:        "Congestão nasal, tosse.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "Aplicar no peito e garganta à noite.",
			EfeitosColaterais: "Irritação local em peles sensíveis.",
			Laboratorio:       "P&G",
			Registro:          "MS 1.0000.0000.008-4",
		},
		"coristina": {
			Nome:              "Coristina D",
			Indicacoes:        "Sintomas de gripe e resfriado.",
			Contraindicacoes:  "Glaucoma, hipertensão grave.",
			Posologia:         "1 comprimido a cada 4 horas.",
			EfeitosColaterais: "Sonolência, boca seca, nervosismo.",
			Laboratorio:       "Mantecorp",
			Registro:          "MS 1.0000.0000.008-5",
		},
		"biotonico": {
			Nome:              "Biotônico Fontoura",
			Indicacoes:        "Suplemento fortificante.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "10-20ml, 2 a 3 vezes ao dia.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "Hypera",
			Registro:          "MS 1.0000.0000.008-6",
		},
		"gelol": {
			Nome:              "Gelol",
			Indicacoes:        "Dores musculares e contusões.",
			Contraindicacoes:  "Feridas abertas, queimaduras.",
			Posologia:         "Aplicar 3 a 4 vezes ao dia.",
			EfeitosColaterais: "Irritação local, vermelhidão.",
			Laboratorio:       "Hypera",
			Registro:          "MS 1.0000.0000.008-7",
		},
		"dexametasona_creme": {
			Nome:              "Dexametasona Creme",
			Indicacoes:        "Inflamações e alergias cutâneas.",
			Contraindicacoes:  "Infecções cutâneas não tratadas.",
			Posologia:         "Aplicar 2 a 3 vezes ao dia.",
			EfeitosColaterais: "Atrofia cutânea, manchas brancas.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.008-8",
		},
		"tandrilax": {
			Nome:              "Tandrilax",
			Indicacoes:        "Dores musculares e inflamação.",
			Contraindicacoes:  "Úlcera gástrica, problemas hepáticos.",
			Posologia:         "1 comprimido 3 vezes ao dia.",
			EfeitosColaterais: "Sonolência, tontura, dor estomacal.",
			Laboratorio:       "Cimed",
			Registro:          "MS 1.0000.0000.008-9",
		},
		"clorexidina": {
			Nome:              "Clorexidina 0.12%",
			Indicacoes:        "Antisséptico bucal.",
			Contraindicacoes:  "Hipersensibilidade à clorexidina.",
			Posologia:         "Bochechar por 30 segundos, 2 vezes ao dia.",
			EfeitosColaterais: "Manchas nos dentes, alteração do paladar.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.009-0",
		},
		"listerine": {
			Nome:              "Listerine Cool Mint",
			Indicacoes:        "Higiene bucal, mau hálito.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "Bochechar por 30 segundos, 2 vezes ao dia.",
			EfeitosColaterais: "Sensação de queimação, alteração do paladar.",
			Laboratorio:       "Johnson & Johnson",
			Registro:          "MS 1.0000.0000.009-1",
		},
		"estomazil": {
			Nome:              "Estomazil",
			Indicacoes:        "Azia, má digestão.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "1 envelope dissolvido em água quando necessário.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "Hypera",
			Registro:          "MS 1.0000.0000.009-2",
		},
		"eno": {
			Nome:              "Eno",
			Indicacoes:        "Azia, má digestão.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "1 envelope dissolvido em água quando necessário.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "GSK",
			Registro:          "MS 1.0000.0000.009-3",
		},
		"sal_de_fruta": {
			Nome:              "Sal de Fruta",
			Indicacoes:        "Azia, má digestão.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "1 envelope dissolvido em água quando necessário.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "EMS",
			Registro:          "MS 1.0000.0000.009-4",
		},
		"epocler": {
			Nome:              "Epocler",
			Indicacoes:        "Proteção hepática, má digestão.",
			Contraindicacoes:  "Hipersensibilidade aos componentes.",
			Posologia:         "1 ampola 2 a 3 vezes ao dia.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "Hypera",
			Registro:          "MS 1.0000.0000.009-5",
		},
		"redoxon": {
			Nome:              "Redoxon",
			Indicacoes:        "Suplementação de vitamina C.",
			Contraindicacoes:  "Hipersensibilidade à vitamina C.",
			Posologia:         "1 comprimido efervescente ao dia.",
			EfeitosColaterais: "Raramente causa efeitos adversos.",
			Laboratorio:       "Bayer",
			Registro:          "MS 1.0000.0000.009-6",
		},
		"naldecon_dia": {
			Nome:              "Naldecon Dia",
			Indicacoes:        "Sintomas de gripe e resfriado.",
			Contraindicacoes:  "Glaucoma, hipertensão grave.",
			Posologia:         "1 cápsula a cada 8 horas.",
			EfeitosColaterais: "Sonolência leve, boca seca.",
			Laboratorio:       "Bristol",
			Registro:          "MS 1.0000.0000.009-7",
		},
		"naldecon_noite": {
			Nome:              "Naldecon Noite",
			Indicacoes:        "Sintomas de gripe e resfriado.",
			Contraindicacoes:  "Glaucoma, hipertensão grave.",
			Posologia:         "1 cápsula antes de dormir.",
			EfeitosColaterais: "Sonolência acentuada, boca seca.",
			Laboratorio:       "Bristol",
			Registro:          "MS 1.0000.0000.009-8",
		},
		"multigrip": {
			Nome:              "Multigrip",
			Indicacoes:        "Sintomas de gripe e resfriado.",
			Contraindicacoes:  "Glaucoma, hipertensão grave.",
			Posologia:         "1 comprimido a cada 6 horas.",
			EfeitosColaterais: "Sonolência, boca seca, tontura.",
			Laboratorio:       "Neo Química",
			Registro:          "MS 1.0000.0000.009-9",
		},
	}

	// Procura o medicamento no mock de forma mais flexível
	for nomeMock, info := range medicamentos {
		// Verifica se o nome do medicamento contém a palavra-chave ou vice-versa
		if strings.Contains(nomeLower, nomeMock) || strings.Contains(nomeMock, nomeLower) {
			return info, nil
		}
	}

	// Se não encontrou, retorna erro
	return nil, fmt.Errorf("bula não encontrada para: %s", nome)
}
