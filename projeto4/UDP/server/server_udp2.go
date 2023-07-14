package main
/*
import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	//"strconv"
)

type NumbersResponse struct {
	PrimeNumbers    []int `json:"prime_numbers"`
	NonPrimeNumbers []int `json:"non_prime_numbers"`
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
	ln, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 8080})
	errorFound(err)
	defer ln.Close()

	fmt.Println("Aguardando conexões UDP...")

	for {
		buffer := make([]byte, 1024)
		n, addr, err := ln.ReadFromUDP(buffer)
		errorFound(err)

		go handleConnection(ln, addr, buffer[:n])
	}
}

func handleConnection(ln *net.UDPConn, addr *net.UDPAddr, buffer []byte) {

	// Deserializar a requisição do cliente para a estrutura de requisição
	var request []int
	err := json.Unmarshal(buffer, &request)
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
*/