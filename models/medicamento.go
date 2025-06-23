package models

import (
	"database/sql"
	"errors"

	// "io/ioutil" // Não será mais necessário diretamente aqui se saveDB e loadDB forem removidas
	"log"
	"os"
	"strings"
	"time"

	"medicontrol/sqlutils" // Para carregar a query de criação da tabela

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Driver SQLite
)

// Medicamento representa a estrutura de um medicamento
type Medicamento struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Fabricante   string    `json:"fabricante"`
	Tipo         string    `json:"tipo"` // "comprimido", "suspensão", "injetável", etc.
	CodigoANVISA string    `json:"codigo_anvisa"`
	Quantidade   int       `json:"quantidade"`
	Validade     string    `json:"validade"` // Formato: YYYY-MM-DD
	Preco        float64   `json:"preco"`
	CriadoEm     time.Time `json:"criado_em"`
	CategoriaID  string    `json:"categoria_id"`
	Categoria    Categoria `json:"categoria"` // Para incluir dados da categoria aninhados
}

// Movimentacao representa uma entrada ou saída de medicamento
type Movimentacao struct {
	ID            string    `json:"id"`
	MedicamentoID string    `json:"medicamento_id"`
	Tipo          string    `json:"tipo"` // "entrada" ou "saida"
	Quantidade    int       `json:"quantidade"`
	Data          time.Time `json:"data"`
	Observacao    string    `json:"observacao"`
}

var sqlDB *sql.DB // Variável global para a conexão com o banco de dados SQL

// InitDB inicializa o banco de dados SQLite
func InitDB() error {
	log.Println("Iniciando banco de dados SQLite...")

	// Criar diretório data se não existir
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("Erro ao criar diretório data: %v", err)
		return err
	}

	sqliteDbPath := "./data/medicontrol.db"
	var err error
	sqlDB, err = sql.Open("sqlite3", sqliteDbPath)
	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados SQLite: %v", err)
		return err
	}

	// Verificar a conexão
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados SQLite (ping): %v", err)
		return err
	}

	log.Println("Conectado ao banco de dados SQLite com sucesso.")

	// Criar tabela de categorias se não existir
	if err := criarTabelaCategorias(); err != nil {
		// O erro já é logado dentro da função
		return err
	}

	// Criar tabela de medicamentos se não existir
	queryCreateTable := sqlutils.GetQuery("criar_tabela_medicamentos")
	if queryCreateTable == "" {
		log.Fatal("Query 'criar_tabela_medicamentos' não encontrada.")
	} else {
		_, err = sqlDB.Exec(queryCreateTable)
		if err != nil {
			log.Fatalf("Erro ao criar tabela 'medicamentos': %v", err)
			return err
		}
		log.Println("Tabela 'medicamentos' verificada/criada com sucesso.")
	}

	// Aplicar migrações para garantir que o esquema esteja atualizado
	if err := applyMigrations(); err != nil {
		log.Fatalf("Erro ao aplicar migrações de banco de dados: %v", err)
		return err
	}

	// Criar tabela de movimentações se não existir
	if err := criarTabelaMovimentacoes(); err != nil {
		return err
	}

	// Criar tabelas de vendas se não existirem
	if err := criarTabelaVendas(); err != nil {
		return err
	}

	log.Println("Banco de dados SQLite inicializado com sucesso.")
	return nil
}

