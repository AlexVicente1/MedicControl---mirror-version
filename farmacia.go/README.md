# Sistema de Farmácia

Sistema simples para gerenciamento de produtos farmacêuticos.

## Requisitos

- Go 1.21 ou superior
- Navegador web moderno

## Instalação

1. Clone o repositório:
```bash
git clone <seu-repositorio>
cd farmacia
```

2. Instale as dependências:
```bash
go mod tidy
```

## Executando o sistema

1. Inicie o servidor:
```bash
go run main.go
```

2. Abra o navegador e acesse:
```
http://localhost:8080
```

## Estrutura do Projeto

```
.
├── data/
│   └── produtos.json     # Banco de dados de produtos
├── static/
│   ├── scripts.js        # JavaScript do frontend
│   └── style.css         # Estilos CSS
├── templates/
│   └── index.html        # Página principal
├── main.go               # Servidor backend
├── go.mod               # Dependências Go
└── README.md            # Este arquivo
```

## Funcionalidades

- Visualização de produtos por categoria
- Filtragem por subcategorias
- Interface responsiva
- Animações suaves
- Design moderno e intuitivo

## Categorias de Produtos

1. Medicamentos Controlados
2. Medicamentos Comuns
3. Cosméticos
4. Higiene Pessoal
5. Suplementos
6. Saúde e Bem-estar
7. Mãe e Bebê

## Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request 
