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
	ClientId int
	Numbers  []int
}

type NumbersResponse struct {
	PrimeNumbers    []int
	NonPrimeNumbers []int
}

const NUM_REPS = 10000
const TAM_AMOSTRA = 3000

const NUM_REPS_DEMO = 5
const TAM_AMOSTRA_DEMO = 20

const qos = 1

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

	// generate random numbers
	numbers := generateRandomNumbers(TAM_AMOSTRA, 81, TAM_AMOSTRA/2)

	// Open file for writing start times
	fileName := fmt.Sprintf("MQTT_start_time_client_%d_%d.txt", numClients, clientID)
	sFile, err := os.Create(fileName)
	errorFound(err)

	for i := 0; i < NUM_REPS; i++ {
		request := NumbersRequest{Numbers: numbers}
		request.ClientId = clientID
		requestBytes, err := json.Marshal(request)
		if err != nil {
			log.Printf("Error encoding request: %v", err)
			continue
		}

		startTime := time.Now().UnixNano()
		_, err = sFile.WriteString(fmt.Sprintf("%d\n", startTime))

		token := client.Publish("prime/request", qos, false, requestBytes)
		token.Wait()

		if token.Error() != nil {
			fmt.Println(token.Error())
			fmt.Println("ABORTAR OPERAÇÃO")
			os.Exit(1)
		}

		time.Sleep(100 * time.Millisecond)

	}
}

func main() {
	numClients := flag.Int("clients", 1, "Number of clients")
	clientID := flag.Int("id", 0, "Client ID")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883") // Update with your Mosquitto server details
	opts.SetClientID(fmt.Sprintf("prime_client_%d", *clientID))
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	defer client.Disconnect(250)

	// Open file for writing end times
	fileName := fmt.Sprintf("MQTT_end_time_client_%d_%d.txt", *numClients, *clientID)
	eFile, err := os.Create(fileName)
	errorFound(err)

	// subscrever a um topico & usar um handler para receber as mensagens
	responseTopic := fmt.Sprintf("prime/response/%d", *clientID)
	println(responseTopic)
	token = client.Subscribe(responseTopic, qos, func(c mqtt.Client, m mqtt.Message) {
		responseHandler(c, m, eFile)
	})
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	serverConnection(client, token, *clientID, *numClients)
	//serverConnectionDemo(client, numbers)

	fmt.Scanln()
}

func responseHandler(c mqtt.Client, m mqtt.Message, eFile *os.File) {
	response := NumbersResponse{}
	err := json.Unmarshal(m.Payload(), &response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	endTime := time.Now().UnixNano()
	_, err = eFile.WriteString(fmt.Sprintf("%d\n", endTime))

	// Display prime and non-prime numbers
	fmt.Println("Prime Numbers:", response.PrimeNumbers)
	fmt.Println("Non-Prime Numbers:", response.NonPrimeNumbers)
}
