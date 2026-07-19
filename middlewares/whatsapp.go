package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// WhatsAppMessage struct represents the message structure
type WhatsAppMessage struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

// RabbitMQPublisherMiddleware publishes a WhatsApp message to the RabbitMQ queue
func WhatsApp(number, message string) (string, error) {
	// Load RabbitMQ URL from environment variable or hardcode for testing
	err := godotenv.Load("./app/.env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return "", fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return "", fmt.Errorf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	queueName := "whatsapp_queue"
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return "", fmt.Errorf("Failed to declare a queue: %v", err)
	}

	// Create the WhatsAppMessage structure
	msg := WhatsAppMessage{
		Number:  number,
		Message: message,
	}

	// Marshal the message into JSON format
	messageBody, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("Failed to serialize message: %v", err)
	}

	// Publish the message to RabbitMQ
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		})
	if err != nil {
		return "", fmt.Errorf("Failed to publish message to queue: %v", err)
	}

	log.Printf("Message sent to %s: %s", number, message)
	return "Message successfully queued for " + number, nil
}

func main() {

	url := "https://warayang.com/api/send_message"
	method := "POST"

	payload := strings.NewReader("token=token%20service&number=08123456xxxx&message=pesan%20yang%20dikirim&date=yyyy-mm-dd&time=hh%3Amm%3Ass")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
