package httpserver

import (
	"context"
	"file_upload_project/core/config"
	minio_service "file_upload_project/core/services/minIO"
	"file_upload_project/core/services/rabbitmq"
	"file_upload_project/http_server/routes"
)

func Start() {
	config.LoadEnv()

	contexto := context.Background()

	minio_service.Start(contexto)

	err := rabbitmq.Start()

	if err != nil {
		panic(err)
	}

	err = rabbitmq.Consumer(rabbitmq.Channel)

	if err != nil {
		panic(err)
	}

	routes.Routes()
}
