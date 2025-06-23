// Função para ofuscar strings
function obfuscateString(str) {
    return str.split('').map(char => {
        return '\\x' + char.charCodeAt(0).toString(16).padStart(2, '0');
    }).join('');
}

// Função para ofuscar código
function obfuscateCode(code) {
    // Substitui strings por versões ofuscadas
    code = code.replace(/'([^']*)'/g, function(match, str) {
        return "'" + obfuscateString(str) + "'";
    });
    code = code.replace(/"([^"]*)"/g, function(match, str) {
        return '"' + obfuscateString(str) + '"';
    });

    // Adiciona variáveis aleatórias
    const randomVars = [];
    for (let i = 0; i < 10; i++) {
        const varName = '_' + Math.random().toString(36).substr(2, 9);
        randomVars.push(varName);
    }

    // Adiciona código de decodificação
    const decoder = `
        function _d(s) {
            return s.replace(/\\\\x([0-9a-f]{2})/gi, function(m, p1) {
                return String.fromCharCode(parseInt(p1, 16));
            });
        }
    `;

    // Adiciona variáveis aleatórias e código de decodificação
    code = decoder + '\n' + randomVars.join(';\n') + ';\n' + code;

    return code;
}


// Exporta a função
window.obfuscateCode = obfuscateCode; 
