package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	//"flag"
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
	/*
	// Argumentos de linha de comando
	numClients := flag.Int("clients", 1, "Number of clients")
	flag.Parse()
	*/

	wg := sync.WaitGroup{}

	ln, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 8080})
	errorFound(err)
	defer ln.Close()

	fmt.Println("Aguardando conexões UDP...")

	for i:=0; i < 1; i++ {
		wg.Add(1)
		go handleConnection(ln, &wg)
	}

	/* Tem que ver como pega esse num clients do powershell pra usar aqui
	for i:=0; i < int(numClients); i++ {
		go handleConnection(ln)
	} */

	wg.Wait()
}

func handleConnection(ln *net.UDPConn, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < NUM_REPS; i++ {
		// Receber requisição do cliente
		buffer := make([]byte, 8092)
		n, addr, err := ln.ReadFromUDP(buffer)
		errorFound(err)

		// Deserializar a requisição do cliente para a estrutura de requisição
		var request []int
		err = json.Unmarshal(buffer[:n], &request)
		errorFound(err)

		// Separar números primos e não primos
		primeNumbers, nonPrimeNumbers := separatePrimeNumbers(request)

		// Montar a estrutura de resposta
		response := NumbersResponse{
			PrimeNumbers:    primeNumbers,
			NonPrimeNumbers: nonPrimeNumbers,
		}

		// Serializar a estrutura de resposta para JSON
		responseJSON, err := json.Marshal(response)
		errorFound(err)

		// Enviar resposta para o cliente
		_, err = ln.WriteToUDP(responseJSON, addr)
		errorFound(err)
	}

}
