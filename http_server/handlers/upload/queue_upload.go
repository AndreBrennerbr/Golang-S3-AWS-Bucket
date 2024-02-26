package upload

import "file_upload_project/core/services/rabbitmq"

func queueUpload() {
	rabbitmq.CreatePublisher()
}
