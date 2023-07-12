package main

import (
	"math/rand"
	"fmt"
	"net"
	"encoding/json"
	"os"
	"sync"
	"time"
	//"strconv"
)

type NumbersResponse struct {
	PrimeNumbers    []int `json:"prime_numbers"`
	NonPrimeNumbers []int `json:"non_prime_numbers"`
}

const sendingNumber int = 10

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

func serverConnection(addr *net.UDPAddr, conn *net.UDPConn, numbersJSON []byte, wg *sync.WaitGroup, mu *sync.Mutex) {
	
	defer wg.Done()

	var response NumbersResponse
	//var PrimeNumbers    []int
	//var NonPrimeNumbers []int
	
	// Enviar requisição para o servidor
	mu.Lock()
	_, err := conn.WriteToUDP(numbersJSON, addr)
	errorFound(err)

	// Aguardar resposta do servidor
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	errorFound(err)
	mu.Unlock()

	// Deserializar a resposta do servidor para a estrutura de resposta
	err = json.Unmarshal(buffer[:n], &response)
	errorFound(err)

	mu.Lock()
	// Exibir lista de números primos
	fmt.Println("Números Primos:", response.PrimeNumbers)

	// Exibir lista de números não primos
	fmt.Println("Números Não Primos:", response.NonPrimeNumbers)
	mu.Unlock()
}

func main() {

	// Gerar lista de números aleatórios
	numbers := generateRandomNumbers(50, 81, 1000)
	wg := sync.WaitGroup{}
	mu := new(sync.Mutex)
	

	// Serializar a estrutura de requisição para JSON
	numbersJSON, err := json.Marshal(numbers)
	errorFound(err)

	startTime := time.Now().UnixNano()
	for i := 0; i < sendingNumber; i++{
		
		// Resolver endereço do servidor
		serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
		errorFound(err)

		// Estabelecer conexão UDP
		conn, err := net.ListenUDP("udp", nil)
		errorFound(err)
		wg.Add(1)

		go serverConnection(serverAddr, conn, numbersJSON, &wg, mu)
	}

	wg.Wait()
	elapsedTime := time.Now().UnixNano() - startTime
	fmt.Printf("%d\n", elapsedTime)

}
