package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Produto representa um produto no sistema
type Produto struct {
	ID           int     `json:"id"`
	Nome         string  `json:"nome"`
	Categoria    string  `json:"categoria"`
	Subcategoria string  `json:"subcategoria"`
	Preco        float64 `json:"preco"`
	Quantidade   int     `json:"quantidade"`
	Laboratorio  string  `json:"laboratorio"`
	Lote         string  `json:"lote"`
	Vencimento   string  `json:"vencimento"`
}

// BulaInfo representa as informações da bula de um medicamento
type BulaInfo struct {
	Nome              string `json:"nome"`
	Laboratorio       string `json:"laboratorio"`
	Registro          string `json:"registro"`
	Principio         string `json:"principio"`
	ClasseTerapeutica string `json:"classe_terapeutica"`
	Indicacoes        string `json:"indicacoes"`
	Contraindicacoes  string `json:"contraindicacoes"`
	Posologia         string `json:"posologia"`
	EfeitosColaterais string `json:"efeitos_colaterais"`
}

// Banco de dados em memória
var produtos []Produto

func main() {
	// Carregar produtos do arquivo JSON
	if err := carregarProdutos(); err != nil {
		log.Fatal("Erro ao carregar produtos:", err)
	}

	r := mux.NewRouter()

	// Servir arquivos estáticos
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Rotas da API
	r.HandleFunc("/api/produtos", getProdutosHandler).Methods("GET")
	r.HandleFunc("/api/produtos", addProdutoHandler).Methods("POST")
	r.HandleFunc("/api/produtos/{id}", deleteProdutoHandler).Methods("DELETE")
	r.HandleFunc("/api/bula/{nome}", getBulaHandler).Methods("GET")

	// Rota principal
	r.HandleFunc("/", indexHandler)

	// Iniciar servidor
	log.Println("Servidor iniciando na porta 8080...")
	go abrirNavegador("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// API Handlers
func getProdutosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"produtos": produtos,
	})
}

