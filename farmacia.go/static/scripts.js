// Configuração das cores por categoria
const categoriaConfig = {
    'medicamentos-controlados': '#dc3545',
    'medicamentos-comuns': '#28a745',
    'cosmeticos': '#fd7e14',
    'higiene-pessoal': '#17a2b8',
    'suplementos': '#6f42c1',
    'saude-bem-estar': '#20c997',
    'mae-bebe': '#e83e8c'
};

// Variável global para armazenar os produtos
let produtos = [];

// Inicialização
document.addEventListener('DOMContentLoaded', async () => {
    await carregarProdutos();
    inicializarEventos();
    
    // Mostrar todos os produtos de todas as subcategorias
    const subcategoriasPorCategoria = {
        'medicamentos-controlados': ['Ansiolíticos', 'Antidepressivos', 'Psicotrópicos', 'Opioides'],
        'medicamentos-comuns': ['Analgésicos', 'Anti-inflamatórios', 'Antialérgicos', 'Antigripais', 'Vitaminas'],
        'cosmeticos': ['Skincare', 'Proteção Solar', 'Maquiagem', 'Cabelos', 'Perfumaria'],
        'higiene-pessoal': ['Higiene Oral', 'Higiene Corporal', 'Cuidados Íntimos', 'Desodorantes', 'Sabonetes'],
        'suplementos': ['Proteínas', 'Vitaminas e Minerais', 'Aminoácidos', 'Emagrecedores', 'Energéticos'],
        'saude-bem-estar': ['Ortopedia', 'Primeiros Socorros', 'Medição e Diagnóstico', 'Cuidados com Idosos', 'Fitness e Recuperação'],
        'mae-bebe': ['Cuidados com o Bebê', 'Alimentação Infantil', 'Higiene Bebê', 'Gestante', 'Puericultura']
    };

    // Para cada categoria e suas subcategorias
    Object.entries(subcategoriasPorCategoria).forEach(([categoria, subcategorias]) => {
        // Criar container para todas as subcategorias
        const container = document.getElementById(`produtos-${categoria}`);
        container.style.display = 'block';
        container.innerHTML = '';

        // Criar grid principal
        const gridPrincipal = document.createElement('div');
        gridPrincipal.className = 'categoria-produtos';

        // Para cada subcategoria
        subcategorias.forEach(subcategoria => {
            // Filtrar produtos da subcategoria
            const produtosFiltrados = produtos.filter(produto => 
                produto.categoria.toLowerCase() === categoria.toLowerCase() && 
                produto.subcategoria === subcategoria
            );

            // Se houver produtos nesta subcategoria
            if (produtosFiltrados.length > 0) {
                // Criar seção da subcategoria
                const subcategoriaSection = document.createElement('div');
                subcategoriaSection.className = 'subcategoria-section';
                
                // Adicionar título da subcategoria
                const titulo = document.createElement('h3');
                titulo.className = 'subcategoria-titulo';
                titulo.style.color = categoriaConfig[categoria];
                titulo.innerHTML = `<i class="fas fa-tag"></i> ${subcategoria}`;
                subcategoriaSection.appendChild(titulo);

                // Criar grid de produtos
                const grid = document.createElement('div');
                grid.className = 'produtos-grid';

                // Adicionar produtos
                produtosFiltrados.forEach(produto => {
                    const card = criarCardProduto(produto, categoriaConfig[categoria]);
                    grid.appendChild(card);
                });

                subcategoriaSection.appendChild(grid);
                gridPrincipal.appendChild(subcategoriaSection);
            }
        });

        container.appendChild(gridPrincipal);
    });
});

// Função para carregar produtos da API
async function carregarProdutos() {
    try {
        const response = await fetch('/api/produtos');
        const data = await response.json();
        produtos = data.produtos || [];
        console.log('Produtos carregados:', produtos.length);
    } catch (error) {
        console.error('Erro ao carregar produtos:', error);
        mostrarErro('Erro ao carregar produtos. Por favor, recarregue a página.');
    }
}

// Função para inicializar eventos
function inicializarEventos() {
    // Adiciona eventos aos badges de subcategoria
    document.querySelectorAll('.badge').forEach(badge => {
        badge.addEventListener('click', () => {
            const categoria = badge.closest('.categoria-card').querySelector('.produtos-container').id.replace('produtos-', '');
            const subcategoria = badge.textContent.trim();
            
            // Remove seleção de todos os badges
            document.querySelectorAll('.badge').forEach(b => b.classList.remove('active'));
            // Adiciona seleção ao badge clicado
            badge.classList.add('active');
            
            mostrarProdutos(categoria, subcategoria);
        });
    });
}

