package main

import (
	"math/rand"
	//"strconv"
)

type NumbersRequest struct {
	Numbers []int `json:"numbers"`
}

type NumbersResponse struct {
	PrimeNumbers    []int `json:"prime_numbers"`
	NonPrimeNumbers []int `json:"non_prime_numbers"`
}

func generateRandomNumbers(count int, baseSeed int64, interval int) []int {
	numbers := make([]int, count)
	seed := baseSeed
	rand.Seed(seed)
	for i := 0; i < count; i++ {
		seed += int64(i)
		rand.Seed(seed)
		numbers[i] = rand.Intn(interval) + 1
		println(numbers[i])
	}
	return numbers
}

func main() {

	// Gerar lista de números aleatórios
	numbers := generateRandomNumbers(10, 81, 100)

	println(numbers)

	/*
		// Montar a estrutura de requisição
		request := NumbersRequest{
			Numbers: numbers,
		}

		// Serializar a estrutura de requisição para JSON
		requestJSON, err := json.Marshal(request)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Conectar ao servidor
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		// Enviar requisição para o servidor
		_, err = conn.Write(requestJSON)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Aguardar resposta do servidor
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Deserializar a resposta do servidor para a estrutura de resposta
		var response NumbersResponse
		err = json.Unmarshal(buffer[:n], &response)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Exibir lista de números primos
		fmt.Println("Números Primos:", response.PrimeNumbers)

		// Exibir lista de números não primos
		fmt.Println("Números Não Primos:", response.NonPrimeNumbers)*/
}
