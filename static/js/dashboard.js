// Verificar autenticação
if (!localStorage.getItem('token')) {
    window.location.href = '/';
}

// Headers padrão para requisições
const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${localStorage.getItem('token')}`
};

// Funções de utilidade
function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
    document.querySelectorAll('.sidebar a').forEach(link => link.classList.remove('active'));
    document.getElementById(`${pageId}-page`).classList.add('active');
    document.querySelector(`[data-page="${pageId}"]`).classList.add('active');
}

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = 'block';
    } else {
        console.error(`Modal ${modalId} não encontrado`);
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = 'none';
    } else {
        console.error(`Modal ${modalId} não encontrado`);
    }
}

// Adicionar estas funções de utilidade
function showError(message) {
    alert(message); // Podemos melhorar isso depois com um componente de toast
}

function showSuccess(message) {
    alert(message); // Podemos melhorar isso depois com um componente de toast
}

// Carregar medicamentos
async function loadMedicamentos(searchTerm = '') {
    try {
        console.log('Iniciando carregamento de medicamentos...');
        
        let url = '/api/medicamentos';
        if (searchTerm) {
            url += `?search=${encodeURIComponent(searchTerm)}`;
        }
        
        const response = await fetch(url, { 
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        
        console.log('Response status:', response.status);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const medicamentos = await response.json();
        console.log('Medicamentos carregados:', medicamentos);
        
        const tbody = document.querySelector('#medicamentosTable tbody');
        if (!tbody) {
            console.error('Elemento tbody não encontrado');
            return;
        }
        
        tbody.innerHTML = '';
        
        if (!Array.isArray(medicamentos)) {
            console.error('Resposta não é um array:', medicamentos);
            throw new Error('Formato de resposta inválido');
        }
        
        if (medicamentos.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="5" style="text-align: center;">Nenhum medicamento cadastrado</td>
                </tr>
            `;
            return;
        }
        
        medicamentos.forEach(med => {
            const precoFormatado = (med && typeof med.preco === 'number') 
                ? med.preco.toFixed(2).replace('.', ',') 
                : '0,00';

            tbody.innerHTML += `
                <tr>
                    <td>${med.nome || '-'}</td>
                    <td>${med.fabricante || '-'}</td>
                    <td>${med.codigo_anvisa || '-'}</td>
                    <td>${med.quantidade || 0}</td>
                    <td>R$ ${precoFormatado}</td>
                    <td>
                        <button onclick="editMedicamento('${med.id}')" class="secondary-btn">Editar</button>
                        <button onclick="deleteMedicamento('${med.id}')" class="error-btn">Excluir</button>
                    </td>
                </tr>
            `;
        });
        
        // Atualizar select de medicamentos para movimentações
        const select = document.getElementById('medicamentoSelect');
        if (select) {
            select.innerHTML = '<option value="">Selecione um medicamento</option>';
            medicamentos.forEach(med => {
                select.innerHTML += `<option value="${med.id}">${med.nome}</option>`;
            });
        }
        
        // Adicionar listener para o formulário de medicamento
        const medicamentoForm = document.getElementById('medicamentoForm');
        if (medicamentoForm) {
            medicamentoForm.addEventListener('submit', handleMedicamentoSubmit);
        }
        
    } catch (error) {
        console.error('Erro ao carregar medicamentos:', error);
        const errorMessage = error.message || 'Erro desconhecido';
        alert(`Erro ao carregar medicamentos: ${errorMessage}`);
    }
}

// Função auxiliar para obter o nome do medicamento pelo ID
function getNomeMedicamento(medicamentoId) {
    // Se o objeto models não estiver disponível, retorna o ID
    if (typeof models === 'undefined' || !models.GetMedicamento) {
        console.warn('Objeto models não está disponível');
        return `Medicamento #${medicamentoId}`;
    }
    
    const medicamento = models.GetMedicamento(medicamentoId);
    return medicamento ? medicamento.nome : `Medicamento #${medicamentoId}`;
}

