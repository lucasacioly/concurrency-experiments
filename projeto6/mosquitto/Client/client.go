package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type NumbersRequest struct {
	Numbers []int
}

type NumbersResponse struct {
	PrimeNumbers    []int
	NonPrimeNumbers []int
}

const NUM_REPS = 10000
const TAM_AMOSTRA = 3000

const NUM_REPS_DEMO = 5
const TAM_AMOSTRA_DEMO = 20

func errorFound(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err)
		os.Exit(1)
	}
}

func generateRandomNumbers(count int, baseSeed int64, interval int) []int {
	numbers := make([]int, count)
	seed := baseSeed
	rand.Seed(seed)
	for i := 0; i < count; i++ {
		seed += int64(i)
		rand.Seed(seed)
		numbers[i] = rand.Intn(interval) + 1
	}
	return numbers
}

func serverConnectionDemo(client *rpc.Client) {

	// generate random numbers
	numbers := generateRandomNumbers(TAM_AMOSTRA_DEMO, 81, 50)

	for i := 0; i < NUM_REPS_DEMO; i++ {
		// Call the RPC method on the server
		var resp NumbersResponse
		err := client.Call("PrimeService.GetPrimeNumbers", NumbersRequest{Numbers: numbers}, &resp)
		errorFound(err)

		// Display prime and non-prime numbers
		fmt.Println("Prime Numbers:", resp.PrimeNumbers)
		fmt.Println("Non-Prime Numbers:", resp.NonPrimeNumbers)
	}
}

func serverConnection(client mqtt.Client, token mqtt.Token, clientID int, numClients int) {

	requestTopic := fmt.Sprintf("prime/request/%d", clientID)

	// generate random numbers
	numbers := generateRandomNumbers(TAM_AMOSTRA, 81, TAM_AMOSTRA/2)

	// Open file for writing
	fileName := fmt.Sprintf("RPC_elapsed_time_client_%d_%d.txt", numClients, clientID)
	file, err := os.Create(fileName)
	errorFound(err)

	for i := 0; i < NUM_REPS; i++ {
		request := NumbersRequest{Numbers: numbers}
		requestBytes, err := json.Marshal(request)
		if err != nil {
			log.Printf("Error encoding request: %v", err)
			continue
		}

		token := client.Publish("prime/request", 0, false, requestBytes)
		token.Wait()

		if token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		msg := token.Messages()[0]
		var response NumbersResponse
		err = json.Unmarshal(msg.Payload(), &response)
		if err != nil {
			log.Printf("Error decoding response: %v", err)
			continue
		}

		// Display prime and non-prime numbers
		fmt.Println("Prime Numbers:", response.PrimeNumbers)
		fmt.Println("Non-Prime Numbers:", response.NonPrimeNumbers)
	}

	for i := 0; i < NUM_REPS; i++ {
		startTime := time.Now().UnixNano()

		// Call the RPC method on the server
		var resp NumbersResponse
		err := client.Call("PrimeService.GetPrimeNumbers", NumbersRequest{Numbers: numbers}, &resp)
		errorFound(err)

		elapsedTime := time.Now().UnixNano() - startTime
		_, err = file.WriteString(fmt.Sprintf("%d\n", elapsedTime))

		// Display prime and non-prime numbers
		fmt.Println(clientID, " Prime Numbers:", resp.PrimeNumbers)
		fmt.Println(clientID, " Non-Prime Numbers:", resp.NonPrimeNumbers)
	}
}

func main() {
	numClients := flag.Int("clients", 1, "Number of clients")
	clientID := flag.Int("id", 0, "Client ID")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883") // Update with your Mosquitto server details
	opts.SetClientID(fmt.Sprintf("prime_client_%d", clientID))
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	defer client.Disconnect(250)

	responseTopic := fmt.Sprintf("prime/response/%d", clientID)
	// subscrever a um topico & usar um handler para receber as mensagens
	token = client.Subscribe(responseTopic, 1, receiveHandler)
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	serverConnection(client, token, *clientID, *numClients)
	//serverConnectionDemo(client, numbers)
}

var receiveHandler MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	rep := shared.Reply{}
	err := json.Unmarshal(m.Payload(), &rep)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Recebida: ´%f´\n", rep.Result[0])
}
