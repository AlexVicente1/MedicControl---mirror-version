package sqlutils

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	queries = make(map[string]string)
	once    sync.Once
)

// LoadSQLFiles carrega todas as queries SQL da pasta 'sql' para a memória.
// Esta função é projetada para ser chamada uma vez durante a inicialização da aplicação.
func LoadSQLFiles(sqlDir string) error {
	once.Do(func() {
		err := filepath.WalkDir(sqlDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
				content, readErr := os.ReadFile(path)
				if readErr != nil {
					log.Printf("Erro ao ler o arquivo SQL %s: %v", path, readErr)
					return readErr // ou continue, dependendo da política de erro desejada
				}
				queryName := strings.TrimSuffix(d.Name(), ".sql")
				queries[queryName] = string(content)
				log.Printf("Carregada query: %s", queryName)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Erro fatal ao carregar arquivos SQL da pasta %s: %v", sqlDir, err)
			// Em um cenário de produção, você pode querer retornar o erro em vez de Fatalf
			// para permitir um tratamento de erro mais granular pelo chamador.
		}
	})
	if len(queries) == 0 {
		log.Printf("Aviso: Nenhum arquivo .sql encontrado em %s ou a pasta não existe.", sqlDir)
		// Poderia ser um erro se arquivos SQL são esperados.
		// return fmt.Errorf("nenhum arquivo .sql encontrado em %s", sqlDir)
	}
	return nil
}

// GetQuery retorna uma query SQL carregada pelo nome.
// Retorna uma string vazia se a query não for encontrada.
func GetQuery(name string) string {
	query, ok := queries[name]
	if !ok {
		log.Printf("Aviso: Query SQL '%s' não encontrada.", name)
		return "" // Ou retorne um erro
	}

	return query
}
