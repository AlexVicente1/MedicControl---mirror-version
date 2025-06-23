package models

import (
	"database/sql"
	"errors"
	"log"

	"medicontrol/sqlutils"

	"github.com/google/uuid"
)

// Categoria representa uma categoria de medicamento
type Categoria struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

// criarTabelaCategorias cria a tabela de categorias no banco de dados se ela não existir.
func criarTabelaCategorias() error {
	query := sqlutils.GetQuery("criar_tabela_categorias")
	if query == "" {
		log.Fatal("Query 'criar_tabela_categorias' não encontrada.")
	}

	_, err := sqlDB.Exec(query)
	if err != nil {
		log.Printf("Erro ao criar tabela 'categorias': %v", err)
		return err
	}
	log.Println("Tabela 'categorias' verificada/criada com sucesso.")
	return nil
}

// AddCategoria adiciona uma nova categoria e retorna seu ID.
func AddCategoria(nome string) (string, error) {
	// Verificar se a categoria já existe
	categoriaExistente, err := GetCategoriaByNome(nome)
	if err != nil {
		// Se o erro não for 'não encontrado', temos um problema real
		if err != sql.ErrNoRows {
			log.Printf("Erro ao verificar existência da categoria '%s': %v", nome, err)
			return "", err
		}
	}
	if categoriaExistente != nil {
		log.Printf("Categoria '%s' já existe com ID: %s. Não será adicionada novamente.", nome, categoriaExistente.ID)
		return categoriaExistente.ID, nil // Retorna o ID da categoria existente
	}

	categoria := Categoria{
		ID:   uuid.New().String(),
		Nome: nome,
	}

	query := sqlutils.GetQuery("inserir_categoria")
	if query == "" {
		return "", errors.New("query 'inserir_categoria' não encontrada")
	}

	_, err = sqlDB.Exec(query, categoria.ID, categoria.Nome)
	if err != nil {
		log.Printf("Erro ao inserir categoria '%s': %v", nome, err)
		return "", err
	}

	log.Printf("Categoria '%s' adicionada com sucesso com ID: %s", nome, categoria.ID)
	return categoria.ID, nil
}

// GetCategoriaByNome busca uma categoria pelo nome.
func GetCategoriaByNome(nome string) (*Categoria, error) {
	query := sqlutils.GetQuery("selecionar_categoria_por_nome")
	if query == "" {
		return nil, errors.New("query 'selecionar_categoria_por_nome' não encontrada")
	}

	var cat Categoria
	err := sqlDB.QueryRow(query, nome).Scan(&cat.ID, &cat.Nome)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Não é um erro, apenas não encontrou
		}
		return nil, err
	}
	return &cat, nil
}

// GetAllCategorias retorna todas as categorias do banco de dados.
func GetAllCategorias() ([]Categoria, error) {
	query := sqlutils.GetQuery("selecionar_todas_categorias")
	if query == "" {
		return nil, errors.New("query 'selecionar_todas_categorias' não encontrada")
	}

	rows, err := sqlDB.Query(query)
	if err != nil {
		log.Printf("Erro ao buscar categorias: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categorias []Categoria
	for rows.Next() {
		var cat Categoria
		if err := rows.Scan(&cat.ID, &cat.Nome); err != nil {
			log.Printf("Erro ao escanear linha da categoria: %v", err)
			continue
		}
		categorias = append(categorias, cat)
	}

	return categorias, nil
}
