package messagequeue

import (
	"context"
	"encoding/json"
	"fmt"
	"history-service/models"
	"history-service/utils"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
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
		"code-submission", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
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

func ConsumeSubmissionMessages(client *firestore.Client, createSubmission func(
	*firestore.Client, context.Context, models.SubmissionHistory) (
	*firestore.DocumentRef, error)) {
	ctx := context.Background()

	// Consume messages from the queue
	msgs, err := rabbitMQChannel.Consume(
		codeSubmissionQueue.Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	utils.FailOnError(err, "RabbitMQ: Failed to register a consumer")

	// Create a channel to block indefinitely
	forever := make(chan bool)

	// Start a goroutine to handle incoming messages
	go func() {
		for d := range msgs {
			log.Printf("RabbitMQ: Received a message")

			// Parse request
			var submissionHistory models.SubmissionHistory
			if err := json.Unmarshal(d.Body, &submissionHistory); err != nil {
				log.Printf("RabbitMQ: Error decoding JSON: %v", err)
				continue
			}

			_, err := createSubmission(client, ctx, submissionHistory)
			if err != nil {
				log.Printf("RabbitMQ: %v", err)
			}
		}
	}()

	log.Printf("RabbitMQ: [*] Waiting for messages.")
	<-forever
}
