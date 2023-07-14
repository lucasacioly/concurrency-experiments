package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
	"sync"
)

type NumbersResponse struct {
	PrimeNumbers    []int `json:"prime_numbers"`
	NonPrimeNumbers []int `json:"non_prime_numbers"`
}

const NUM_REPS = 10000

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

// , wg *sync.WaitGroup, mu *sync.Mutex
func serverConnection(addr *net.UDPAddr, conn *net.UDPConn, numbersJSON []byte, clientID int, numClients int, wg *sync.WaitGroup) {

	defer wg.Done()

	var response NumbersResponse
	var checker int64 = 0

	// Abrir arquivos para escrita
	fileName := fmt.Sprintf("UDP_elapsed_time_client_%d_%d.txt", numClients, clientID)
	file, err := os.Create(fileName)
	errorFound(err)

	fileName2 := fmt.Sprintf("UDP_total-medium_time_client_%d_%d.txt", numClients, clientID)
	file2, err := os.Create(fileName2)
	errorFound(err)

	startTimeAll := time.Now().UnixNano()

	for i := 0; i < NUM_REPS; i++ {
		// Iniciar contagem antes do envio
		startTime := time.Now().UnixNano()

		// Enviar requisição para o servidor
		_, err := conn.WriteToUDP(numbersJSON, addr)
		errorFound(err)

		// Aguardar resposta do servidor
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		errorFound(err)

		// Deserializar a resposta do servidor para a estrutura de resposta
		err = json.Unmarshal(buffer[:n], &response)
		errorFound(err)

		/*
		// Exibir lista de números primos
		fmt.Println("Números Primos:", response.PrimeNumbers)

		// Exibir lista de números não primos
		fmt.Println("Números Não Primos:", response.NonPrimeNumbers)
		*/

		// registrar tempo decorrido
		elapsedTime := time.Now().UnixNano() - startTime
		if elapsedTime != 0{
			_, err = file.WriteString(fmt.Sprintf("%d\n", elapsedTime))
			errorFound(err)
			checker++
		}

	}

	elapsedTimeAll := time.Now().UnixNano() - startTimeAll
	elapsedTimeMedium := elapsedTimeAll/checker

	// Escrever elapsedTime no arquivo geral
	_, err = file2.WriteString(fmt.Sprintf("Tempo de duração geral: %d\nTempo de duração médio de cada envio: %d", elapsedTimeAll, elapsedTimeMedium))
	errorFound(err)

	// Exibir lista de números primos
	fmt.Println("Números Primos:", response.PrimeNumbers)

	// Exibir lista de números não primos
	fmt.Println("Números Não Primos:", response.NonPrimeNumbers)


}

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Argumentos de linha de comando
	numClients := flag.Int("clients", 1, "Number of clients")
	clientID := flag.Int("id", 0, "Client ID")
	flag.Parse()

	// Gerar lista de números aleatórios
	numbers := generateRandomNumbers(800, 81, 1500)

	// Serializar a estrutura de requisição para JSON
	numbersJSON, err := json.Marshal(numbers)
	errorFound(err)

	// Resolver endereço do servidor
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	errorFound(err)

	// Estabelecer conexão UDP
	conn, err := net.ListenUDP("udp", nil)
	errorFound(err)

	serverConnection(serverAddr, conn, numbersJSON, *clientID, *numClients, &wg)

	wg.Wait()

}