func addProdutoHandler(w http.ResponseWriter, r *http.Request) {
	var produto Produto
	if err := json.NewDecoder(r.Body).Decode(&produto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Gerar novo ID
	maxID := 0
	for _, p := range produtos {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	produto.ID = maxID + 1

	produtos = append(produtos, produto)
	salvarProdutos()

	w.WriteHeader(http.StatusCreated)
}

func deleteProdutoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := 0
	fmt.Sscanf(vars["id"], "%d", &id)

	for i, p := range produtos {
		if p.ID == id {
			produtos = append(produtos[:i], produtos[i+1:]...)
			salvarProdutos()
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func carregarProdutos() error {
	// Criar diretório data se não existir
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	// Tentar ler o arquivo de produtos
	data, err := os.ReadFile(filepath.Join("data", "produtos.json"))
	if err != nil {
		if os.IsNotExist(err) {
			// Se o arquivo não existe, criar um array vazio
			produtos = []Produto{}
			return nil
		}
		return err
	}

	var produtosData struct {
		Produtos []Produto `json:"produtos"`
	}
	if err := json.Unmarshal(data, &produtosData); err != nil {
		return err
	}
	produtos = produtosData.Produtos
	return nil
}

func salvarProdutos() error {
	data, err := json.MarshalIndent(map[string]interface{}{
		"produtos": produtos,
	}, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("data", "produtos.json"), data, 0644)
}

// Função para abrir o navegador automaticamente
func abrirNavegador(url string) {
	time.Sleep(1 * time.Second)
	var comando string
	var args []string

	switch runtime.GOOS {
	case "windows":
		comando = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "linux":
		comando = "xdg-open"
		args = []string{url}
	case "darwin":
		comando = "open"
		args = []string{url}
	}

	exec.Command(comando, args...).Start()
}

// Função para buscar bula da base local
func buscarBulaLocal(nome string) (*BulaInfo, error) {
	// Ler o arquivo de bulas
	data, err := os.ReadFile("data/bulas.json")
	if err != nil {
		log.Printf("Erro ao ler arquivo de bulas: %v", err)
		return nil, err
	}

	// Parse do JSON
	var bulasDB struct {
		Bulas map[string]BulaInfo `json:"bulas"`
	}
	if err := json.Unmarshal(data, &bulasDB); err != nil {
		log.Printf("Erro ao fazer parse do JSON de bulas: %v", err)
		return nil, err
	}

	// Converter o nome para minúsculas para comparação
	nomeLower := strings.ToLower(nome)

	// Lista de palavras-chave para busca
	keywords := []string{
		"rivotril", "clonazepam",
		"dipirona", "novalgina",
		"paracetamol", "tylenol",
		"omeprazol",
		"sertralina", "zoloft",
		"protetor", "solar", "fps",
		"vitamina", "suplemento",
		"sabonete", "higiene",
		"whey", "protein", "proteina",
		"shampoo", "xampu",
	}

	// Primeiro tenta encontrar correspondência exata
	for key, bula := range bulasDB.Bulas {
		if strings.Contains(nomeLower, key) {
			return &bula, nil
		}
	}

	// Se não encontrar correspondência exata, tenta por palavras-chave
	for _, keyword := range keywords {
		if strings.Contains(nomeLower, keyword) {
			for key, bula := range bulasDB.Bulas {
				if strings.Contains(key, keyword) {
					return &bula, nil
				}
			}
		}
	}

	// Se não encontrar nenhuma correspondência, procura por categoria
	categorias := map[string]string{
		"protetor": "protetor solar",
		"vitamina": "vitamina c",
		"sabonete": "sabonete",
		"shampoo":  "shampoo",
		"whey":     "whey protein",
	}

	for categoria, key := range categorias {
		if strings.Contains(nomeLower, categoria) {
			if bula, ok := bulasDB.Bulas[key]; ok {
				return &bula, nil
			}
		}
	}

	// Se ainda não encontrou, retorna o mock
	return buscarBulaMock(nome), nil
}

// Handler para buscar bula
func getBulaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nomeMedicamento := vars["nome"]

	// Primeiro tenta buscar da base local
	bula, err := buscarBulaLocal(nomeMedicamento)
	if err != nil {
		// Se der erro na base local, tenta usar o mock
		bula = buscarBulaMock(nomeMedicamento)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bula)
}

// Função para buscar bula mock (quando não encontrar na ANVISA)
func buscarBulaMock(nome string) *BulaInfo {
	nomeLower := strings.ToLower(nome)

	// Mapa de bulas mockadas para testes
	bulas := map[string]*BulaInfo{
		"dipirona": {
			Nome:              "Dipirona Monoidratada",
			Laboratorio:       "Medley",
			Registro:          "1.0047.0118",
			Principio:         "Dipirona Monoidratada",
			ClasseTerapeutica: "Analgésico e Antitérmico",
			Indicacoes:        "Indicado como analgésico e antitérmico. Utilizado no tratamento de dores de diversos tipos como: dor de cabeça, dor de dente, dor muscular, cólicas menstruais e febre.",
			Contraindicacoes:  "Não deve ser utilizado em pacientes com alergia à dipirona ou a qualquer componente da fórmula; problemas de sangue ou que fazem tratamento com medicamentos que podem causar alterações no sangue; deficiência de uma enzima chamada G6PD; porfiria aguda; grávidas no primeiro e último trimestre de gestação.",
			Posologia:         "Adultos e adolescentes acima de 15 anos: 1 comprimido de 500mg a 1g, 3 a 4 vezes ao dia. Dose máxima diária de 4g. Crianças e adolescentes de 5 a 14 anos: 1 comprimido de 500mg, 3 a 4 vezes ao dia.",
			EfeitosColaterais: "Podem ocorrer reações alérgicas, desde leves (vermelhidão na pele) até graves; problemas no sangue (agranulocitose); queda da pressão arterial; náuseas; vômitos; irritação no estômago.",
		},
		"rivotril": {
			Nome:              "Rivotril (Clonazepam)",
			Laboratorio:       "Roche",
			Registro:          "1.0100.0075",
			Principio:         "Clonazepam",
			ClasseTerapeutica: "Anticonvulsivante e Ansiolítico",
			Indicacoes:        "Tratamento de crises epilépticas e convulsivas. Transtornos de ansiedade como síndrome do pânico, fobia social, transtornos de humor. Tratamento de espasmos e movimentos involuntários.",
			Contraindicacoes:  "Pacientes com hipersensibilidade ao clonazepam ou a outros benzodiazepínicos. Pacientes com insuficiência respiratória grave ou apneia do sono. Pacientes com insuficiência hepática grave. Glaucoma agudo de ângulo fechado.",
			Posologia:         "A dose deve ser individualizada. Iniciar com doses baixas: 0,5mg, 2 a 3 vezes ao dia, aumentando gradualmente conforme necessidade e tolerância. Dose máxima diária recomendada: 4mg para adultos. O medicamento deve ser usado sob estrita supervisão médica.",
			EfeitosColaterais: "Sonolência, fadiga, diminuição da coordenação motora, alterações da memória, tontura, depressão respiratória, dependência física e psíquica com uso prolongado. Pode prejudicar a capacidade de dirigir veículos e operar máquinas.",
		},
		"paracetamol": {
			Nome:              "Paracetamol",
			Laboratorio:       "EMS",
			Registro:          "1.0235.0264",
			Principio:         "Paracetamol (Acetaminofeno)",
			ClasseTerapeutica: "Analgésico e Antitérmico",
			Indicacoes:        "Alívio temporário de dores leves a moderadas como dores de cabeça, musculares, articulares, nas costas, dores de dente, cólicas menstruais e redução da febre.",
			Contraindicacoes:  "Pacientes com hipersensibilidade ao paracetamol ou a qualquer componente da fórmula. Doença grave do fígado. Não deve ser utilizado por período prolongado ou em doses maiores que as recomendadas.",
			Posologia:         "Adultos e crianças acima de 12 anos: 1 comprimido de 500mg a 750mg a cada 4 a 6 horas. Não exceder 4g (8 comprimidos de 500mg) em 24 horas. Crianças: consultar o médico para dosagem adequada.",
			EfeitosColaterais: "Quando usado nas doses recomendadas, os efeitos colaterais são raros. Podem ocorrer reações alérgicas leves. Em doses muito altas pode causar lesão no fígado. Uso com álcool pode aumentar o risco de danos ao fígado.",
		},
		"omeprazol": {
			Nome:              "Omeprazol",
			Laboratorio:       "EMS",
			Registro:          "1.0235.0446",
			Principio:         "Omeprazol",
			ClasseTerapeutica: "Inibidor da Bomba de Prótons",
			Indicacoes:        "Tratamento de úlcera gástrica e duodenal, esofagite de refluxo, síndrome de Zollinger-Ellison, prevenção de úlceras causadas por medicamentos anti-inflamatórios não esteroidais (AINEs).",
			Contraindicacoes:  "Hipersensibilidade conhecida ao omeprazol ou a qualquer componente da fórmula. Pacientes que estejam tomando medicamentos contendo nelfinavir. Gravidez e amamentação devem ser avaliadas pelo médico.",
			Posologia:         "Úlcera duodenal: 20mg uma vez ao dia por 4 semanas. Úlcera gástrica: 20mg uma vez ao dia por 8 semanas. Refluxo gastroesofágico: 20mg uma vez ao dia por 4 a 8 semanas. Dose de manutenção: 10mg a 20mg uma vez ao dia.",
			EfeitosColaterais: "Dor de cabeça, diarreia, constipação, dor abdominal, náusea, vômito, gases, boca seca, tontura. Raramente: reações alérgicas, alterações nas células do sangue, problemas no fígado, alterações nos níveis de magnésio no sangue.",
		},
		"sertralina": {
			Nome:              "Sertralina",
			Laboratorio:       "EMS",
			Registro:          "1.0235.0862",
			Principio:         "Cloridrato de Sertralina",
			ClasseTerapeutica: "Antidepressivo (ISRS)",
			Indicacoes:        "Tratamento da depressão, transtorno obsessivo-compulsivo (TOC), transtorno do pânico, transtorno de estresse pós-traumático, fobia social e síndrome pré-menstrual.",
			Contraindicacoes:  "Uso concomitante com medicamentos IMAO (inibidores da monoaminoxidase). Pacientes com hipersensibilidade à sertralina. Uso durante a gravidez e amamentação deve ser avaliado pelo médico.",
			Posologia:         "Depressão e TOC: Iniciar com 50mg uma vez ao dia, podendo ser aumentada gradualmente até máximo de 200mg por dia. Pânico, estresse pós-traumático e fobia social: Iniciar com 25mg uma vez ao dia na primeira semana, depois 50mg uma vez ao dia.",
			EfeitosColaterais: "Náusea, diarreia, tremor, insônia, sonolência, boca seca, diminuição do apetite, sudorese, disfunção sexual. Pode ocorrer agitação e ansiedade nas primeiras semanas de tratamento.",
		},
	}

	// Procura por correspondência parcial no nome
	for bulaName, bula := range bulas {
		if strings.Contains(nomeLower, bulaName) {
			return bula
		}
	}

	// Se não encontrar correspondência exata, tenta encontrar por categoria
	categoriaBulas := map[string]*BulaInfo{
		"vitamina": {
			Nome:              nome,
			Laboratorio:       "Diversos",
			Registro:          "Consulte a embalagem",
			Principio:         "Complexo Vitamínico",
			ClasseTerapeutica: "Suplemento Vitamínico",
			Indicacoes:        "Suplementação vitamínica em casos de deficiência, prevenção de carências nutricionais, auxílio na manutenção da saúde.",
			Contraindicacoes:  "Hipersensibilidade aos componentes da fórmula. Algumas vitaminas podem interagir com medicamentos, consulte seu médico.",
			Posologia:         "Geralmente 1 comprimido ao dia ou conforme orientação médica. A dose pode variar de acordo com a concentração e composição do produto.",
			EfeitosColaterais: "Geralmente bem tolerado quando usado nas doses recomendadas. Pode ocorrer desconforto gastrointestinal, dor de cabeça, alteração na cor da urina.",
		},
		"protetor solar": {
			Nome:              nome,
			Laboratorio:       "Diversos",
			Registro:          "Consulte a embalagem",
			Principio:         "Filtros Solares UVA e UVB",
			ClasseTerapeutica: "Protetor Solar",
			Indicacoes:        "Proteção da pele contra os raios solares UVA e UVB, prevenção do envelhecimento precoce e do câncer de pele.",
			Contraindicacoes:  "Hipersensibilidade aos componentes da fórmula. Em caso de irritação, descontinuar o uso.",
			Posologia:         "Aplicar generosamente sobre a pele limpa 15 minutos antes da exposição ao sol. Reaplicar a cada 2 horas ou após mergulho ou suor excessivo.",
			EfeitosColaterais: "Pode causar irritação, vermelhidão ou coceira em peles sensíveis. Em caso de contato com os olhos, lavar com água em abundância.",
		},
	}

	// Tenta encontrar por categoria
	for categoria, bula := range categoriaBulas {
		if strings.Contains(nomeLower, categoria) {
			return bula
		}
	}

	// Retorna uma bula genérica se não encontrar
	return &BulaInfo{
		Nome:              nome,
		Laboratorio:       "Não encontrado",
		Registro:          "N/A",
		Principio:         "Não encontrado",
		ClasseTerapeutica: "Não encontrado",
		Indicacoes:        "Por favor, consulte a bula física do medicamento ou procure orientação do farmacêutico.",
		Contraindicacoes:  "Por favor, consulte a bula física do medicamento ou procure orientação do farmacêutico.",
		Posologia:         "Por favor, consulte a bula física do medicamento ou procure orientação do farmacêutico.",
		EfeitosColaterais: "Por favor, consulte a bula física do medicamento ou procure orientação do farmacêutico.",
	}
}
