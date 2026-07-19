// RabbitMQPublisherMiddleware is a middleware to publish messages to RabbitMQ
package middleware

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// func RabbitMQPublisherMiddleware(c *fiber.Ctx) error {
// 	type EmailRequest struct {
// 		To      string `json:"to"`
// 		Subject string `json:"subject"`
// 		Body    string `json:"body"`
// 	}

// 	var emailReq EmailRequest
// 	if err := c.BodyParser(&emailReq); err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
// 	}

// 	// Publish the message to RabbitMQ
// 	err := publishMessageToQueue(emailReq.To, emailReq.Subject, emailReq.Body)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString("Failed to send message to queue")
// 	}

// 	// Call the next handler in the chain (if any)
// 	return c.Next()
// }

// Helper function to publish message to RabbitMQ
func PublishMessageToQueue(to, subject, body string) error {
	err := godotenv.Load("./app/.env")
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mail_queue", // queue name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}
	// jsonBody, _ := json.Marshal("body": body})
	message := "To: " + to + ", Subject: " + subject + ", Body: " + body

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return err
}

func Mail(to, subject, body_mail string) (string, error) {

	url := "https://api.farmasiunissula.com/mail/send-email"
	method := "POST"

	// fmt.Println(to, subject, body_mail)
	payload := strings.NewReader(`{
    "to": "` + to + `",
    "subject": "` + subject + `",
    "body": "` + body_mail + `",
    "from": "SSO Farmasi UNISSULA"
}`)

	fmt.Println(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "1", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic YXBwLWZhcm1hc2k6ZmFybWFzaTIwMjU=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "2", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "3", err
	}
	fmt.Println(string(body))
	resultmail := string(body)
	return resultmail, nil
}
