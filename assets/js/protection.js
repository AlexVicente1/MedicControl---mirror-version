// Proteção contra cópia
document.addEventListener('DOMContentLoaded', function() {
    // Previne seleção de texto
    document.addEventListener('selectstart', function(e) {
        e.preventDefault();
    });

    // Previne cópia
    document.addEventListener('copy', function(e) {
        e.preventDefault();
        const warning = "Todos os direitos são reservados a AlexCodeWorks. Tentativa de cópia detectada!";
        e.clipboardData.setData('text/plain', warning);
        
        // Mostra aviso visual
        const warningDiv = document.createElement('div');
        warningDiv.style.cssText = `
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background-color: #ff0000;
            color: white;
            padding: 20px;
            border-radius: 5px;
            z-index: 9999;
            font-weight: bold;
            text-align: center;
            box-shadow: 0 0 10px rgba(0,0,0,0.5);
        `;
        warningDiv.innerHTML = `
            <h3>⚠️ Aviso de Segurança ⚠️</h3>
            <p>${warning}</p>
            <p>Esta ação foi registrada.</p>
        `;
        document.body.appendChild(warningDiv);
        
        // Remove o aviso após 3 segundos
        setTimeout(() => {
            warningDiv.remove();
        }, 3000);
    });

    // Previne clique direito
    document.addEventListener('contextmenu', function(e) {
        e.preventDefault();
    });

    
    // Previne DevTools
    document.addEventListener('keydown', function(e) {
        if (e.ctrlKey && e.shiftKey && (e.key === 'I' || e.key === 'i' || e.key === 'J' || e.key === 'j')) {
            e.preventDefault();
        }
    });
}); 
