package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateMedicamento(t *testing.T) {
	// Configurar o router de teste
	r := gin.Default()
	r.POST("/medicamentos", CreateMedicamento)

	// Criar um medicamento de teste
	medicamento := Medicamento{
		Nome:       "Paracetamol",
		Descricao:  "Analgésico e antitérmico",
		Fabricante: "EMS",
		Dosagem:    "500mg",
		Forma:      "Comprimido",
		Validade:   "2025-12-31",
		Quantidade: 100,
	}

	// Converter para JSON
	medJSON, _ := json.Marshal(medicamento)

	// Criar requisição de teste
	req, _ := http.NewRequest("POST", "/medicamentos", bytes.NewBuffer(medJSON))
	req.Header.Set("Content-Type", "application/json")

	// Criar response recorder
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o status da resposta
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verificar se o medicamento foi criado
	var createdMed Medicamento
	json.Unmarshal(w.Body.Bytes(), &createdMed)
	assert.Equal(t, medicamento.Nome, createdMed.Nome)
	assert.Equal(t, medicamento.Descricao, createdMed.Descricao)
	assert.Equal(t, medicamento.Fabricante, createdMed.Fabricante)
}

func TestGetMedicamentos(t *testing.T) {
	// Configurar o router de teste
	r := gin.Default()
	r.GET("/medicamentos", GetMedicamentos)

	// Criar response recorder
	w := httptest.NewRecorder()

	// Criar requisição de teste
	req, _ := http.NewRequest("GET", "/medicamentos", nil)

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o status da resposta
	assert.Equal(t, http.StatusOK, w.Code)

	// Verificar se a resposta é um JSON válido
	var medicamentos []Medicamento
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &medicamentos))
}
