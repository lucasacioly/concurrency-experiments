package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type NumbersRequest struct {
	Numbers []int
}

type NumbersResponse struct {
	PrimeNumbers    []int
	NonPrimeNumbers []int
}

func connectToRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func consumeAndRespond(ch *amqp.Channel) {
	filaResp, err := ch.QueueDeclare(
		"prime_queue", // Queue name
		false,         // Durable
		false,         // Delete when unused
		false,         // Exclusive
		false,         // No-wait
		nil,           // Arguments
	)
	errorFound(err)

	msgs, err := ch.Consume(
		filaResp.Name, // Queue name
		"",            // Consumer
		true,          // Auto-ack
		false,         // Exclusive
		false,         // No-local
		false,         // No-wait
		nil,           // Args
	)
	errorFound(err)

	for msg := range msgs {
		var request NumbersRequest
		err := json.Unmarshal(msg.Body, &request)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}

		var response NumbersResponse
		response.PrimeNumbers, response.NonPrimeNumbers = separatePrimeNumbers(request.Numbers)

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
			continue
		}

		err = ch.Publish(
			"",          // Exchange
			msg.ReplyTo, // Routing key
			false,       // Mandatory
			false,       // Immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: msg.CorrelationId,
				Body:          responseBytes,
			})
		errorFound(err)
	}
}

func errorFound(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err)
		os.Exit(1)
	}
}

func isPrime(num int) bool {
	if num == 1 {
		return false
	}

	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}

	return true
}

func separatePrimeNumbers(numbers []int) ([]int, []int) {
	var primeNumbers []int
	var nonPrimeNumbers []int

	for _, num := range numbers {
		if isPrime(num) {
			primeNumbers = append(primeNumbers, num)
		} else {
			nonPrimeNumbers = append(nonPrimeNumbers, num)
		}
	}

	return primeNumbers, nonPrimeNumbers
}

func main() {
	conn, ch, err := connectToRabbitMQ()
	errorFound(err)
	defer conn.Close()
	defer ch.Close()

	// Consume messages and respond
	go consumeAndRespond(ch)

	fmt.Println("Aguardando conexões via RabbitMQ...")
	select {}
}