// Carregar movimentações
async function loadMovimentacoes() {
    try {
        console.log('Carregando movimentações...');
        const response = await fetch('/api/movimentacoes', { 
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        
        console.log('Resposta da API de movimentações:', response.status);
        
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || `Erro HTTP: ${response.status}`);
        }
        
        const movimentacoes = await response.json();
        console.log('Movimentações recebidas:', movimentacoes);
        
        const tbody = document.querySelector('#movimentacoesTable tbody');
        if (!tbody) {
            console.error('Elemento tbody não encontrado para movimentações');
            return;
        }
        
        tbody.innerHTML = '';
        
        if (!Array.isArray(movimentacoes)) {
            console.error('Resposta não é um array:', movimentacoes);
            throw new Error('Formato de resposta inválido');
        }
        
        if (movimentacoes.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="5" style="text-align: center;">Nenhuma movimentação registrada</td>
                </tr>
            `;
            return;
        }
        
        movimentacoes.forEach(mov => {
            try {
                // Formatar a data
                const data = new Date(mov.data);
                const dataFormatada = data.toLocaleString('pt-BR');
                
                // Adicionar linha na tabela
                tbody.innerHTML += `
                    <tr>
                        <td>${dataFormatada}</td>
                        <td>${mov.nome_medicamento || 'Medicamento não encontrado'}</td>
                        <td>${mov.tipo === 'entrada' ? 'Entrada' : 'Saída'}</td>
                        <td>${mov.quantidade}</td>
                        <td>${mov.observacao || '-'}</td>
                    </tr>
                `;
            } catch (error) {
                console.error('Erro ao processar movimentação:', mov, error);
                // Adicionar linha de erro para esta movimentação
                tbody.innerHTML += `
                    <tr style="color: red;">
                        <td colspan="5">Erro ao exibir movimentação: ${error.message}</td>
                    </tr>
                `;
            }
        });
    } catch (error) {
        console.error('Erro ao carregar movimentações:', error);
        const errorMessage = error.message || 'Erro desconhecido';
        alert(`Erro ao carregar movimentações: ${errorMessage}`);
        
        // Mostrar mensagem de erro na tabela
        const tbody = document.querySelector('#movimentacoesTable tbody');
        if (tbody) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="5" style="color: red; text-align: center;">
                        Erro ao carregar movimentações: ${errorMessage}
                    </td>
                </tr>
            `;
        }
    }
}

// Carregar relatórios
async function loadRelatorios() {
    try {
        // Carregar medicamentos com baixo estoque
        const responseBaixoEstoque = await fetch('/api/relatorios/baixo-estoque?limite=50', { headers });
        if (!responseBaixoEstoque.ok) {
            throw new Error('Erro ao carregar relatório de baixo estoque');
        }
        const medicamentosBaixoEstoque = await responseBaixoEstoque.json();
        
        const listaBaixa = document.getElementById('medicamentosBaixa');
        listaBaixa.innerHTML = ''; // Limpar a lista antiga

        if (Array.isArray(medicamentosBaixoEstoque) && medicamentosBaixoEstoque.length > 0) {
            medicamentosBaixoEstoque.forEach(med => {
                const li = document.createElement('li');
                li.textContent = `${med.nome} (${med.fabricante}) - Quantidade: ${med.quantidade}`;
                listaBaixa.appendChild(li);
            });
        } else {
            listaBaixa.innerHTML = '<li>Nenhum medicamento com baixo estoque.</li>';
        }

        // Carregar total de vendas (manter a lógica existente e melhorá-la depois)
        const responseVendas = await fetch('/api/relatorios/vendas', { headers });
        if (responseVendas.ok) {
            const dataVendas = await responseVendas.json();
            document.getElementById('totalVendas').textContent = dataVendas.total_vendas || 0;
        } else {
             document.getElementById('totalVendas').textContent = 'Erro';
        }

    } catch (error) {
        alert(error.message || 'Erro ao carregar relatórios');
    }
}

// Event Listeners
document.querySelectorAll('.sidebar a').forEach(link => {
    link.addEventListener('click', (e) => {
        e.preventDefault();
        const page = e.target.dataset.page;
        showPage(page);
        
        if (page === 'medicamentos') loadMedicamentos();
        else if (page === 'movimentacoes') loadMovimentacoes();
        else if (page === 'relatorios') loadRelatorios();
    });
});

document.getElementById('logoutBtn').addEventListener('click', () => {
    localStorage.removeItem('token');
    window.location.href = '/';
});

document.getElementById('addMedicamentoBtn').addEventListener('click', () => {
    openModal('medicamentoModal');
});

