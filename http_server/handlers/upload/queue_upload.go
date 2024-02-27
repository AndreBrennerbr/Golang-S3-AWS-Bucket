package upload

import (
	"file_upload_project/core/services/rabbitmq"
	"net/http"
)

func queueUpload(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")

	if err != nil {
		if err != nil {
			http.Error(w, "Erro ao obter o arquivo do formulário", http.StatusBadRequest)
			return
		}
	}

	rabbitmq.CreatePublisher(rabbitmq.Channel, file)
}
