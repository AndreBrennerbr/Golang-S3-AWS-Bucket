package upload

import (
	"net/http"
	"path/filepath"
)

func validateTypeOfFile(w http.ResponseWriter, r *http.Request) bool {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter o arquivo do formulário", http.StatusBadRequest)
		return false
	}

	defer file.Close()

	acceptedExtensions := []string{".pdf", ".doc", ".docx", ".txt"}
	extension := filepath.Ext(handler.Filename)

	if !contains(acceptedExtensions, extension) {
		http.Error(w, "Tipo de arquivo não permitido", http.StatusBadRequest)
		return false
	}

	return true
}

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}
