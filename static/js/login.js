// Configurações da API
const API_URL = 'http://localhost:8080/api';

document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const errorMessage = document.getElementById('error-message');
            
            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ username, password }),
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    // Salvar token no localStorage
                    localStorage.setItem('token', data.token);
                    // Função para alternar a visibilidade da senha
                    function togglePassword() {
                        const passwordInput = document.getElementById('password');
                        const icon = document.querySelector('.toggle-password');
                        
                        if (passwordInput.type === 'password') {
                            passwordInput.type = 'text';
                            icon.classList.remove('fa-eye');
                            icon.classList.add('fa-eye-slash');
                        } else {
                            passwordInput.type = 'password';
                            icon.classList.remove('fa-eye-slash');
                            icon.classList.add('fa-eye');
                        }
                    }

                    // Login functionality will be added here
                    window.location.href = '/dashboard.html';
                } else {
                    errorMessage.textContent = data.error || 'Erro ao fazer login';
                }
            } catch (error) {
                errorMessage.textContent = 'Erro ao conectar ao servidor';
                console.error('Erro:', error);
            }
        });
    }
});

// Função para verificar se o usuário está autenticado
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token && window.location.pathname !== '/index.html') {
        window.location.href = '/index.html';
    }
}

// Função para fazer logout
function logout() {
    localStorage.removeItem('token');
    window.location.href = '/index.html';
}

// Função para fazer requisições autenticadas
async function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/index.html';
        return;
    }

    const headers = {
        ...options.headers,
        'Authorization': `Bearer ${token}`
    };

    try {
        const response = await fetch(url, { ...options, headers });
        if (response.status === 401) {
            // Token inválido ou expirado
            logout();
            return;
        }
        return response;
    } catch (error) {
        console.error('Erro na requisição:', error);
        throw error;
    }
} 
