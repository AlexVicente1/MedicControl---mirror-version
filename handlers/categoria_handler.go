package handlers

import (
	"log"
	"medicontrol/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListarCategorias retorna a lista de todas as categorias de medicamentos.
func ListarCategorias(c *gin.Context) {
	log.Println("Recebida requisição para listar categorias")

	// Adicionar headers CORS para permitir acesso de diferentes origens
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Lidar com requisições pre-flight OPTIONS
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	categorias, err := models.GetAllCategorias()
	if err != nil {
		log.Printf("Erro ao buscar categorias do banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar sua solicitação"})
		return
	}

	log.Printf("Categorias recuperadas: %d", len(categorias))

	// Garantir que a resposta seja um array vazio em vez de nulo, se não houver categorias
	if categorias == nil {
		categorias = []models.Categoria{}
	}

	c.JSON(http.StatusOK, categorias)
}
