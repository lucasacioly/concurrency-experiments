package main

import (
	"fmt"
	"sync"
	"time"
)

var NUM_CADEIRAS_VAGAS int = 10
var mu sync.Mutex

func Cliente(chCadeiras chan int, chAcorda chan int, id int) {
	fmt.Println("Cliente ", id, " chegou")

	// Verificar se o canal chAcorda está vazio e acordar o barbeiro se necessário
	select {
	case chCadeiras <- id:
		fmt.Println("Cliente ", id, " entrou na fila")
	case chAcorda <- 1:
		fmt.Println("Cliente ", id, " acordou o barbeiro")
	default:
		fmt.Println("Cliente ", id, " Foi embora, fila cheia")
	}

}

func cortar_cabelo(cliente int) {
	// Lógica para cortar o cabelo do cliente
	fmt.Println("Cortando cabelo do cliente", cliente)
	time.Sleep(time.Second) // Simulação do tempo de corte do cabelo
	fmt.Println("Cabelo do cliente", cliente, "cortado")
}

func Barbeiro(chAcorda chan int, chCadeiras chan int) {
	dormindo := 1
	for {
		if dormindo == 1 {
			<-chAcorda
			// Barbeiro é acordado por um cliente
			fmt.Println("Barbeiro acordado")
			clienteId := <-chCadeiras
			cortar_cabelo(clienteId)
		} else {
			select {
			case clienteId := <-chCadeiras:
				// Cortar o cabelo do cliente
				cortar_cabelo(clienteId)

			default:
				// A fila está vazia, dormir
				dormindo = 0
				fmt.Println("Dormindo...")
			}
		}
	}
}

func main() {

	chAcorda := make(chan int)
	chCadeiras := make(chan int, NUM_CADEIRAS_VAGAS)

	go Barbeiro(chAcorda, chCadeiras)

	clientId := 0

	for i := 0; i < 15; i++ {
		go Cliente(chCadeiras, chAcorda, clientId)
		clientId += 1
	}

	time.Sleep(30 * time.Second)

	for i := 0; i < 15; i++ {
		go Cliente(chCadeiras, chAcorda, clientId)
		clientId += 1
	}

	_, _ = fmt.Scanln()
}
