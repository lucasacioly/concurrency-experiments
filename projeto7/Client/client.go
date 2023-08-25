package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"projeto6/pb"
	"time"

	"google.golang.org/grpc"
)

const NUM_REPS = 10000
const TAM_AMOSTRA = 3000
const RANGE = 1500

const NUM_REPS_DEMO = 5
const TAM_AMOSTRA_DEMO = 20
const RANGE_DEMO = 100

func errorFound(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err)
		os.Exit(1)
	}
}

func generateRandomNumbers(count int, baseSeed int64, interval int) []int32 {
	numbers := make([]int32, count)
	seed := baseSeed
	rand.Seed(seed)
	for i := 0; i < count; i++ {
		seed += int64(i)
		rand.Seed(seed)
		numbers[i] = int32(rand.Intn(interval) + 1)
	}
	return numbers
}

func serverConnection(client pb.PrimeServiceClient, clientID int, numClients int) {
	numbers := generateRandomNumbers(TAM_AMOSTRA, 81, RANGE/2)

	// Open file for writing
	fileName := fmt.Sprintf("gRPC_elapsed_time_client_%d_%d.txt", numClients, clientID)
	file, err := os.Create(fileName)
	errorFound(err)
	defer file.Close()

	for i := 0; i < NUM_REPS; i++ {
		startTime := time.Now()

		resp, err := client.SeparatePrimeNumbers(context.Background(), &pb.NumbersRequest{Numbers: numbers})
		errorFound(err)

		elapsedTime := time.Since(startTime).Nanoseconds()
		_, err = file.WriteString(fmt.Sprintf("%d\n", elapsedTime))

		// Display prime and non-prime numbers
		fmt.Println(clientID, " Prime Numbers:", resp.PrimeNumbers)
		fmt.Println(clientID, " Non-Prime Numbers:", resp.NonPrimeNumbers)
	}
}

func serverConnectionDemo(client pb.PrimeServiceClient, clientID int, numClients int) {
	numbers := generateRandomNumbers(TAM_AMOSTRA_DEMO, 81, RANGE_DEMO/2)

	for i := 0; i < NUM_REPS_DEMO; i++ {

		resp, err := client.SeparatePrimeNumbers(context.Background(), &pb.NumbersRequest{Numbers: numbers})
		errorFound(err)

		// Display prime and non-prime numbers
		fmt.Println(clientID, " Prime Numbers:", resp.PrimeNumbers)
		fmt.Println(clientID, " Non-Prime Numbers:", resp.NonPrimeNumbers)
	}
}

func main() {
	numClients := flag.Int("clients", 1, "Number of clients")
	clientID := flag.Int("id", 0, "Client ID")
	flag.Parse()

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	errorFound(err)
	defer conn.Close()

	client := pb.NewPrimeServiceClient(conn)
	serverConnection(client, *clientID, *numClients)
	//serverConnectionDemo(client, *clientID, *numClients)
}
