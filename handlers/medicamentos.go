package handlers

import (
	"net/http"
	"strconv"
	"time"

	"medicontrol/models"
	"medicontrol/services"

	"log"

	"github.com/gin-gonic/gin"
)

// ListarMedicamentos retorna a lista de todos os medicamentos
func ListarMedicamentos(c *gin.Context) {
	log.Println("Recebida requisição para listar medicamentos")

	// Adicionar headers CORS específicos
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	searchQuery := c.Query("search")
	var meds []models.Medicamento
	var err error

	// A lógica foi unificada dentro de models.BuscarMedicamentos
	// Se a searchQuery for vazia, ela retornará todos os medicamentos.
	log.Printf("Buscando medicamentos com o termo: '%s'", searchQuery)
	meds, err = models.BuscarMedicamentos(searchQuery)
	if err != nil {
		log.Printf("Erro ao buscar medicamentos do banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar sua solicitação"})
		return
	}

	log.Printf("Medicamentos recuperados: %d", len(meds))

	// Garantir que a resposta seja um array válido mesmo se meds for nil
	if meds == nil {
		log.Println("Lista de medicamentos é nil (após tentativa de recuperação), retornando array vazio")
		meds = []models.Medicamento{}
	}

	c.JSON(http.StatusOK, meds)
}

// ObterMedicamento retorna um medicamento específico
func ObterMedicamento(c *gin.Context) {
	id := c.Param("id")
	med := models.GetMedicamento(id)
	if med == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medicamento não encontrado"})
		return
	}
	c.JSON(http.StatusOK, med)
}

// BuscarDadosAnvisa busca os dados do medicamento na ANVISA
func BuscarDadosAnvisa(c *gin.Context) {
	codigo := c.Param("codigo")

	// Verificar se o medicamento já existe no sistema
	med := models.GetMedicamentoByCodigo(codigo)
	if med != nil {
		c.JSON(http.StatusOK, gin.H{
			"nome":       med.Nome,
			"fabricante": med.Fabricante,
			"exists":     true,
		})
		return
	}

	// Se não existir, buscar na API da ANVISA
	dados, err := services.BuscarDadosAnvisa(codigo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nome":       dados.Nome,
		"fabricante": dados.Fabricante,
		"exists":     false,
	})
}

// CriarMedicamento adiciona um novo medicamento
func CriarMedicamento(c *gin.Context) {
	var med models.Medicamento
	if err := c.ShouldBindJSON(&med); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar dados da ANVISA se não existirem
	if med.Nome == "" || med.Fabricante == "" {
		dados, err := services.BuscarDadosAnvisa(med.CodigoANVISA)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao buscar dados na ANVISA: " + err.Error()})
			return
		}
		med.Nome = dados.Nome
		med.Fabricante = dados.Fabricante
	}

	if err := models.AddMedicamento(&med); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, med)
}

// AtualizarMedicamento atualiza um medicamento existente
func AtualizarMedicamento(c *gin.Context) {
	id := c.Param("id")
	var med models.Medicamento
	if err := c.ShouldBindJSON(&med); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	med.ID = id
	if err := models.UpdateMedicamento(&med); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, med)
}

// DeletarMedicamento remove um medicamento
func DeletarMedicamento(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeleteMedicamento(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegistrarMovimentacao registra uma entrada ou saída de medicamento
func RegistrarMovimentacao(c *gin.Context) {
	var mov models.Movimentacao
	if err := c.ShouldBindJSON(&mov); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mov.Data = time.Now()
	if err := models.RegistrarMovimentacao(mov); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mov)
}

// ListarMovimentacoes retorna o histórico de movimentações
func ListarMovimentacoes(c *gin.Context) {
	movs, err := models.GetMovimentacoes()
	if err != nil {
		log.Printf("Erro ao buscar movimentações: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar histórico de movimentações"})
		return
	}

	if movs == nil {
		// Garantir que a resposta seja sempre um array, mesmo que vazio
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}

	c.JSON(http.StatusOK, movs)
}

// ObterRelatorioBaixoEstoque retorna uma lista de medicamentos com baixo estoque.
func ObterRelatorioBaixoEstoque(c *gin.Context) {
	// Definir um limite padrão, mas permitir que seja sobrescrito por um query param
	limiteStr := c.DefaultQuery("limite", "50")
	limite, err := strconv.Atoi(limiteStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'limite' inválido"})
		return
	}

	medicamentos, err := models.GetMedicamentosBaixoEstoque(limite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar relatório de baixo estoque"})
		return
	}

	if medicamentos == nil {
		c.JSON(http.StatusOK, []models.Medicamento{})
		return
	}

	c.JSON(http.StatusOK, medicamentos)
}

// ObterTotalVendas retorna o total de medicamentos vendidos (unidades).
func ObterTotalVendas(c *gin.Context) {
	total, err := models.GetTotalVendas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao calcular total de vendas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_vendas": total})
}
