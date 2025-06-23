package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"medicontrol/handlers"
	"medicontrol/models"
	"medicontrol/sqlutils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Chave secreta para assinar os tokens JWT
var jwtSecret = []byte("seu_segredo_super_secreto")

// User representa a estrutura de um usuário
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims representa as claims do JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Credenciais do admin (em produção, isso estaria em um banco de dados)
var adminUser = User{
	Username: "admin",
	// Senha: senha123 (hash gerado com bcrypt)
	Password: "$2a$10$YourHashedPasswordHere", // Será substituído na inicialização
}

// LoginRequest representa a estrutura do corpo da requisição de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// authMiddleware verifica se o token JWT é válido
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:] // Remove "Bearer "

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	log.Println("Attempting to start MediControl server...")
	log.Println("Iniciando o servidor MediControl...")

	// Carregar arquivos SQL
	if err := sqlutils.LoadSQLFiles("sql"); err != nil {
		log.Fatalf("Erro ao carregar arquivos SQL: %v", err)
	}
	log.Println("Arquivos SQL carregados com sucesso.")

	// Inicializar o banco de dados
	if err := models.InitDB(); err != nil {
		log.Fatal("Erro ao inicializar banco de dados:", err)
	}

	// Corrigir preços de medicamentos existentes (migração de dados)
	if err := models.CorrigirPrecosMedicamentos(); err != nil {
		log.Printf("Aviso: Ocorreu um erro durante a correção de preços: %v", err)
	}

	// Importar medicamentos do arquivo JSON
	if err := models.ImportarMedicamentos(); err != nil {
		log.Printf("Erro ao importar medicamentos: %v", err)
	} else {
		log.Println("Medicamentos importados com sucesso")
	}

	// Gerar hash da senha do admin na inicialização
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Erro ao gerar hash da senha:", err)
	}
	adminUser.Password = string(hashedPassword)
	log.Println("Hash da senha do admin gerado com sucesso")

	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	log.Println("Middleware CORS configurado")

	// Obter o diretório de trabalho atual
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Erro ao obter o diretório de trabalho:", err)
	}

	// Servir arquivos estáticos em rotas específicas
	r.StaticFile("/", filepath.Join(wd, "static", "index.html"))
	r.StaticFile("/dashboard.html", filepath.Join(wd, "static", "dashboard.html"))
	r.Static("/css", filepath.Join(wd, "static", "css"))
	r.Static("/js", filepath.Join(wd, "static", "js"))

	// Configurar rota para a pasta de logos
	logoPath := filepath.Join(wd, "logo")
	log.Println("Configurando rota para a pasta de logos...")
	log.Println("Diretório de logos:", logoPath)

	// Verificar se o diretório de logos existe
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		log.Printf("Aviso: Diretório de logos não encontrado: %s", logoPath)
	} else {
		files, err := os.ReadDir(logoPath)
		if err != nil {
			log.Printf("Erro ao ler diretório de logos: %v", err)
		} else {
			log.Println("Arquivos encontrados na pasta de logos:")
			for _, file := range files {
				info, _ := file.Info()
				log.Printf("- %s (tamanho: %d bytes)", file.Name(), info.Size())
			}
		}
	}

	r.Static("/logo", logoPath)
	log.Println("Configuração de arquivos estáticos concluída")

	// API routes
	api := r.Group("/api")
	{
		// Rota de login
		api.POST("/login", func(c *gin.Context) {
			log.Println("Recebida requisição de login")
			var loginReq LoginRequest
			if err := c.ShouldBindJSON(&loginReq); err != nil {
				log.Printf("Erro no corpo da requisição: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			// Verificar credenciais
			if loginReq.Username != adminUser.Username {
				log.Printf("Tentativa de login com usuário inválido: %s", loginReq.Username)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			// Verificar senha
			if err := bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(loginReq.Password)); err != nil {
				log.Printf("Senha incorreta para o usuário: %s", loginReq.Username)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			// Criar claims para o token
			claims := &Claims{
				Username: loginReq.Username,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				},
			}

			// Gerar token JWT
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtSecret)
			if err != nil {
				log.Printf("Erro ao gerar token: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
				return
			}

			log.Printf("Login bem-sucedido para o usuário: %s", loginReq.Username)
			c.JSON(http.StatusOK, gin.H{
				"token": tokenString,
			})
		})

		// Rotas protegidas
		protected := api.Group("")
		protected.Use(authMiddleware())
		{
			// Rotas de medicamentos
			protected.GET("/medicamentos", handlers.ListarMedicamentos)
			protected.GET("/medicamentos/:id", handlers.ObterMedicamento)
			protected.POST("/medicamentos", handlers.CriarMedicamento)
			protected.PUT("/medicamentos/:id", handlers.AtualizarMedicamento)
			protected.DELETE("/medicamentos/:id", handlers.DeletarMedicamento)

			// Rotas de categorias
			protected.GET("/categorias", handlers.ListarCategorias)

			// Nova rota para buscar dados da ANVISA
			protected.GET("/anvisa/:codigo", handlers.BuscarDadosAnvisa)

			// Rotas de movimentação
			protected.POST("/movimentacoes", handlers.RegistrarMovimentacao)
			protected.GET("/movimentacoes", handlers.ListarMovimentacoes)

			// Rotas de relatórios
			protected.GET("/relatorios/vendas", handlers.ObterTotalVendas)
			protected.GET("/relatorios/baixo-estoque", handlers.ObterRelatorioBaixoEstoque)

			// Rota para Vendas
			protected.POST("/vendas", handlers.CriarVendaHandler)

			// Rota protegida de teste
			protected.GET("/protected", func(c *gin.Context) {
				log.Println("Acessando rota protegida")
				c.JSON(http.StatusOK, gin.H{
					"message": "Esta é uma rota protegida!",
				})
			})
		}
	}

	// Iniciar o servidor na porta 8080
	log.Println("Servidor iniciado na porta 8080 - Acesse http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
