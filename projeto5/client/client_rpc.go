package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"time"
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
const TAM_AMOSTRA_DEMO = 100

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
	numbers := generateRandomNumbers(TAM_AMOSTRA_DEMO, 81, TAM_AMOSTRA_DEMO/2)

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

func serverConnection(client *rpc.Client, clientID int, numClients int) {

	// generate random numbers
	numbers := generateRandomNumbers(TAM_AMOSTRA, 81, TAM_AMOSTRA/2)

	// Open file for writing
	fileName := fmt.Sprintf("RPC_elapsed_time_client_%d_%d.txt", numClients, clientID)
	file, err := os.Create(fileName)
	errorFound(err)

	for i := 0; i < NUM_REPS; i++ {
		startTime := time.Now().UnixNano()

		// Call the RPC method on the server
		var resp NumbersResponse
		err := client.Call("PrimeService.GetPrimeNumbers", NumbersRequest{Numbers: numbers}, &resp)
		errorFound(err)

		elapsedTime := time.Now().UnixNano() - startTime
		_, err = file.WriteString(fmt.Sprintf("%d\n", elapsedTime))
		errorFound(err)
		/*
			if elapsedTime != 0 {
				_, err = file.WriteString(fmt.Sprintf("%d\n", elapsedTime))
				errorFound(err)
			}*/

		// Display prime and non-prime numbers
		fmt.Println(clientID, " Prime Numbers:", resp.PrimeNumbers)
		fmt.Println(clientID, " Non-Prime Numbers:", resp.NonPrimeNumbers)
	}
}

func main() {
	numClients := flag.Int("clients", 1, "Number of clients")
	clientID := flag.Int("id", 0, "Client ID")
	flag.Parse()

	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	errorFound(err)

	serverConnection(client, *clientID, *numClients)
	//serverConnectionDemo(client, numbers)
}
