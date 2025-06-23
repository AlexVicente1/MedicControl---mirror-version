# Script para gerar PDF da documentação do Medic Control
# Requer: Pandoc e LaTeX (MiKTeX ou TeX Live)

# Verificar se o Pandoc está instalado
try {
    $pandocVersion = pandoc --version
    Write-Host "Pandoc encontrado: $($pandocVersion[0])"
} catch {
    Write-Host "Pandoc não encontrado. Por favor, instale o Pandoc: https://pandoc.org/installing.html"
    exit 1
}

# Verificar se o MiKTeX/TeX Live está instalado
try {
    $xelatexVersion = xelatex --version
    Write-Host "XeLaTeX encontrado: $($xelatexVersion[0])"
} catch {
    Write-Host "XeLaTeX não encontrado. Por favor, instale o MiKTeX: https://miktex.org/download"
    exit 1
}

# Criar diretório de saída se não existir
$outputDir = "..\dist"
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
    Write-Host "Diretório de saída criado: $outputDir"
}

# Nome do arquivo de saída
$outputFile = "$outputDir\medicontrol-documentacao.pdf"

# Gerar PDF
Write-Host "Gerando PDF da documentação..."
pandoc `
    "..\README.md" `
    "api.md" `
    "setup.md" `
    -o $outputFile `
    --pdf-engine=xelatex `
    --css="style.css" `
    --toc `
    --toc-depth=3 `
    --number-sections `
    --highlight-style=tango `
    -V geometry:margin=1in `
    -V colorlinks=true `
    -V linkcolor=blue `
    -V toccolor=blue

# Verificar se o PDF foi gerado com sucesso
if (Test-Path $outputFile) {
    Write-Host "`nDocumentação gerada com sucesso!"
    Write-Host "Arquivo: $outputFile"
    
    # Abrir o PDF automaticamente
    Write-Host "`nAbrindo o PDF..."
    Start-Process $outputFile
} else {
    Write-Host "Erro ao gerar o PDF."
    exit 1
} 