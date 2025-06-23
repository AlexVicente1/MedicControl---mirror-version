package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"medicontrol/sqlutils"
	"time"
)

// Venda representa a tabela 'vendas'
type Venda struct {
	ID     int       `json:"id"`
	Data   time.Time `json:"data"`
	UserID int       `json:"user_id"`
}

// VendaItem representa a tabela 'venda_items'
type VendaItem struct {
	ID            int     `json:"id"`
	VendaID       int     `json:"venda_id"`
	MedicamentoID int     `json:"medicamento_id"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
}

// RegistrarVendaRequest é o que a API recebe para criar uma venda
type RegistrarVendaRequest struct {
	Itens []struct {
		MedicamentoID int `json:"medicamento_id"`
		Quantidade    int `json:"quantidade"`
	} `json:"itens"`
}

// VendaInfo é a struct para os dados de resumo da lista de vendas
type VendaInfo struct {
	ID              int       `json:"id"`
	Data            time.Time `json:"data"`
	UserID          int       `json:"user_id"`
	QuantidadeItens int       `json:"quantidade_itens"`
	TotalVenda      float64   `json:"total_venda"`
}

// ListarVendas busca um resumo de todas as vendas no banco de dados.
func ListarVendas() ([]VendaInfo, error) {
	query := sqlutils.GetQuery("ListarVendas")
	if query == "" {
		return nil, errors.New("query ListarVendas não encontrada")
	}

	rows, err := sqlDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar a query para listar vendas: %w", err)
	}
	defer rows.Close()

	var vendas []VendaInfo
	for rows.Next() {
		var v VendaInfo
		// Corresponde à ordem SELECT da query ListarVendas
		if err := rows.Scan(&v.ID, &v.Data, &v.UserID, &v.QuantidadeItens, &v.TotalVenda); err != nil {
			return nil, fmt.Errorf("erro ao scanear linha da venda: %w", err)
		}
		vendas = append(vendas, v)
	}

	return vendas, nil
}

// RegistrarVenda processa uma nova venda, atualizando o estoque e registrando os itens.
func RegistrarVenda(req RegistrarVendaRequest) (int64, error) {
	tx, err := sqlDB.Begin()
	if err != nil {
		return 0, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback() // Rollback é uma proteção; só tem efeito se Commit não for chamado.

	// 1. Inserir na tabela 'vendas' para gerar um ID de venda.
	queryInsertVenda := sqlutils.GetQuery("InserirVenda")
	if queryInsertVenda == "" {
		return 0, errors.New("query InserirVenda não encontrada")
	}
	res, err := tx.Exec(queryInsertVenda, 1) // Provisório: UserID fixo como 1
	if err != nil {
		return 0, fmt.Errorf("erro ao inserir na tabela de vendas: %w", err)
	}
	vendaID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter o ID da última venda inserida: %w", err)
	}

	// Carregar queries necessárias
	queryInsertVendaItem := sqlutils.GetQuery("InserirVendaItem")
	queryGetMedicamento := sqlutils.GetQuery("ObterMedicamentoCompleto")
	queryUpdateEstoque := sqlutils.GetQuery("AtualizarEstoqueMedicamento")

	// 2. Iterar sobre cada item da requisição.
	for _, itemReq := range req.Itens {
		var med Medicamento
		var categoriaNome sql.NullString
		var preco sql.NullFloat64

		// Buscar dados atuais do medicamento dentro da transação para garantir consistência.
		err := tx.QueryRow(queryGetMedicamento, itemReq.MedicamentoID).Scan(
			&med.ID, &med.Nome, &med.Fabricante, &med.CodigoANVISA, &med.Quantidade,
			&med.Validade, &med.CriadoEm, &med.Categoria.ID, &categoriaNome, &preco,
		)
		if err != nil {
			return 0, fmt.Errorf("medicamento com ID %d não encontrado na transação: %w", itemReq.MedicamentoID, err)
		}
		if preco.Valid {
			med.Preco = preco.Float64
		}

		// Validar estoque.
		if med.Quantidade < itemReq.Quantidade {
			return 0, fmt.Errorf("estoque insuficiente para o medicamento '%s'", med.Nome)
		}

		// Inserir o item na tabela 'venda_items'.
		_, err = tx.Exec(queryInsertVendaItem, vendaID, itemReq.MedicamentoID, itemReq.Quantidade, med.Preco)
		if err != nil {
			return 0, fmt.Errorf("erro ao inserir o item de venda '%s': %w", med.Nome, err)
		}

		// Atualizar o estoque do medicamento.
		novoEstoque := med.Quantidade - itemReq.Quantidade
		_, err = tx.Exec(queryUpdateEstoque, novoEstoque, itemReq.MedicamentoID)
		if err != nil {
			return 0, fmt.Errorf("erro ao atualizar o estoque do medicamento '%s': %w", med.Nome, err)
		}

		log.Printf("Item vendido: %s | Quantidade: %d | Preço Unitário: %.2f", med.Nome, itemReq.Quantidade, med.Preco)
	}

	// Se todos os itens foram processados sem erro, comitar a transação.
	return vendaID, tx.Commit()
}
