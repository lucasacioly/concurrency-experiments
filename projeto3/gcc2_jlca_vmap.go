package main

import (
	"fmt"
	"sync"
	"time"
)

/*

O Jantar dos Selvagens

Uma variação curiosa do problema dos produtores/consumidores foi proposta
em [Andrews, 1991] com o nome de Jantar dos Selvagens: uma tribo de selvagens está
jantando ao redor de um grande caldeirão contendo N porções de missionário cozido.
Quando um selvagem quer comer, ele se serve de uma porção no caldeirão, a menos
que este esteja vazio. Nesse caso, o selvagem primeiro acorda o cozinheiro da tribo e
espera que ele encha o caldeirão de volta, para então se servir novamente. Após encher
o caldeirão, o cozinheiro volta a dormir.

Restrições de sincronização:

Selvagens não podem se servir ao mesmo tempo (mas podem comer ao mesmo
tempo);

• Selvagens não podem se servir se o caldeirão estiver vazio;
• O cozinheiro só pode encher o caldeirão quando ele estiver vazio.


*/

const (
	numPessoas int = 5
	addPorçoes int = 10 //numero de porções para adicionar
)

var (
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
		fmt.Println("Pessoa ", id, " encheu o prato. Restam  ", caldeirao, "  porções")
		mutex.Unlock()

		comer(id)
	}
}

func comer(id int) {
	time.Sleep(3000 * time.Millisecond)
}

func main() {

	wg.Add(1)
	go cozinheiro()

	for i := 0; i < numPessoas; i++ {
		go pessoa(i)
	}

	wg.Wait()
}
