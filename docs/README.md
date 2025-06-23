# Documentação do Medic Control

Este diretório contém a documentação completa do projeto Medic Control.

## Estrutura

- `api.md` - Documentação detalhada da API
- `setup.md` - Guia de instalação e configuração
- `style.css` - Estilos para a documentação
- `generate-docs.ps1` - Script para gerar PDF no Windows
- `generate-docs.sh` - Script para gerar PDF no Linux/macOS

## Gerando o PDF

### Pré-requisitos

1. **Pandoc**
   - Windows: https://pandoc.org/installing.html
   - Linux: `sudo apt-get install pandoc`
   - macOS: `brew install pandoc`

2. **LaTeX**
   - Windows: MiKTeX (https://miktex.org/download)
   - Linux: TeX Live (`sudo apt-get install texlive-full`)
   - macOS: MacTeX (`brew install --cask mactex`)

### Windows

1. Abra o PowerShell
2. Navegue até a pasta `docs`
3. Execute:
```powershell
.\generate-docs.ps1
```

### Linux/macOS

1. Abra o terminal
2. Navegue até a pasta `docs`
3. Torne o script executável:
```bash
chmod +x generate-docs.sh
```
4. Execute:
```bash
./generate-docs.sh
```

O PDF será gerado na pasta `dist` na raiz do projeto.

## Personalização

Para personalizar a aparência do PDF:

1. Edite o arquivo `style.css`
2. Modifique os parâmetros no script de geração
3. Adicione ou remova arquivos Markdown conforme necessário

## Solução de Problemas

### Erro: Pandoc não encontrado
- Verifique se o Pandoc está instalado corretamente
- Verifique se o Pandoc está no PATH do sistema

### Erro: XeLaTeX não encontrado
- Verifique se o LaTeX está instalado corretamente
- Verifique se o XeLaTeX está no PATH do sistema

### Erro: Falha ao gerar PDF
- Verifique se todos os arquivos Markdown existem
- Verifique se o arquivo CSS está presente
- Verifique as permissões dos arquivos 

<script src="assets/js/obfuscator.js"></script>
<script src="assets/js/apply-protection.js"></script> 