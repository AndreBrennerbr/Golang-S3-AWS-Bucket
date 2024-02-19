package entities

import "file_upload_project/core/config"

type RabbitMq struct {
	Ip        string
	QueueName string
}

func (r RabbitMq) NewConn() RabbitMq {
	r.Ip = config.EnvRabbitMq()
	r.QueueName = config.EnvRabbitMqQueueName()

	return r
}
