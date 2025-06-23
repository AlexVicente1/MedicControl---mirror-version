document.addEventListener('DOMContentLoaded', () => {
    // Esta verificação garante que o código só será executado se o token estiver disponível.
    const token = localStorage.getItem('token');
    if (!token) return;

    const headers = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    };

    // Funções de feedback para o usuário.
    // Elas procuram por divs com IDs 'error-message' e 'success-message' no seu HTML.
    // (Estas divs já existem no dashboard.html)
    function showError(message) {
        const errorDiv = document.getElementById('error-message');
        if (errorDiv) {
            errorDiv.textContent = message;
            errorDiv.style.display = 'block';
            setTimeout(() => {
                errorDiv.style.display = 'none';
            }, 3000);
        } else {
            console.error("Elemento #error-message não encontrado. Fallback para alert.");
            alert(message);
        }
    }

    function showSuccess(message) {
        const successDiv = document.getElementById('success-message');
        if (successDiv) {
            successDiv.textContent = message;
            successDiv.style.display = 'block';
            setTimeout(() => {
                successDiv.style.display = 'none';
            }, 3000);
        } else {
            console.log("Elemento #success-message não encontrado:", message);
            alert(message); // Fallback
        }
    }
    
    // O código do PDV só é ativado se a página de vendas estiver presente.
    if (document.getElementById('vendas-page')) {
        const searchInput = document.getElementById('searchPdv');
        const suggestionsContainer = document.getElementById('suggestionsPdv');
        const carrinhoItemsContainer = document.querySelector('.carrinho-items');
        const totalVendaSpan = document.getElementById('totalVenda');
        const finalizarVendaBtn = document.getElementById('finalizarVendaBtn');
        
        let carrinho = [];
        let searchTimeout;

        // --- LÓGICA DE BUSCA ---
        searchInput.addEventListener('input', () => {
            clearTimeout(searchTimeout);
            const searchTerm = searchInput.value.trim();
            if (searchTerm.length < 2) {
                suggestionsContainer.innerHTML = '';
                suggestionsContainer.style.display = 'none';
                return;
            }
            searchTimeout = setTimeout(() => searchMedicamentos(searchTerm), 300);
        });

        async function searchMedicamentos(term) {
            try {
                const response = await fetch(`/api/medicamentos?search=${encodeURIComponent(term)}`, { headers });
                if (!response.ok) throw new Error('Erro ao buscar medicamentos.');
                const medicamentos = await response.json() || [];
                
                suggestionsContainer.innerHTML = '';
                if (medicamentos.length === 0) {
                    suggestionsContainer.style.display = 'none';
                    return;
                }

                medicamentos.forEach(med => {
                    const div = document.createElement('div');
                    div.className = 'suggestion-item';
                    div.innerHTML = `
                        <strong>${med.nome} (${med.fabricante})</strong><br>
                        <small>Estoque: ${med.quantidade} | Preço: R$ ${med.Preco.toFixed(2)}</small>
                    `;
                    div.onclick = () => adicionarAoCarrinho(med);
                    suggestionsContainer.appendChild(div);
                });
                suggestionsContainer.style.display = 'block';

            } catch (error) {
                console.error(error);
                suggestionsContainer.style.display = 'none';
            }
        }
        
        document.addEventListener('click', (e) => {
            if (!suggestionsContainer.contains(e.target) && e.target !== searchInput) {
                suggestionsContainer.style.display = 'none';
            }
        });

        // --- LÓGICA DO CARRINHO ---
        function adicionarAoCarrinho(medicamento) {
            searchInput.value = '';
            suggestionsContainer.style.display = 'none';

            if (medicamento.quantidade <= 0) {
                showError('Este medicamento está fora de estoque.');
                return;
            }

            const itemExistente = carrinho.find(item => item.id === medicamento.id);
            if (itemExistente) {
                if(itemExistente.quantidade < medicamento.quantidade) {
                    itemExistente.quantidade++;
                } else {
                    showError('Quantidade máxima em estoque atingida para este item.');
                }
            } else {
                carrinho.push({
                    id: medicamento.id,
                    nome: medicamento.nome,
                    preco: medicamento.Preco,
                    quantidade: 1,
                    estoque: medicamento.quantidade
                });
            }
            renderizarCarrinho();
        }

        function renderizarCarrinho() {
            carrinhoItemsContainer.innerHTML = '';
            let total = 0;

            if (carrinho.length === 0) {
                carrinhoItemsContainer.innerHTML = '<p>O carrinho está vazio.</p>';
                finalizarVendaBtn.disabled = true;
                totalVendaSpan.textContent = '0,00';
                return;
            }

            carrinho.forEach(item => {
                const itemTotal = item.preco * item.quantidade;
                total += itemTotal;

                const itemDiv = document.createElement('div');
                itemDiv.className = 'carrinho-item';
                itemDiv.innerHTML = `
                    <div class="carrinho-item-info">
                        <strong>${item.nome}</strong>
                        <span>R$ ${item.preco.toFixed(2).replace('.', ',')} x ${item.quantidade} = R$ ${itemTotal.toFixed(2).replace('.', ',')}</span>
                    </div>
                    <div class="carrinho-item-actions">
                        <input type="number" value="${item.quantidade}" min="1" max="${item.estoque}" data-id="${item.id}" class="quantidade-input">
                        <button class="remover-item-btn" data-id="${item.id}" title="Remover Item">&times;</button>
                    </div>
                `;
                carrinhoItemsContainer.appendChild(itemDiv);
            });

            totalVendaSpan.textContent = total.toFixed(2).replace('.', ',');
            finalizarVendaBtn.disabled = false;
        }
        
        carrinhoItemsContainer.addEventListener('input', (e) => {
            if (e.target.classList.contains('quantidade-input')) {
                const id = parseInt(e.target.dataset.id, 10);
                let novaQuantidade = parseInt(e.target.value, 10);
                const item = carrinho.find(i => i.id === id);

                if (item) {
                     if (isNaN(novaQuantidade) || novaQuantidade < 1) {
                        novaQuantidade = 1;
                        e.target.value = novaQuantidade;
                    }
                    if (novaQuantidade > item.estoque) {
                        novaQuantidade = item.estoque;
                        e.target.value = novaQuantidade;
                        showError(`Estoque máximo para ${item.nome} é ${item.estoque}.`);
                    }
                    item.quantidade = novaQuantidade;
                    renderizarCarrinho();
                }
            }
        });

        carrinhoItemsContainer.addEventListener('click', (e) => {
            if (e.target.closest('.remover-item-btn')) {
                const id = parseInt(e.target.closest('.remover-item-btn').dataset.id, 10);
                carrinho = carrinho.filter(i => i.id !== id);
                renderizarCarrinho();
            }
        });

        // --- FINALIZAR VENDA ---
        finalizarVendaBtn.addEventListener('click', async () => {
            if (carrinho.length === 0) return;

            const vendaData = {
                itens: carrinho.map(item => ({
                    medicamento_id: item.id,
                    quantidade: item.quantidade
                }))
            };

            try {
                finalizarVendaBtn.disabled = true;
                finalizarVendaBtn.textContent = 'Processando...';

                const response = await fetch('/api/vendas', {
                    method: 'POST',
                    headers: headers,
                    body: JSON.stringify(vendaData)
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || 'Não foi possível concluir a venda.');
                }
                
                showSuccess('Venda realizada com sucesso!');
                carrinho = [];
                renderizarCarrinho();
                // A função loadMedicamentos vem do dashboard.js e atualiza a lista principal
                if(typeof loadMedicamentos === 'function'){
                    loadMedicamentos(); 
                }

            } catch (error) {
                showError(`Erro na venda: ${error.message}`);
            } finally {
                finalizarVendaBtn.disabled = false;
                finalizarVendaBtn.textContent = 'Finalizar Venda';
            }
        });
    }
}); 