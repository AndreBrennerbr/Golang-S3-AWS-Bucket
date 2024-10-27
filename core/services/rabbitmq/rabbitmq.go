package rabbitmq

import (
	"context"
	"encoding/base64"
	"file_upload_project/core/entities"
	"io"
	"log"
	"mime/multipart"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ConnData *entities.RabbitMq
	Channel  *amqp.Channel
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Start() (amqp.Queue, error) {
	conn := Connect()
	Channel := Createchannel(conn)
	p := Createqueue(Channel)

	return p, nil
}

func Connect() *amqp.Connection {

	ConnData = &entities.RabbitMq{}
	*ConnData = ConnData.NewConn()

	connection, err := amqp.Dial("amqp:" + ConnData.Ip)

	failOnError(err, "Erro ao conectar no rabbitMq")

	defer connection.Close()

	return connection
}

func Createchannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()

	failOnError(err, "Erro ao cirar channel")

	defer channel.Close()

	return channel
}

func Createqueue(channel *amqp.Channel) amqp.Queue {

	queue_name := ConnData.QueueName

	p, err := channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // args
	)

	failOnError(err, "Erro ao criar fila")

	return p
}

func sendMessage(p amqp.Queue, ch *amqp.Channel, file multipart.File) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conteudo, err := io.ReadAll(file)

	failOnError(err, "Erro ao ler conteudo do arquivo")

	conteudoCodificado := base64.StdEncoding.EncodeToString(conteudo)

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		p.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(conteudoCodificado),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", conteudoCodificado)
}

func readMessage(p amqp.Queue) {

}