// GetMedicamentos retorna todos os medicamentos do banco de dados SQLite
func GetMedicamentos() ([]Medicamento, error) {
	query := sqlutils.GetQuery("selecionar_todos_medicamentos")
	if query == "" {
		log.Println("Erro: Query 'selecionar_todos_medicamentos' não encontrada.")
		return nil, errors.New("query para selecionar medicamentos não encontrada")
	}

	rows, err := sqlDB.Query(query)
	if err != nil {
		log.Printf("Erro ao executar query para buscar medicamentos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var medicamentos []Medicamento
	for rows.Next() {
		var med Medicamento
		// Campos da Categoria precisam ser `sql.NullString` para o caso de LEFT JOIN com categoria nula
		var categoriaID, categoriaNome sql.NullString
		var preco sql.NullFloat64

		err := rows.Scan(
			&med.ID, &med.Nome, &med.Fabricante, &med.Tipo,
			&med.CodigoANVISA, &med.Quantidade, &med.Validade, &med.CriadoEm,
			&preco,
			&categoriaID, &categoriaNome,
		)
		if err != nil {
			log.Printf("Erro ao escanear linha do medicamento: %v", err)
			continue
		}

		if preco.Valid {
			med.Preco = preco.Float64
		} else {
			med.Preco = 0.0
		}

		// Preencher dados da categoria
		if categoriaID.Valid {
			med.Categoria.ID = categoriaID.String
			med.CategoriaID = categoriaID.String
		}
		if categoriaNome.Valid {
			med.Categoria.Nome = categoriaNome.String
		}

		medicamentos = append(medicamentos, med)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Erro após iterar sobre as linhas de medicamentos: %v", err)
		return nil, err
	}

	return medicamentos, nil
}

// GetMedicamento retorna um medicamento específico pelo seu ID.
func GetMedicamento(id string) *Medicamento {
	query := sqlutils.GetQuery("selecionar_medicamento_por_id")
	if query == "" {
		log.Printf("Query 'selecionar_medicamento_por_id' não encontrada.")
		return nil
	}

	var med Medicamento
	var categoriaID, categoriaNome sql.NullString
	var preco sql.NullFloat64

	err := sqlDB.QueryRow(query, id).Scan(
		&med.ID, &med.Nome, &med.Fabricante, &med.Tipo,
		&med.CodigoANVISA, &med.Quantidade, &med.Validade, &med.CriadoEm,
		&preco,
		&categoriaID, &categoriaNome,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Erro ao buscar medicamento por ID %s: %v", id, err)
		}
		return nil
	}

	if preco.Valid {
		med.Preco = preco.Float64
	} else {
		med.Preco = 0.0
	}

	// Preencher dados da categoria
	if categoriaID.Valid {
		med.Categoria.ID = categoriaID.String
		med.CategoriaID = categoriaID.String
	}
	if categoriaNome.Valid {
		med.Categoria.Nome = categoriaNome.String
	}

	return &med
}

// GetMedicamentoByCodigoANVISA retorna um medicamento específico pelo código ANVISA
func GetMedicamentoByCodigoANVISA(codigoANVISA string) *Medicamento {
	query := sqlutils.GetQuery("selecionar_medicamento_por_codigo_anvisa")
	if query == "" {
		log.Printf("Query 'selecionar_medicamento_por_codigo_anvisa' não encontrada.")
		return nil
	}

	var med Medicamento
	var categoriaID, categoriaNome sql.NullString
	var preco sql.NullFloat64

	err := sqlDB.QueryRow(query, codigoANVISA).Scan(
		&med.ID, &med.Nome, &med.Fabricante, &med.Tipo,
		&med.CodigoANVISA, &med.Quantidade, &med.Validade, &med.CriadoEm,
		&preco,
		&categoriaID, &categoriaNome,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Erro ao buscar medicamento por código ANVISA %s: %v", codigoANVISA, err)
		}
		// Se for sql.ErrNoRows, simplesmente retorna nil, o que é o comportamento esperado.
		return nil
	}

	if preco.Valid {
		med.Preco = preco.Float64
	} else {
		med.Preco = 0.0
	}

	// Preencher dados da categoria
	if categoriaID.Valid {
		med.Categoria.ID = categoriaID.String
		med.CategoriaID = categoriaID.String
	}
	if categoriaNome.Valid {
		med.Categoria.Nome = categoriaNome.String
	}

	return &med
}

// GetMedicamentoByCodigo é um wrapper para GetMedicamentoByCodigoANVISA
func GetMedicamentoByCodigo(codigo string) *Medicamento {
	return GetMedicamentoByCodigoANVISA(codigo)
}

