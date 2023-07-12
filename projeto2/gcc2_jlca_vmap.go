package main

/*
import (
	"fmt"
	"sync"
)

var NUM_CADEIRAS_VAGAS int = 10
var DORMINDO bool = true
var CORTANDO bool = false

func desejoCortarCabelo(mu *sync.Mutex, chdorme chan int, chcorta chan int, wg *sync.WaitGroup, id int) {
	mu.Lock()
	NUM_CADEIRAS_VAGAS -= 1
	mu.Unlock()

	if DORMINDO {
		chdorme <- 1
	}

	fmt.Println("Cliente", id, "entrou na barbearia e deseja cortar o cabelo")

	select {
	case chcorta <- id:
		mu.Lock()
		NUM_CADEIRAS_VAGAS += 1
		mu.Unlock()
	}
	wg.Done()
}

func Cliente(mu *sync.Mutex, chdorme chan int, chcorta chan int, wg *sync.WaitGroup, id int) {

	if NUM_CADEIRAS_VAGAS == 0 {
		println("cliente ", id, " foi embora")
		wg.Done()
		return
	}

	if NUM_CADEIRAS_VAGAS > 0 {
		desejoCortarCabelo(mu, chdorme, chcorta, wg, id)
	}

}

func cortar_cabelo(cliente int) {
	// Lógica para cortar o cabelo do cliente
	CORTANDO = true
	fmt.Println("Cortando cabelo do cliente", cliente)
	//time.Sleep(2 * time.Second) // Simulação do tempo de corte do cabelo
	fmt.Println("Cabelo do cliente", cliente, "cortado")
	CORTANDO = false
}

func Barbeiro(mu *sync.Mutex, chdorme chan int, chcorta chan int) {

	for {

		if NUM_CADEIRAS_VAGAS == 10 && CORTANDO == false {
			fmt.Println("Barbeiro dormindo...")
			DORMINDO = true
			select {
			case <-chdorme:
				DORMINDO = false
			}
		} else {
			idcliente := <-chcorta
			cortar_cabelo(idcliente)
		}

	}
}

func run() {

	mu := new(sync.Mutex)
	chdorme := make(chan int)
	chcorta := make(chan int)
	wg := sync.WaitGroup{}

	go Barbeiro(mu, chdorme, chcorta)

	clientId := 0

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go Cliente(mu, chdorme, chcorta, &wg, clientId)
		clientId += 1
	}

	_, _ = fmt.Scanln()
}
*/
