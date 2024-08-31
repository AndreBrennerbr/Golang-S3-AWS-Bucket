package rabbitmq

import (
	"encoding/base64"
	"file_upload_project/core/entities"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/streadway/amqp"
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

func Start() error {
	conn := Connect()
	Channel := Createchannel(conn)
	Createqueue(Channel)

	return nil
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

func Createqueue(channel *amqp.Channel) {

	queue_name := ConnData.QueueName

	_, err := channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // args
	)

	failOnError(err, "Erro ao criar fila")
}

func CreatePublisher(channel *amqp.Channel, content_file multipart.File) error {

	queue_name := ConnData.QueueName
	conteudo, err := io.ReadAll(content_file)
	if err != nil {
		log.Fatal(err)
		return err
	}

	conteudoCodificado := base64.StdEncoding.EncodeToString(conteudo)

	err = channel.Publish(
		"",         // exchange
		queue_name, // key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        []byte(conteudoCodificado),
		},
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func Consumer(channel *amqp.Channel) error {

	queue_name := ConnData.QueueName

	msgs, err := channel.Consume(
		queue_name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        //args
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	go listener(msgs)

	return nil

}

func listener(msgs <-chan amqp.Delivery) {
	// print consumed messages from queue
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
