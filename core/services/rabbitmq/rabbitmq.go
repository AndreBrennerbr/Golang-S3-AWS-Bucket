package rabbitmq

import (
	"context"
	"encoding/base64"
	"file_upload_project/core/entities"
	"file_upload_project/http_server/handlers/upload"
	"io"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
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

func Start() {
	conn := Connect()
	Channel := Createchannel(conn)
	Queue := Createqueue(Channel)
	msgs := createConsumer(Channel, Queue)
	createListener(msgs)
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

	queue, err := channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // args
	)

	failOnError(err, "Erro ao criar fila")

	return queue
}

func createConsumer(ch *amqp.Channel, Queue amqp.Queue) <-chan amqp091.Delivery {

	msgs, err := ch.Consume(
		Queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	failOnError(err, "Failed to register a consumer")

	return msgs
}

func createListener(msgs <-chan amqp091.Delivery) {
	var forever chan struct{}

	go func() {
		for d := range msgs {
			upload.UploadFileMinIO()
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func SendMessage(file *os.File, objectNameTarGz string) {

	conn := Connect()
	Channel := Createchannel(conn)
	Queue := Createqueue(Channel)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conteudo, err := io.ReadAll(file)

	failOnError(err, "Erro ao ler conteudo do arquivo")

	conteudoCodificado := base64.StdEncoding.EncodeToString(conteudo)

	err = Channel.PublishWithContext(ctx,
		"",         // exchange
		Queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(conteudoCodificado),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", conteudoCodificado)
}