document.getElementById('addMovimentacaoBtn').addEventListener('click', () => {
    openModal('movimentacaoModal');
});

// Manipulador para o formulário de medicamento (Criação e Edição)
async function handleMedicamentoSubmit(e) {
    e.preventDefault();

    const medicamentoId = document.getElementById('medicamentoId').value;
    const isEditing = !!medicamentoId;

    const url = isEditing ? `/api/medicamentos/${medicamentoId}` : '/api/medicamentos';
    const method = isEditing ? 'PUT' : 'POST';

    const body = {
        nome: document.getElementById('nome').value,
        fabricante: document.getElementById('fabricante').value,
        tipo: document.getElementById('tipo').value,
        codigo_anvisa: document.getElementById('codigo_anvisa').value,
        quantidade: parseInt(document.getElementById('quantidade').value, 10),
        validade: document.getElementById('validade').value,
        preco: parseFloat(document.getElementById('preco').value) || 0.0,
        categoria_id: document.getElementById('categoriaId').value
    };
    
    // Adiciona o ID ao corpo apenas se estiver editando
    if(isEditing) {
        body.id = medicamentoId;
    }


    try {
        const response = await fetch(url, {
            method,
            headers,
            body: JSON.stringify(body)
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || `Erro ao ${isEditing ? 'atualizar' : 'criar'} medicamento`);
        }

        showSuccess(`Medicamento ${isEditing ? 'atualizado' : 'criado'} com sucesso!`);
        closeModal('medicamentoModal');
        document.getElementById('medicamentoForm').reset();
        document.getElementById('medicamentoId').value = ''; // Limpa o campo oculto
        loadMedicamentos();

    } catch (error) {
        showError(`Erro: ${error.message}`);
    }
}

// Buscar dados da ANVISA
document.getElementById('buscarAnvisa').addEventListener('click', async () => {
    const codigo = document.getElementById('codigoAnvisa').value;
    
    try {
        const response = await fetch(`/api/anvisa/${codigo}`, { headers });
        const data = await response.json();
        
        if (response.ok) {
            document.getElementById('nomeMedicamento').value = data.nome;
            document.getElementById('fabricante').value = data.fabricante;
        } else {
            alert(data.error || 'Erro ao buscar dados da ANVISA');
        }
    } catch (error) {
        alert('Erro ao conectar ao servidor');
    }
});

// Formulário de nova movimentação
document.getElementById('movimentacaoForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const movimentacao = {
        medicamento_id: document.getElementById('medicamentoSelect').value,
        tipo: document.getElementById('tipoMovimentacao').value,
        quantidade: parseInt(document.getElementById('quantidadeMovimentacao').value),
        observacao: document.getElementById('observacao').value
    };

    try {
        const response = await fetch('/api/movimentacoes', {
            method: 'POST',
            headers,
            body: JSON.stringify(movimentacao)
        });

        if (response.ok) {
            closeModal('movimentacaoModal');
            loadMovimentacoes();
            loadMedicamentos(); // Atualizar quantidades
        } else {
            alert('Erro ao registrar movimentação');
        }
    } catch (error) {
        alert('Erro ao conectar ao servidor');
    }
});

