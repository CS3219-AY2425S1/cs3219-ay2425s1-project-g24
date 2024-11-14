package messagequeue

import (
	"execution-service/utils"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	CODE_SUBMISSION_QUEUE_KEY = "code-submission"
	NUM_RETRIES               = 10
)

var (
	codeSubmissionQueue amqp.Queue
	rabbitMQChannel     *amqp.Channel
)

func InitRabbitMQServer() (*amqp.Connection, *amqp.Channel) {
	conn := connectToRabbitMQ()

	// Create a channel
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	rabbitMQChannel = ch

	// Declare a queue
	q, err := ch.QueueDeclare(
		CODE_SUBMISSION_QUEUE_KEY, // name
		false,                     // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	codeSubmissionQueue = q

	return conn, ch
}

func connectToRabbitMQ() *amqp.Connection {
	var conn *amqp.Connection
	var err error
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	for i := 0; i < NUM_RETRIES; i++ { // Retry up to 10 times
		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			log.Println("Connected to RabbitMQ")
			return conn
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 5 seconds... (Attempt %d/%d)", i+1, NUM_RETRIES)
		time.Sleep(5 * time.Second)
	}
	utils.FailOnError(err, fmt.Sprintf("Failed to connect to RabbitMQ after %d attempts", NUM_RETRIES))
	return nil
}

func PublishSubmissionMessage(submission []byte) error {
	err := rabbitMQChannel.Publish(
		"",                       // exchange
		codeSubmissionQueue.Name, // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        submission,
		})
	if err != nil {
		return fmt.Errorf("Failed to publish a message: %v", err)
	}
	log.Printf("RabbitMQ: [x] Sent %s", submission)
	return nil
}