// Função para mostrar produtos
function mostrarProdutos(categoria, subcategoria) {
    console.log(`Mostrando produtos: ${categoria} / ${subcategoria}`);
    
    // Esconde todos os containers de produtos
    document.querySelectorAll('.produtos-container').forEach(container => {
        container.style.display = 'none';
    });
    
    // Filtra os produtos
    const produtosFiltrados = produtos.filter(produto => 
        produto.categoria.toLowerCase() === categoria.toLowerCase() && 
        produto.subcategoria === subcategoria
    );
    
    // Obtém o container da categoria
    const container = document.getElementById(`produtos-${categoria}`);
    if (!container) {
        console.error(`Container não encontrado para ${categoria}`);
        return;
    }
    
    // Limpa o container
    container.innerHTML = '';
    
    // Se não houver produtos, mostra mensagem
    if (produtosFiltrados.length === 0) {
        container.innerHTML = '<div class="sem-produtos">Nenhum produto encontrado nesta subcategoria.</div>';
        container.style.display = 'block';
        return;
    }
    
    // Cria o grid de produtos
    const grid = document.createElement('div');
    grid.className = 'produtos-grid';
    
    // Adiciona os cards de produtos
    produtosFiltrados.forEach(produto => {
        const card = criarCardProduto(produto, categoriaConfig[categoria]);
        grid.appendChild(card);
    });
    
    container.appendChild(grid);
    container.style.display = 'block';
}

// Função para criar card de produto
function criarCardProduto(produto, corCategoria) {
    // Determinar o status do produto baseado na quantidade
    let statusClass = 'status-ok';
    let statusText = 'Em estoque';
    
    if (produto.quantidade <= 0) {
        statusClass = 'status-danger';
        statusText = 'Sem estoque';
    } else if (produto.quantidade <= 10) {
        statusClass = 'status-warning';
        statusText = 'Estoque baixo';
    }

    // Formatar data de vencimento
    const dataVencimento = new Date(produto.vencimento);
    const hoje = new Date();
    const diasParaVencer = Math.ceil((dataVencimento - hoje) / (1000 * 60 * 60 * 24));
    
    let vencimentoStatus = 'status-ok';
    let vencimentoText = 'Válido';
    
    if (diasParaVencer <= 0) {
        vencimentoStatus = 'status-danger';
        vencimentoText = 'Vencido';
    } else if (diasParaVencer <= 30) {
        vencimentoStatus = 'status-warning';
        vencimentoText = 'Próximo ao vencimento';
    }

    const card = document.createElement('div');
    card.className = 'produto-card';
    card.innerHTML = `
        <div class="produto-header">
            <h3>${produto.nome}</h3>
            <span class="categoria-badge" style="background-color: ${corCategoria}">${produto.subcategoria}</span>
        </div>
        <div class="produto-info">
            <div class="info-item">
                <i class="fas fa-flask" style="color: ${corCategoria}"></i>
                <span title="Laboratório">${produto.laboratorio}</span>
            </div>
            <div class="info-item">
                <i class="fas fa-box" style="color: ${corCategoria}"></i>
                <span title="Lote">Lote: ${produto.lote}</span>
            </div>
            <div class="info-item">
                <i class="fas fa-calendar" style="color: ${corCategoria}"></i>
                <div class="produto-status" title="Status do vencimento">
                    <div class="status-indicator ${vencimentoStatus}"></div>
                    <span>Venc: ${produto.vencimento} (${vencimentoText})</span>
                </div>
            </div>
            <div class="info-item">
                <i class="fas fa-cubes" style="color: ${corCategoria}"></i>
                <div class="produto-status" title="Status do estoque">
                    <div class="status-indicator ${statusClass}"></div>
                    <span>Qtd: ${produto.quantidade} (${statusText})</span>
                </div>
            </div>
        </div>
        <div class="produto-footer">
            <div class="produto-preco">
                <i class="fas fa-tag" style="color: ${corCategoria}"></i>
                R$ ${produto.preco.toFixed(2)}
            </div>
            <div class="product-buttons">
                <button class="bula-btn" onclick="mostrarBula('${produto.nome}')">
                    <i class="fas fa-file-medical"></i>
                    Bula
                </button>
            </div>
        </div>`;
    return card;
}

