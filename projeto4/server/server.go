package main

import (
	"encoding/json"
	"fmt"
	"net"
	//"strconv"
)

type NumbersRequest struct {
	Numbers []int `json:"numbers"`
}

type NumbersResponse struct {
	PrimeNumbers    []int `json:"prime_numbers"`
	NonPrimeNumbers []int `json:"non_prime_numbers"`
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	fmt.Println("Aguardando conexões TCP...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Receber requisição do cliente
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Deserializar a requisição do cliente para a estrutura de requisição
	var request NumbersRequest
	err = json.Unmarshal(buffer[:n], &request)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Separar números primos e não primos
	primeNumbers, nonPrimeNumbers := separatePrimeNumbers(request.Numbers)

	// Montar a estrutura de resposta
	response := NumbersResponse{
		PrimeNumbers:    primeNumbers,
		NonPrimeNumbers: nonPrimeNumbers,
	}

	// Serializar a estrutura de resposta para JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Enviar resposta para o cliente
	_, err = conn.Write(responseJSON)
	if err != nil {
		fmt.Println(err)
		return
	}
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

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}

	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}

	return true
}
