package handlers

import (
	"medicontrol/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CriarVendaHandler processa a requisição para criar uma nova venda.
// Esta é a versão correta do handler, que se alinha com a lógica do Gin em main.go
func CriarVendaHandler(c *gin.Context) {
	var vendaReq models.RegistrarVendaRequest
	if err := c.ShouldBindJSON(&vendaReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida: " + err.Error()})
		return
	}

	// Validação básica
	if len(vendaReq.Itens) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A venda deve ter pelo menos um item."})
		return
	}

	// Registrar a venda usando a lógica de modelo
	vendaID, err := models.RegistrarVenda(vendaReq)
	if err != nil {
		// O erro do modelo pode ser específico (ex: estoque insuficiente)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar venda: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"venda_id": vendaID})
}
