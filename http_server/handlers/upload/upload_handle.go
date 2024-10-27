package upload

import (
	"context"
	"encoding/base64"
	minio_service "file_upload_project/core/services/minIO"
	"fmt"
	"net/http"

	"github.com/minio/minio-go/v7"
)

func UploadFileMinIO(objectNameTarGz string, compressed_file string) {

	data, err := base64.StdEncoding.DecodeString(compressed_file)

	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusBadRequest)
		return
	}

	info, err := minio_service.MinioClient.PutObject(
		context.Background(),
		minio_service.BucketName,
		objectNameTarGz,
		data,
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
