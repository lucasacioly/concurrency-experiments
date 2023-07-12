package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numPessoas int = 5
	addPorçoes int = 10 //numero de porsoes para adicionar
)

var (
	// wgpessoas  sync.WaitGroup
	wg         sync.WaitGroup
	mutex      sync.Mutex
	caldeirao  int = 0 // caldeirao
	avisaCheio     = make(chan int)
	avisaVazio     = make(chan int)
)

func cozinheiro() {
	for {
		<-avisaVazio
		caldeirao += addPorçoes
		fmt.Println("Caldeirão cheio")
		avisaCheio <- 1
		fmt.Println("Cozinheiro Dormindo")
	}
}

func pessoa(id int) {
	for {
		mutex.Lock()
		if caldeirao == 0 {
			avisaVazio <- 1
			<-avisaCheio
		}
		caldeirao -= 1
		fmt.Println("Pessoa ", id, " encheu o prato")
		mutex.Unlock()

		comer(id)
	}
}

func comer(id int) {
	time.Sleep(100 * time.Millisecond)
	//fmt.Println("Pessoa ", id, " comeu")
}

func main() {

	wg.Add(1)
	go cozinheiro()

	for i := 0; i < numPessoas; i++ {
		go pessoa(i)
	}

	wg.Wait()
}
