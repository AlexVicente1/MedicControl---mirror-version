package db

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Database representa a estrutura do banco de dados
type Database struct {
	Produtos []Produto `json:"produtos"`
	NextID   int       `json:"nextId"`
	mu       sync.RWMutex
	filePath string
}

// Produto representa um produto no banco de dados
type Produto struct {
	ID           int     `json:"id"`
	Nome         string  `json:"nome"`
	Lote         string  `json:"lote"`
	Vencimento   string  `json:"vencimento"`
	Laboratorio  string  `json:"laboratorio"`
	Preco        float64 `json:"preco"`
	Quantidade   int     `json:"quantidade"`
	Categoria    string  `json:"categoria"`
	Subcategoria string  `json:"subcategoria"`
	DataCadastro string  `json:"dataCadastro"`
}

var instance *Database
var once sync.Once

// GetInstance retorna a instância única do banco de dados
func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{
			Produtos: make([]Produto, 0),
			NextID:   1,
			filePath: "data/produtos.json",
		}
		instance.init()
	})
	return instance
}

// init inicializa o banco de dados
func (db *Database) init() {
	// Criar diretório data se não existir
	dir := filepath.Dir(db.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	// Tentar carregar dados existentes
	if err := db.load(); err != nil {
		// Se o arquivo não existe, criar um novo
		if os.IsNotExist(err) {
			// Tentar importar do produtos_completo.json
			if err := db.ImportProdutosCompletos(); err != nil {
				// Se não conseguir importar, criar estrutura vazia
				db.Produtos = []Produto{}
				db.NextID = 1
				db.save()
			}
		} else {
			panic(err)
		}
	}
}

// load carrega os dados do arquivo
func (db *Database) load() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := os.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, db)
}

// save salva os dados no arquivo
func (db *Database) save() error {
	data, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.filePath, data, 0644)
}

// AddProduto adiciona um novo produto
func (db *Database) AddProduto(p Produto) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	p.ID = db.NextID
	db.NextID++
	db.Produtos = append(db.Produtos, p)
	return db.save()
}

// GetProdutos retorna todos os produtos
func (db *Database) GetProdutos() []Produto {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Log para debug
	log.Printf("Buscando produtos do banco de dados. Total: %d", len(db.Produtos))

	produtos := make([]Produto, len(db.Produtos))
	copy(produtos, db.Produtos)

	// Normalizar as categorias para o formato do frontend
	for i := range produtos {
		switch produtos[i].Categoria {
		case "Saúde e Bem-estar":
			produtos[i].Categoria = "saude-bem-estar"
		case "Mãe e Bebê":
			produtos[i].Categoria = "mae-bebe"
		case "Medicamentos Controlados":
			produtos[i].Categoria = "medicamentos-controlados"
		case "Medicamentos Comuns":
			produtos[i].Categoria = "medicamentos-comuns"
		case "Higiene Pessoal":
			produtos[i].Categoria = "higiene-pessoal"
		default:
			// Converter para lowercase e substituir espaços por hífens
			produtos[i].Categoria = strings.ToLower(strings.ReplaceAll(produtos[i].Categoria, " ", "-"))
		}

		// Log para debug de cada produto
		log.Printf("Produto: %s, Categoria: %s, Subcategoria: %s",
			produtos[i].Nome, produtos[i].Categoria, produtos[i].Subcategoria)
	}

	return produtos
}

// DeleteProduto remove um produto pelo ID
func (db *Database) DeleteProduto(id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, p := range db.Produtos {
		if p.ID == id {
			db.Produtos = append(db.Produtos[:i], db.Produtos[i+1:]...)
			break
		}
	}
	return db.save()
}

// SearchProdutos busca produtos com filtros
func (db *Database) SearchProdutos(termo, categoria, subcategoria string) []Produto {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var resultados []Produto
	termoLower := strings.ToLower(termo)
	categoriaLower := strings.ToLower(categoria)
	subcategoriaLower := strings.ToLower(subcategoria)

	for _, p := range db.Produtos {
		if (termo == "" ||
			strings.Contains(strings.ToLower(p.Nome), termoLower) ||
			strings.Contains(strings.ToLower(p.Lote), termoLower) ||
			strings.Contains(strings.ToLower(p.Laboratorio), termoLower)) &&
			(categoria == "" || strings.ToLower(p.Categoria) == categoriaLower) &&
			(subcategoria == "" || strings.ToLower(p.Subcategoria) == subcategoriaLower) {
			resultados = append(resultados, p)
		}
	}

	return resultados
}

// ImportProdutosCompletos importa produtos do arquivo produtos_completo.json
func (db *Database) ImportProdutosCompletos() error {
	// Ler o arquivo produtos_completo.json
	data, err := os.ReadFile("data/produtos_completo.json")
	if err != nil {
		return err
	}

	// Decodificar o JSON
	var produtosImportados []struct {
		ID             int     `json:"id"`
		Nome           string  `json:"nome"`
		Categoria      string  `json:"categoria"`
		Subcategoria   string  `json:"subcategoria"`
		Preco          float64 `json:"preco"`
		Estoque        int     `json:"estoque"`
		Descricao      string  `json:"descricao"`
		DataCadastro   string  `json:"dataCadastro"`
		DataValidade   string  `json:"dataValidade"`
		Laboratorio    string  `json:"laboratorio"`
		Controlado     bool    `json:"controlado"`
		PrincipioAtivo string  `json:"principioAtivo"`
	}

	if err := json.Unmarshal(data, &produtosImportados); err != nil {
		return err
	}

	// Converter para o formato do banco de dados
	db.mu.Lock()
	defer db.mu.Unlock()

	db.Produtos = []Produto{} // Limpar produtos existentes
	db.NextID = 1

	// Converter cada produto
	for _, p := range produtosImportados {
		produto := Produto{
			ID:           p.ID,
			Nome:         p.Nome,
			Categoria:    p.Categoria,
			Subcategoria: p.Subcategoria,
			Preco:        p.Preco,
			Quantidade:   p.Estoque,
			DataCadastro: p.DataCadastro,
			Vencimento:   p.DataValidade,
			Laboratorio:  p.Laboratorio,
			Lote:         "LOTE-" + time.Now().Format("20060102"),
		}
		db.Produtos = append(db.Produtos, produto)
		if p.ID >= db.NextID {
			db.NextID = p.ID + 1
		}
	}

	return db.save()
}
