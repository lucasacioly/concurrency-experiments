package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type NumbersRequest struct {
	Numbers []int
}

type NumbersResponse struct {
	PrimeNumbers    []int
	NonPrimeNumbers []int
}

const NUM_REPS = 10

type PrimeService NumbersResponse

func (ps *PrimeService) GetPrimeNumbers(req *NumbersRequest, resp *NumbersResponse) (err error) {
	primeNumbers, nonPrimeNumbers := separatePrimeNumbers(req.Numbers)
	resp.PrimeNumbers = primeNumbers
	resp.NonPrimeNumbers = nonPrimeNumbers
	return nil
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
	// cria um novo servidor RPC e registra PrimeService
	primeService := new(PrimeService)
	rpc.Register(primeService)
	rpc.HandleHTTP()

	// cria o listener
	ln, err := net.Listen("tcp", ":8080")
	errorFound(err)

	fmt.Println("Aguardando conexões RPC...")

	http.Serve(ln, nil)
}