// Função para mostrar erro
function mostrarErro(mensagem) {
    const erro = document.createElement('div');
    erro.className = 'alert alert-danger';
    erro.textContent = mensagem;
    document.querySelector('.container').prepend(erro);
    setTimeout(() => erro.remove(), 5000);
}

// Função para abrir modal de cadastro
function showCadastroModal() {
    // Implementar quando necessário
    console.log('Função de abrir modal de cadastro será implementada em breve');
}

// Função para fazer logout
function logout() {
    window.location.reload();
}

// Função para buscar e exibir a bula
async function mostrarBula(nomeProduto) {
    try {
        const response = await fetch(`/api/bula/${encodeURIComponent(nomeProduto)}`);
        const bula = await response.json();

        // Criar o modal se não existir
        let modal = document.getElementById('bulaModal');
        if (!modal) {
            modal = document.createElement('div');
            modal.id = 'bulaModal';
            modal.className = 'modal';
            modal.innerHTML = `
                <div class="modal-content">
                    <span class="close">&times;</span>
                    <div class="bula-content">
                        <h2>Bula do Medicamento</h2>
                        <div class="bula-info">
                            <h3>Nome</h3>
                            <p id="bulaNome"></p>
                            
                            <h3>Laboratório</h3>
                            <p id="bulaLaboratorio"></p>
                            
                            <h3>Registro ANVISA</h3>
                            <p id="bulaRegistro"></p>
                            
                            <h3>Princípio Ativo</h3>
                            <p id="bulaPrincipio"></p>
                            
                            <h3>Classe Terapêutica</h3>
                            <p id="bulaClasse"></p>
                            
                            <h3>Indicações</h3>
                            <p id="bulaIndicacoes"></p>
                            
                            <h3>Contraindicações</h3>
                            <p id="bulaContraindicacoes"></p>
                            
                            <h3>Posologia</h3>
                            <p id="bulaPosologia"></p>
                            
                            <h3>Efeitos Colaterais</h3>
                            <p id="bulaEfeitos"></p>
                        </div>
                    </div>
                </div>
            `;
            document.body.appendChild(modal);

            // Adicionar evento para fechar o modal
            const span = modal.querySelector('.close');
            span.onclick = function() {
                modal.style.display = "none";
            }

            // Fechar modal ao clicar fora dele
            window.onclick = function(event) {
                if (event.target == modal) {
                    modal.style.display = "none";
                }
            }
        }

        // Preencher as informações da bula
        document.getElementById('bulaNome').textContent = bula.nome;
        document.getElementById('bulaLaboratorio').textContent = bula.laboratorio;
        document.getElementById('bulaRegistro').textContent = bula.registro;
        document.getElementById('bulaPrincipio').textContent = bula.principio;
        document.getElementById('bulaClasse').textContent = bula.classe_terapeutica;
        document.getElementById('bulaIndicacoes').textContent = bula.indicacoes;
        document.getElementById('bulaContraindicacoes').textContent = bula.contraindicacoes;
        document.getElementById('bulaPosologia').textContent = bula.posologia;
        document.getElementById('bulaEfeitos').textContent = bula.efeitos_colaterais;

        // Mostrar o modal
        modal.style.display = "block";
    } catch (error) {
        console.error('Erro ao buscar bula:', error);
        alert('Não foi possível carregar a bula do medicamento.');
    }
}

// Adicionar botão de bula em cada produto
function adicionarBotaoBula() {
    const produtos = document.querySelectorAll('.product-card');
    produtos.forEach(produto => {
        const nomeProduto = produto.querySelector('h3').textContent;
        const botoesProduto = produto.querySelector('.product-buttons');
        
        if (!botoesProduto.querySelector('.bula-btn')) {
            const bulaBtn = document.createElement('button');
            bulaBtn.className = 'bula-btn';
            bulaBtn.innerHTML = '<i class="fas fa-file-medical"></i> Bula';
            bulaBtn.onclick = () => mostrarBula(nomeProduto);
            botoesProduto.appendChild(bulaBtn);
        }
    });
}

// Modificar a função displayProducts para incluir o botão de bula
function displayProducts(products) {
    // ... existing code ...
    
    // Depois de exibir os produtos, adicionar os botões de bula
    adicionarBotaoBula();
} 