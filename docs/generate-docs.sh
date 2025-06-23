#!/bin/bash

# Script para gerar PDF da documentação do Medic Control
# Requer: Pandoc e LaTeX (TeX Live)

# Verificar se o Pandoc está instalado
if ! command -v pandoc &> /dev/null; then
    echo "Pandoc não encontrado. Por favor, instale o Pandoc: https://pandoc.org/installing.html"
    exit 1
else
    echo "Pandoc encontrado: $(pandoc --version | head -n 1)"
fi

# Verificar se o XeLaTeX está instalado
if ! command -v xelatex &> /dev/null; then
    echo "XeLaTeX não encontrado. Por favor, instale o TeX Live: https://www.tug.org/texlive/"
    exit 1
else
    echo "XeLaTeX encontrado: $(xelatex --version | head -n 1)"
fi

# Criar diretório de saída se não existir
output_dir="../dist"
mkdir -p "$output_dir"
echo "Diretório de saída: $output_dir"

# Nome do arquivo de saída
output_file="$output_dir/medicontrol-documentacao.pdf"

# Gerar PDF
echo "Gerando PDF da documentação..."
pandoc \
    "../README.md" \
    "api.md" \
    "setup.md" \
    -o "$output_file" \
    --pdf-engine=xelatex \
    --css="style.css" \
    --toc \
    --toc-depth=3 \
    --number-sections \
    --highlight-style=tango \
    -V geometry:margin=1in \
    -V colorlinks=true \
    -V linkcolor=blue \
    -V toccolor=blue

# Verificar se o PDF foi gerado com sucesso
if [ -f "$output_file" ]; then
    echo -e "\nDocumentação gerada com sucesso!"
    echo "Arquivo: $output_file"
    
    # Abrir o PDF automaticamente
    echo -e "\nAbrindo o PDF..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        open "$output_file"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        xdg-open "$output_file"
    fi
else
    echo "Erro ao gerar o PDF."
    exit 1
fi 