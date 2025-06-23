// Aplica proteção em todas as páginas
document.addEventListener('DOMContentLoaded', function() {
    // Carrega o script de proteção
    const protectionScript = document.createElement('script');
    protectionScript.src = 'assets/js/protection.js';
    document.head.appendChild(protectionScript);

    // Ofusca todo o código JavaScript
    const scripts = document.getElementsByTagName('script');
    for (let script of scripts) {
        if (script.src) {
            // Carrega e ofusca scripts externos
            fetch(script.src)
                .then(response => response.text())
                .then(code => {
                    const obfuscatedCode = window.obfuscateCode(code);
                    const newScript = document.createElement('script');
                    newScript.textContent = obfuscatedCode;
                    script.parentNode.replaceChild(newScript, script);
                });
        } else if (script.textContent) {
            // Ofusca scripts inline
            const obfuscatedCode = window.obfuscateCode(script.textContent);
            script.textContent = obfuscatedCode;
        }
    }

    // Adiciona marca d'água
    const watermark = document.createElement('div');
    watermark.style.cssText = `
        position: fixed;
        bottom: 10px;
        right: 10px;
        color: rgba(0, 0, 0, 0.1);
        font-size: 12px;
        pointer-events: none;
        user-select: none;
        z-index: 9999;
    `;
    watermark.textContent = '© AlexCodeWorks - Todos os direitos reservados';
    document.body.appendChild(watermark);
    
}); 
