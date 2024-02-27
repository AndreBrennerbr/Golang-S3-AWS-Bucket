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

func Start() error {
	conn, err := Connect()
	if err != nil {
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	Channel, err := Createchannel(conn)
	if err != nil {
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	err = Createqueue(Channel)
	if err != nil {
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func Connect() (*amqp.Connection, error) {

	ConnData = &entities.RabbitMq{}
	*ConnData = ConnData.NewConn()

	connection, err := amqp.Dial("amqp:" + ConnData.Ip)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer connection.Close()

	return connection, nil
}

func Createchannel(connection *amqp.Connection) (*amqp.Channel, error) {
	channel, err := connection.Channel()

	if err != nil {
		panic(err)
		return nil, err
	}

	defer channel.Close()

	return channel, nil
}

func Createqueue(channel *amqp.Channel) error {

	queue_name := ConnData.QueueName

	_, err := channel.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		panic(err)
	}
	return nil
}

func CreatePublisher(channel *amqp.Channel, content_file multipart.File) error {

	queue_name := ConnData.QueueName
	conteudo, err := io.ReadAll(content_file)
	if err != nil {
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
		panic(err)
		return err
	}

	return nil
}

func consumer(channel *amqp.Channel) error {

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
		panic(err)
	}

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
