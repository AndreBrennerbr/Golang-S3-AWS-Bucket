package delete_handler

import (
	"context"
	minio_service "file_upload_project/core/services/minIO"
	"fmt"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
)

func DeletObject(w http.ResponseWriter, r *http.Request) {

	caminho := strings.Split(r.URL.Path, "/")

	if len(caminho) <= 1 {
		http.Error(w, "Erro ao ler o path", http.StatusInternalServerError)
		return
	}

	objectName := caminho[2]

	err := minio_service.MinioClient.RemoveObject(context.Background(), minio_service.BucketName, objectName, minio.RemoveObjectOptions{})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(w, "Deletado com sucesso")
	w.WriteHeader(http.StatusOK)
}
