package upload

import (
	rabbitmq "file_upload_project/core/services/rabbitmq"
	"net/http"
	"regexp"
)

func QueueUpload(w http.ResponseWriter, r *http.Request) {

	if !validateTypeOfFile(w, r) {
		return
	}

	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Erro ao obter o arquivo do formulário", http.StatusBadRequest)
		return
	}

	objectName := handler.Filename

	compressed_file, err := compress(w, r, file, objectName)

	if err != nil {
		http.Error(w, "Erro ao comprimir arquivo", http.StatusInternalServerError)
		return
	}

	defer compressed_file.Close()
	defer file.Close()

	objectNameTarGz := regexp.MustCompile(`\.(.*)$`).ReplaceAllString(objectName, ".tar.gz")

	rabbitmq.SendMessage(compressed_file, objectNameTarGz)
}