// BuscarMedicamentos busca medicamentos por nome, fabricante ou código ANVISA.
// Esta é a versão FINAL, robusta e que lida com campos nulos corretamente.
func BuscarMedicamentos(termoBusca string) ([]Medicamento, error) {
	query := `
		SELECT 
			m.ID, m.Nome, m.Fabricante, m.Tipo, m.CodigoANVISA, 
			m.Quantidade, m.Validade, m.Preco, m.CriadoEm, 
			c.ID as CategoriaID, c.Nome as CategoriaNome
		FROM medicamentos m
		LEFT JOIN categorias c ON m.CategoriaID = c.ID`

	var args []interface{}
	if termoBusca != "" {
		query += " WHERE m.Nome LIKE ? OR m.Fabricante LIKE ? OR m.CodigoANVISA LIKE ?"
		searchTerm := "%" + termoBusca + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	rows, err := sqlDB.Query(query, args...)
	if err != nil {
		log.Printf("Erro ao buscar medicamentos com termo '%s': %v", termoBusca, err)
		return nil, err
	}
	defer rows.Close()

	var medicamentos []Medicamento
	for rows.Next() {
		var med Medicamento
		// Usar tipos sql.Null* para todos os campos que podem ser nulos
		var fabricante, tipo, validade, categoriaID, categoriaNome sql.NullString
		var preco sql.NullFloat64

		err := rows.Scan(
			&med.ID,
			&med.Nome,
			&fabricante,
			&tipo,
			&med.CodigoANVISA,
			&med.Quantidade,
			&validade,
			&preco,
			&med.CriadoEm,
			&categoriaID,
			&categoriaNome,
		)
		if err != nil {
			log.Printf("Erro ao escanear linha da busca por medicamento: %v", err)
			continue // Pula para a próxima linha em caso de erro
		}

		// Atribuir valores somente se eles forem válidos (não nulos)
		med.Fabricante = fabricante.String
		med.Tipo = tipo.String
		med.Validade = validade.String
		if preco.Valid {
			med.Preco = preco.Float64
		} else {
			med.Preco = 0.0 // Valor padrão se o preço for nulo
		}
		if categoriaID.Valid {
			med.Categoria.ID = categoriaID.String
			med.CategoriaID = categoriaID.String
		}
		if categoriaNome.Valid {
			med.Categoria.Nome = categoriaNome.String
		}

		medicamentos = append(medicamentos, med)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Erro após iterar sobre os resultados da busca por medicamentos: %v", err)
		return nil, err
	}

	return medicamentos, nil
}

// AddMedicamento adiciona um novo medicamento ao banco de dados SQLite
func AddMedicamento(med *Medicamento) error {
	// Garante que o ID seja gerado se estiver vazio
	if med.ID == "" {
		med.ID = uuid.New().String()
	}
	med.CriadoEm = time.Now()

	query := sqlutils.GetQuery("inserir_medicamento")
	if query == "" {
		return errors.New("query 'inserir_medicamento' não encontrada")
	}

	_, err := sqlDB.Exec(query,
		med.ID,
		med.Nome,
		med.Fabricante,
		med.Tipo,
		med.CodigoANVISA,
		med.Quantidade,
		med.Validade,
		med.Preco,
		med.CriadoEm,
		med.CategoriaID,
	)

	if err != nil {
		log.Printf("Erro ao inserir medicamento no banco de dados: %v", err)
		return err
	}
	return nil
}

// UpdateMedicamento atualiza um medicamento existente no banco de dados SQLite
func UpdateMedicamento(med *Medicamento) error {
	query := sqlutils.GetQuery("atualizar_medicamento")
	if query == "" {
		return errors.New("query 'atualizar_medicamento' não encontrada")
	}

	_, err := sqlDB.Exec(query,
		med.Nome,
		med.Fabricante,
		med.Tipo,
		med.CodigoANVISA,
		med.Quantidade,
		med.Validade,
		med.Preco,
		med.CategoriaID,
		med.ID,
	)

	if err != nil {
		log.Printf("Erro ao atualizar medicamento no banco de dados: %v", err)
		return err
	}
	return nil
}

// DeleteMedicamento remove um medicamento do banco de dados SQLite
func DeleteMedicamento(id string) error {
	query := sqlutils.GetQuery("deletar_medicamento")
	if query == "" {
		return errors.New("query 'deletar_medicamento' não encontrada")
	}

	_, err := sqlDB.Exec(query, id)

	if err != nil {
		log.Printf("Erro ao deletar medicamento do banco de dados: %v", err)
		return err
	}
	return nil
}

// RegistrarMovimentacao registra uma entrada ou saída de medicamento e atualiza o estoque
func RegistrarMovimentacao(mov Movimentacao) error {
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}

	// Primeiro, busca o medicamento para verificar o estoque
	medicamento := GetMedicamento(mov.MedicamentoID)
	if medicamento == nil {
		tx.Rollback()
		return errors.New("medicamento não encontrado para a movimentação")
	}

	// Calcula a nova quantidade
	novaQuantidade := medicamento.Quantidade
	if mov.Tipo == "entrada" {
		novaQuantidade += mov.Quantidade
	} else if mov.Tipo == "saida" {
		if novaQuantidade < mov.Quantidade {
			tx.Rollback()
			return errors.New("quantidade em estoque insuficiente para a saída")
		}
		novaQuantidade -= mov.Quantidade
	}

	// Atualiza a quantidade do medicamento
	updateQuery := "UPDATE medicamentos SET Quantidade = ? WHERE ID = ?"
	_, err = tx.Exec(updateQuery, novaQuantidade, mov.MedicamentoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insere o registro da movimentação
	mov.ID = uuid.New().String()
	mov.Data = time.Now()
	insertQuery := sqlutils.GetQuery("inserir_movimentacao")
	if insertQuery == "" {
		tx.Rollback()
		return errors.New("query 'inserir_movimentacao' não encontrada")
	}

	_, err = tx.Exec(insertQuery, mov.ID, mov.MedicamentoID, mov.Tipo, mov.Quantidade, mov.Data, mov.Observacao)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetMovimentacoes retorna todas as movimentações com detalhes do medicamento
func GetMovimentacoes() ([]map[string]interface{}, error) {
	query := sqlutils.GetQuery("selecionar_todas_movimentacoes")
	if query == "" {
		return nil, errors.New("query 'selecionar_todas_movimentacoes' não encontrada")
	}

	rows, err := sqlDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movimentacoes []map[string]interface{}
	for rows.Next() {
		var mov Movimentacao
		var nomeMedicamento, tipoMedicamento string
		err := rows.Scan(&mov.ID, &mov.MedicamentoID, &mov.Tipo, &mov.Quantidade, &mov.Data, &mov.Observacao, &nomeMedicamento, &tipoMedicamento)
		if err != nil {
			log.Printf("Erro ao escanear movimentação: %v", err)
			continue
		}
		movimentacoes = append(movimentacoes, map[string]interface{}{
			"id":               mov.ID,
			"medicamento_id":   mov.MedicamentoID,
			"nome_medicamento": nomeMedicamento,
			"tipo_medicamento": tipoMedicamento,
			"tipo":             mov.Tipo,
			"quantidade":       mov.Quantidade,
			"data":             mov.Data,
			"observacao":       mov.Observacao,
		})
	}
	return movimentacoes, nil
}

// GetTotalVendas retorna a soma de todas as quantidades de itens de venda.
func GetTotalVendas() (int, error) {
	query := sqlutils.GetQuery("contar_total_vendas")
	if query == "" {
		return 0, errors.New("query 'contar_total_vendas' não encontrada")
	}

	var totalVendas sql.NullInt64 // Usar NullInt64 para o caso de não haver vendas
	err := sqlDB.QueryRow(query).Scan(&totalVendas)
	if err != nil {
		// Se não houver linhas, o que significa nenhuma venda, o total é 0.
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	if totalVendas.Valid {
		return int(totalVendas.Int64), nil
	}
	return 0, nil
}

// criarTabelaMovimentacoes cria a tabela 'movimentacoes' se ela não existir.
func criarTabelaMovimentacoes() error {
	query := sqlutils.GetQuery("criar_tabela_movimentacoes")
	if query == "" {
		return errors.New("query 'criar_tabela_movimentacoes' não encontrada")
	}
	_, err := sqlDB.Exec(query)
	if err != nil {
		log.Printf("Erro ao criar tabela 'movimentacoes': %v", err)
	}
	return err
}

// criarTabelaVendas cria as tabelas 'vendas' e 'venda_items' se não existirem.
func criarTabelaVendas() error {
	queryVendas := sqlutils.GetQuery("criar_tabela_vendas")
	if queryVendas == "" {
		return errors.New("query 'criar_tabela_vendas' não encontrada")
	}
	if _, err := sqlDB.Exec(queryVendas); err != nil {
		log.Printf("Erro ao criar tabela 'vendas': %v", err)
		return err
	}

	queryVendaItems := sqlutils.GetQuery("criar_tabela_venda_items")
	if queryVendaItems == "" {
		return errors.New("query 'criar_tabela_venda_items' não encontrada")
	}
	if _, err := sqlDB.Exec(queryVendaItems); err != nil {
		log.Printf("Erro ao criar tabela 'venda_items': %v", err)
		return err
	}

	return nil
}

// applyMigrations aplica migrações no banco de dados, como adicionar novas colunas.
func applyMigrations() error {
	// Adiciona a coluna 'preco' se ela não existir
	if err := addColumnIfNotExists("medicamentos", "Preco", "REAL DEFAULT 0.0"); err != nil {
		return err
	}
	// Adiciona a coluna 'categoria_id' se ela não existir
	if err := addColumnIfNotExists("medicamentos", "CategoriaID", "TEXT"); err != nil {
		return err
	}
	return nil
}

// addColumnIfNotExists verifica se uma coluna existe e a adiciona se não existir.
func addColumnIfNotExists(tableName, columnName, columnType string) error {
	// Query para verificar se a coluna existe
	// Note que PRAGMA_table_info é específico do SQLite
	rows, err := sqlDB.Query("PRAGMA table_info(" + tableName + ")")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var type_ string
		var notnull bool
		var dflt_value interface{}
		var pk int
		if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
			return err
		}
		if strings.EqualFold(name, columnName) {
			return nil // Coluna já existe
		}
	}

	// Coluna não encontrada, então a adiciona
	_, err = sqlDB.Exec("ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + columnType)
	if err != nil {
		// Log para ajudar a depurar erros de sintaxe SQL, etc.
		log.Printf("Erro ao adicionar coluna '%s' à tabela '%s': %v", columnName, tableName, err)
	} else {
		log.Printf("Coluna '%s' adicionada à tabela '%s' com sucesso.", columnName, tableName)
	}
	return err
}

// GetMedicamentosBaixoEstoque retorna medicamentos com quantidade abaixo de um limite.
func GetMedicamentosBaixoEstoque(limite int) ([]Medicamento, error) {
	query := sqlutils.GetQuery("selecionar_medicamentos_baixo_estoque")
	if query == "" {
		return nil, errors.New("query 'selecionar_medicamentos_baixo_estoque' não encontrada")
	}
	rows, err := sqlDB.Query(query, limite)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medicamentos []Medicamento
	for rows.Next() {
		var med Medicamento
		err := rows.Scan(&med.ID, &med.Nome, &med.Fabricante, &med.Quantidade)
		if err != nil {
			log.Printf("Erro ao escanear medicamento com baixo estoque: %v", err)
			continue
		}
		medicamentos = append(medicamentos, med)
	}
	return medicamentos, nil
}
