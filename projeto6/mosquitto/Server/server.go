package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type NumbersRequest struct {
	Numbers []int
}

type NumbersResponse struct {
	PrimeNumbers    []int
	NonPrimeNumbers []int
}

const NUM_REPS = 10

type PrimeService NumbersResponse

func (ps *PrimeService) GetPrimeNumbers(req *NumbersRequest, resp *NumbersResponse) (err error) {
	primeNumbers, nonPrimeNumbers := separatePrimeNumbers(req.Numbers)
	resp.PrimeNumbers = primeNumbers
	resp.NonPrimeNumbers = nonPrimeNumbers
	return nil
}

func onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	var request NumbersRequest
	err := json.Unmarshal(msg.Payload(), &request)
	if err != nil {
		log.Printf("Error decoding message: %v", err)
		return
	}

	var response NumbersResponse
	response.PrimeNumbers, response.NonPrimeNumbers = separatePrimeNumbers(request.Numbers)

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}

	responseTopic := strings.Replace(msg.Topic(), "request", "response", 1)
	token := client.Publish(responseTopic, 0, false, responseBytes)
	token.Wait()
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
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883") // Update with your Mosquitto server details
	opts.SetClientID("prime_server")

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	defer client.Disconnect(250)

	requestTopic := "prime/request"
	token = client.Subscribe(requestTopic, 0, onMessageReceived)
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	fmt.Println("Aguardando conexões via MQTT...")
	select {}
}