// Função para buscar medicamentos
async function searchMedicamentos(searchTerm) {
    try {
        const response = await fetch(`/api/medicamentos?search=${encodeURIComponent(searchTerm)}`, { 
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const medicamentos = await response.json();
        const tbody = document.querySelector('#medicamentosTable tbody');
        tbody.innerHTML = '';
        
        if (!Array.isArray(medicamentos) || medicamentos.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="5" style="text-align: center;">Nenhum medicamento encontrado</td>
                </tr>
            `;
            return;
        }
        
        medicamentos.forEach(med => {
            tbody.innerHTML += `
                <tr>
                    <td>${med.nome || '-'}</td>
                    <td>${med.fabricante || '-'}</td>
                    <td>${med.codigo_anvisa || '-'}</td>
                    <td>${med.quantidade || 0}</td>
                    <td>
                        <button onclick="editMedicamento('${med.id}')" class="secondary-btn">Editar</button>
                        <button onclick="deleteMedicamento('${med.id}')" class="error-btn">Excluir</button>
                    </td>
                </tr>
            `;
        });
    } catch (error) {
        console.error('Erro ao buscar medicamentos:', error);
        showError('Erro ao buscar medicamentos. Tente novamente.');
    }
}

// Inicialização quando a página carregar
document.addEventListener('DOMContentLoaded', async () => {
    console.log('Página carregada, iniciando...');
    
    try {
        await loadMedicamentos();
    } catch (error) {
        console.error('Erro na inicialização:', error);
    }
    
    // Configurar event listeners
    const addMedicamentoBtn = document.getElementById('addMedicamentoBtn');
    if (addMedicamentoBtn) {
        addMedicamentoBtn.addEventListener('click', () => {
            console.log('Abrindo modal de novo medicamento');
            openModal('medicamentoModal');
        });
    }
    
    // Configurar busca
    const searchInput = document.getElementById('searchMedicamento');
    const searchButton = document.getElementById('searchButton');
    let searchTimeout;
    
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            clearTimeout(searchTimeout);
            const searchTerm = e.target.value.trim();
            
            if (searchTerm.length === 0) {
                loadMedicamentos();
                return;
            }
            
            searchTimeout = setTimeout(() => {
                searchMedicamentos(searchTerm);
            }, 300);
        });
    }
    
    if (searchButton) {
        searchButton.addEventListener('click', () => {
            const searchTerm = searchInput ? searchInput.value.trim() : '';
            if (searchTerm) {
                searchMedicamentos(searchTerm);
            } else {
                loadMedicamentos();
            }
        });
    }
    
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            localStorage.removeItem('token');
            window.location.href = '/';
        });
    }

    loadMovimentacoes();
    loadCategorias(); // Adicionado para carregar as categorias no select
});

// Excluir medicamento
async function deleteMedicamento(id) {
    if (!confirm('Tem certeza que deseja excluir este medicamento? Esta ação não pode ser desfeita.')) {
        return;
    }

    try {
        const response = await fetch(`/api/medicamentos/${id}`, {
            method: 'DELETE',
            headers
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Erro ao excluir medicamento');
        }

        showSuccess('Medicamento excluído com sucesso!');
        loadMedicamentos(); // Recarrega a lista
    } catch (error) {
        showError(`Erro: ${error.message}`);
    }
}

// Editar medicamento
// Nota: Esta função depende que exista um endpoint GET /api/medicamentos/:id
// e que o modal de formulário tenha os IDs corretos.
async function editMedicamento(id) {
    try {
        const response = await fetch(`/api/medicamentos/${id}`, { headers });
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Não foi possível carregar os dados do medicamento. Status: ${response.status}. Detalhes: ${errorText}`);
        }
        const med = await response.json();

        if (!med) {
            throw new Error("Medicamento não encontrado.");
        }

        document.getElementById('medicamentoId').value = med.id;
        document.getElementById('nome').value = med.nome || '';
        document.getElementById('fabricante').value = med.fabricante || '';
        document.getElementById('tipo').value = med.tipo || '';
        document.getElementById('codigo_anvisa').value = med.codigo_anvisa || '';
        document.getElementById('quantidade').value = med.quantidade || 0;
        document.getElementById('validade').value = med.validade || '';
        document.getElementById('preco').value = (med.preco || 0).toFixed(2);
        
        const categoriaSelect = document.getElementById('categoriaId');
        if (categoriaSelect) {
            categoriaSelect.value = med.categoria.id || '';
        }

        document.querySelector('#medicamentoModal h2').textContent = 'Editar Medicamento';
        openModal('medicamentoModal');

    } catch (error) {
        showError(`Erro ao preparar edição: ${error.message}`);
    }
}

// Carregar categorias no formulário
async function loadCategorias() {
    try {
        const response = await fetch('/api/categorias', { headers });
        if (!response.ok) {
            throw new Error('Erro ao buscar categorias');
        }
        const categorias = await response.json();
        const select = document.getElementById('categoriaId');
        if (select) {
            select.innerHTML = '<option value="">Selecione a Categoria</option>';
            if (Array.isArray(categorias)) {
                categorias.forEach(cat => {
                    select.innerHTML += `<option value="${cat.ID}">${cat.Nome}</option>`;
                });
            }
        }
    } catch (error) {
        console.error('Erro ao carregar categorias:', error);
    }
}
