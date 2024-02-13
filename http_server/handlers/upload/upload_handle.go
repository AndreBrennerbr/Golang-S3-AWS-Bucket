package upload

import (
	"context"
	minio_service "file_upload_project/core/services/minIO"
	"fmt"
	"net/http"
	"regexp"

	"github.com/minio/minio-go/v7"
)

func UploadFileMinIO(w http.ResponseWriter, r *http.Request) {

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

	info, err := minio_service.MinioClient.PutObject(
		context.Background(),
		minio_service.BucketName,
		objectNameTarGz,
		compressed_file,
		-1,
		minio.PutObjectOptions{},
	)

	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Upload concluido com sucesso "+info.ETag)
	return
}
