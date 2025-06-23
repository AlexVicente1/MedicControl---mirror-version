package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User representa um usuário do sistema
type User struct {
	Username      string
	PasswordHash  string
	Role          string
	LastLogin     time.Time
	CreatedAt     time.Time
	Active        bool
	LoginAttempts int
	LastAttempt   time.Time
}

// Session representa uma sessão de usuário
type Session struct {
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// UserStore gerencia o armazenamento de usuários
type UserStore struct {
	users map[string]User
	mutex sync.RWMutex
}

// SessionManager gerencia as sessões ativas
type SessionManager struct {
	sessions map[string]Session
	mutex    sync.RWMutex
}

// LoginRequest representa os dados de login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest representa os dados de registro
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var (
	store   *UserStore
	manager *SessionManager
	once    sync.Once
)

const (
	maxLoginAttempts = 5
	lockoutDuration  = 15 * time.Minute
)

// GetUserStore retorna a instância única do UserStore
func GetUserStore() *UserStore {
	once.Do(func() {
		store = &UserStore{
			users: make(map[string]User),
		}
		// Criar usuário admin padrão
		store.RegisterUser("admin", "admin123", "admin")
	})
	return store
}

// GetSessionManager retorna o gerenciador de sessões
func GetSessionManager() *SessionManager {
	if manager == nil {
		manager = &SessionManager{
			sessions: make(map[string]Session),
		}
		// Inicia uma goroutine para limpar sessões expiradas periodicamente
		go func() {
			for {
				time.Sleep(1 * time.Hour)
				manager.cleanExpiredSessions()
			}
		}()
	}
	return manager
}

// RegisterUser registra um novo usuário
func (us *UserStore) RegisterUser(username, password, role string) error {
	if username == "" || password == "" {
		return errors.New("usuário e senha são obrigatórios")
	}

	// Validação de senha
	if err := validatePassword(password); err != nil {
		return err
	}

	us.mutex.Lock()
	defer us.mutex.Unlock()

	// Verifica se o usuário já existe
	if _, exists := us.users[username]; exists {
		return errors.New("usuário já existe")
	}

	// Gera o hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Cria o novo usuário
	us.users[username] = User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
		Active:       true,
	}

	return nil
}

// ValidateUser valida as credenciais do usuário
func (us *UserStore) ValidateUser(username, password string) bool {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	user, exists := us.users[username]
	if !exists || !user.Active {
		return false
	}

	// Verifica se a conta está bloqueada
	if user.LoginAttempts >= maxLoginAttempts {
		if time.Since(user.LastAttempt) < lockoutDuration {
			return false
		}
		// Reseta as tentativas após o período de bloqueio
		user.LoginAttempts = 0
	}

	// Verifica a senha
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Incrementa tentativas de login falhas
		user.LoginAttempts++
		user.LastAttempt = time.Now()
		us.users[username] = user
		return false
	}

	// Login bem-sucedido: reseta contadores e atualiza último login
	user.LoginAttempts = 0
	user.LastLogin = time.Now()
	us.users[username] = user
	return true
}

// GetUserRole retorna o papel do usuário
func (us *UserStore) GetUserRole(username string) string {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	if user, exists := us.users[username]; exists {
		return user.Role
	}
	return ""
}

// validatePassword verifica se a senha atende aos requisitos mínimos
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("a senha deve ter pelo menos 8 caracteres")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char >= '!' && char <= '/' || char >= ':' && char <= '@':
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("a senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	}

	return nil
}

// CreateSession cria uma nova sessão para um usuário
func (sm *SessionManager) CreateSession(username string) (string, error) {
	// Gerar token aleatório
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	sm.mutex.Lock()
	sm.sessions[token] = Session{
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	sm.mutex.Unlock()

	return token, nil
}

// ValidateSession verifica se uma sessão é válida
func (sm *SessionManager) ValidateSession(token string) bool {
	sm.mutex.RLock()
	session, exists := sm.sessions[token]
	sm.mutex.RUnlock()

	if !exists {
		return false
	}

	// Verifica se a sessão não expirou
	if time.Now().After(session.ExpiresAt) {
		sm.mutex.Lock()
		delete(sm.sessions, token)
		sm.mutex.Unlock()
		return false
	}

	return true
}

// cleanExpiredSessions remove todas as sessões expiradas
func (sm *SessionManager) cleanExpiredSessions() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	now := time.Now()
	for token, session := range sm.sessions {
		if now.After(session.ExpiresAt) {
			delete(sm.sessions, token)
		}
	}
}

// RemoveSession remove uma sessão específica
func (sm *SessionManager) RemoveSession(token string) {
	sm.mutex.Lock()
	delete(sm.sessions, token)
	sm.mutex.Unlock()
}

// GetSession retorna uma sessão específica
func (sm *SessionManager) GetSession(token string) Session {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.sessions[token]
}

// LoginHandler processa o login do usuário
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Renderiza a página de login
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Erro ao renderizar página de login", http.StatusInternalServerError)
		}
		return
	case http.MethodPost:
		var req LoginRequest

		// Tenta ler como JSON primeiro
		if r.Header.Get("Content-Type") == "application/json" {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Dados inválidos", http.StatusBadRequest)
				return
			}
		} else {
			// Se não for JSON, assume que é form-data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
				return
			}
			req.Username = r.FormValue("username")
			req.Password = r.FormValue("password")
		}

		// Valida as credenciais
		if !GetUserStore().ValidateUser(req.Username, req.Password) {
			if r.Header.Get("Content-Type") == "application/json" {
				http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			} else {
				// Se for form-data, redireciona de volta para o login com mensagem de erro
				tmpl := template.Must(template.ParseFiles("templates/login.html"))
				data := struct {
					Error string
				}{
					Error: "Usuário ou senha inválidos",
				}
				tmpl.Execute(w, data)
			}
			return
		}

		// Cria uma nova sessão
		token, err := GetSessionManager().CreateSession(req.Username)
		if err != nil {
			http.Error(w, "Erro ao criar sessão", http.StatusInternalServerError)
			return
		}

		// Define o cookie de sessão
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		if r.Header.Get("Content-Type") == "application/json" {
			// Se for JSON, retorna resposta JSON
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Login realizado com sucesso",
				"role":    GetUserStore().GetUserRole(req.Username),
			})
		} else {
			// Se for form-data, redireciona para a página inicial
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// RegisterHandler processa o registro de um novo usuário
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verifica se o usuário atual é admin
	cookie, err := r.Cookie("session_token")
	if err != nil || !GetSessionManager().ValidateSession(cookie.Value) {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Registra o novo usuário
	if err := GetUserStore().RegisterUser(req.Username, req.Password, req.Role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Usuário registrado com sucesso",
	})
}

// LogoutHandler encerra a sessão do usuário
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Remove o cookie de sessão
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout realizado com sucesso",
	})
}

// AuthMiddleware é o middleware que protege as rotas
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtém o token do cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Valida o token
		if !GetSessionManager().ValidateSession(cookie.Value) {
			// Remove o cookie inválido
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-1 * time.Hour),
				HttpOnly: true,
				Path:     "/",
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Se chegou aqui, a sessão é válida
		next.ServeHTTP(w, r)
	})
}
