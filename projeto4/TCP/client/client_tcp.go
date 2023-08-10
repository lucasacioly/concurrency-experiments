package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

type NumbersResponse struct {
	PrimeNumbers    []int `json:"prime_numbers"`
	NonPrimeNumbers []int `json:"non_prime_numbers"`
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
		//println(numbers[i])
	}
	return numbers
}

func serverConnection(conn net.Conn, clientID int, numClients int) {

	defer conn.Close()

	// Gerar lista de números aleatórios
	numbers := generateRandomNumbers(TAM_AMOSTRA, 81, TAM_AMOSTRA/2)

	// Serializar a estrutura de requisição para JSON
	numbersJSON, err := json.Marshal(numbers)
	errorFound(err)

	var response NumbersResponse

	// Abrir arquivo para escrita
	fileName := fmt.Sprintf("TCP_elapsed_time_client_%d_%d.txt", numClients, clientID)
	file, err := os.Create(fileName)
	errorFound(err)

	for i := 0; i < NUM_REPS; i++ {
		//time.Sleep(time.Millisecond)
		// Iniciar contagem antes do envio
		startTime := time.Now().UnixNano()

		// Enviar requisição para o servidor
		_, err := conn.Write(numbersJSON)
		errorFound(err)

		// Aguardar resposta do servidor
		buffer := make([]byte, 16384)
		n, err := conn.Read(buffer)

		// registrar tempo decorrido
		elapsedTime := time.Now().UnixNano() - startTime

		_, err = file.WriteString(fmt.Sprintf("%d\n", elapsedTime))
		errorFound(err)

		// Deserializar a resposta do servidor para a estrutura de resposta
		err = json.Unmarshal(buffer[:n], &response)
		errorFound(err)

		// Exibir lista de números primos
		fmt.Println("Números Primos:", response.PrimeNumbers)

		// Exibir lista de números não primos
		fmt.Println("Números Não Primos:", response.NonPrimeNumbers)

	}
}

func main() {

	// Argumentos de linha de comando
	numClients := flag.Int("clients", 1, "Number of clients")
	clientID := flag.Int("id", 0, "Client ID")
	flag.Parse()
	println(*numClients, *clientID)

	// Conectar ao servidor
	conn, err := net.Dial("tcp", "localhost:8080")
	errorFound(err)

	serverConnection(conn, *clientID, *numClients)

	println("TERMINOU")

}
