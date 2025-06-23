class SearchMedicamentos {
    constructor() {
        this.searchInput = document.getElementById('searchMedicamento');
        this.clearButton = document.querySelector('.clear-search');
        this.medicamentosTable = document.querySelector('#medicamentosTable tbody');
        this.suggestionsContainer = document.querySelector('.search-suggestions');
        this.loadingIndicator = document.querySelector('.search-loading');
        this.currentRequest = null;
        this.debounceTimeout = null;
        
        this.init();
    }
    
    init() {
        // Event listeners
        this.searchInput.addEventListener('input', (e) => this.handleSearchInput(e));
        this.searchInput.addEventListener('focus', () => this.showSuggestions());
        this.searchInput.addEventListener('blur', () => {
            // Pequeno atraso para permitir o clique nas sugestões
            setTimeout(() => this.hideSuggestions(), 200);
        });
        
        // Tecla Enter na pesquisa
        this.searchInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                this.performSearch(this.searchInput.value);
            }
        });
        
        // Limpar busca quando clicar no ícone de X
        if (this.clearButton) {
            this.clearButton.addEventListener('click', (e) => {
                e.preventDefault();
                this.clearSearch();
            });
        }
        
        // Mostrar/ocultar botão de limpar
        this.searchInput.addEventListener('input', () => {
            this.toggleClearButton();
        });
    }
    
    handleSearchInput(e) {
        const searchTerm = e.target.value.trim();
        const clearButton = this.searchInput.parentElement.querySelector('.clear-search');
        
        // Mostrar/ocultar botão de limpar
        if (clearButton) {
            clearButton.style.display = searchTerm ? 'flex' : 'none';
        }
        
        // Se o campo estiver vazio, limpar resultados
        if (!searchTerm) {
            this.clearSearch();
            return;
        }
        
        // Cancelar a requisição anterior se existir
        if (this.currentRequest) {
            this.currentRequest.abort();
        }
        
        // Limpar o timeout anterior
        clearTimeout(this.debounceTimeout);
        
        // Iniciar o indicador de carregamento
        this.showLoading(true);
        
        // Usar debounce para evitar muitas requisições
        this.debounceTimeout = setTimeout(() => {
            this.fetchSuggestions(searchTerm);
        }, 300);
    }
    
    fetchSuggestions(term) {
        // Cancelar requisição anterior se existir
        if (this.currentRequest) {
            this.currentRequest.abort();
        }
        
        // Criar um novo AbortController para a requisição atual
        const controller = new AbortController();
        const signal = controller.signal;
        this.currentRequest = controller;
        
        fetch(`/api/medicamentos?search=${encodeURIComponent(term)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            signal
        })
        .then(response => {
            if (!response.ok) throw new Error('Erro na busca');
            return response.json();
        })
        .then(data => {
            this.displaySuggestions(data);
        })
        .catch(error => {
            if (error.name !== 'AbortError') {
                console.error('Erro ao buscar sugestões:', error);
            }
        })
        .finally(() => {
            this.showLoading(false);
            this.currentRequest = null;
        });
    }
    
    displaySuggestions(medicamentos) {
        if (!this.suggestionsContainer) return;
        
        this.suggestionsContainer.innerHTML = '';
        
        if (!Array.isArray(medicamentos) || medicamentos.length === 0) {
            const noResults = document.createElement('div');
            noResults.className = 'search-suggestion-item';
            noResults.textContent = 'Nenhum medicamento encontrado';
            this.suggestionsContainer.appendChild(noResults);
        } else {
            medicamentos.forEach(med => {
                const item = document.createElement('div');
                item.className = 'search-suggestion-item';
                item.innerHTML = `
                    <div><strong>${med.nome || 'Sem nome'}</strong></div>
                    <div class="text-muted">${med.fabricante || 'Fabricante não informado'}</div>
                `;
                
                item.addEventListener('mousedown', (e) => {
                    e.preventDefault();
                    this.searchInput.value = med.nome;
                    this.performSearch(med.nome);
                });
                
                this.suggestionsContainer.appendChild(item);
            });
        }
        
        this.showSuggestions();
    }
    
    performSearch(term) {
        const searchTerm = term.trim();
        
        if (!searchTerm) {
            this.clearSearch();
            return;
        }
        
        // Atualizar a URL sem recarregar a página
        const url = new URL(window.location);
        url.searchParams.set('q', searchTerm);
        window.history.pushState({}, '', url);
        
        // Mostrar loading
        this.showLoading(true);
        
        // Fazer a requisição de busca
        fetch(`/api/medicamentos?search=${encodeURIComponent(searchTerm)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        })
        .then(response => {
            if (!response.ok) throw new Error('Erro na busca');
            return response.json();
        })
        .then(data => {
            this.updateResults(data);
        })
        .catch(error => {
            console.error('Erro na busca:', error);
            this.showError('Erro ao buscar medicamentos');
        })
        .finally(() => {
            this.showLoading(false);
        });
    }
    
    updateResults(medicamentos) {
        if (!this.medicamentosTable) return;
        
        this.medicamentosTable.innerHTML = '';
        
        if (!Array.isArray(medicamentos) || medicamentos.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td colspan="5" style="text-align: center;">
                    Nenhum medicamento encontrado
                </td>
            `;
            this.medicamentosTable.appendChild(row);
            return;
        }
        
        medicamentos.forEach(med => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${med.nome || '-'}</td>
                <td>${med.fabricante || '-'}</td>
                <td>${med.codigo_anvisa || '-'}</td>
                <td>${med.quantidade || 0}</td>
                <td>
                    <button onclick="editMedicamento('${med.id}')" class="secondary-btn">Editar</button>
                    <button onclick="deleteMedicamento('${med.id}')" class="error-btn">Excluir</button>
                </td>
            `;
            this.medicamentosTable.appendChild(row);
        });
    }
    
    toggleClearButton() {
        if (this.clearButton) {
            this.clearButton.style.display = this.searchInput.value.trim() ? 'flex' : 'none';
        }
    }
    
    clearSearch() {
        this.searchInput.value = '';
        this.searchInput.focus();
        this.hideSuggestions();
        this.toggleClearButton();
        
        // Atualizar a URL
        const url = new URL(window.location);
        url.searchParams.delete('q');
        window.history.pushState({}, '', url);
        
        // Recarregar a lista completa
        this.performSearch('');
    }
    
    showLoading(show) {
        if (this.loadingIndicator) {
            this.loadingIndicator.style.display = show ? 'block' : 'none';
        }
    }
    
    showError(message) {
        // Implementar um toast ou mensagem de erro mais bonito
        alert(message);
    }
    
    showSuggestions() {
        if (this.suggestionsContainer && this.searchInput.value.trim()) {
            this.suggestionsContainer.style.display = 'block';
        }
    }
    
    hideSuggestions() {
        if (this.suggestionsContainer) {
            this.suggestionsContainer.style.display = 'none';
        }
    }
}

// Inicializar a busca quando o DOM estiver carregado
document.addEventListener('DOMContentLoaded', () => {
    // Verificar se estamos na página de medicamentos
    if (document.getElementById('searchMedicamento')) {
        const search = new SearchMedicamentos();
        
        // Verificar se há um termo de busca na URL
        const urlParams = new URLSearchParams(window.location.search);
        const searchTerm = urlParams.get('q');
        
        if (searchTerm) {
            document.getElementById('searchMedicamento').value = searchTerm;
            search.performSearch(searchTerm);
        }
    }
});
