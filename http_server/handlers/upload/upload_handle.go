package upload

import (
	"context"
	minio_service "file_upload_project/core/services/minIO"
	"fmt"
	"net/http"

	"github.com/minio/minio-go/v7"
)

func UploadFileMinIO(w http.ResponseWriter, r *http.Request) {

	if !ValidadeTypeOfFile(w, r) {
		return
	}

	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Erro ao obter o arquivo do formulário", http.StatusBadRequest)
		return
	}

	objectName := handler.Filename

	info, err := minio_service.MinioClient.PutObject(context.Background(), minio_service.BucketName, objectName, file, -1, minio.PutObjectOptions{})

	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Upload concluido com sucesso "+info.ETag)
	return
}
