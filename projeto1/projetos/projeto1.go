package projetos

/*
O programa é um multiplicador de
matrizes quadradas, tendo uma
função para realizar a multiplicação de
forma sequencial e uma função para
realizar a multiplicação de forma
concorrente. Na função concorrente
realizamos a multiplicação de cada
linha em uma thread, sendo o número de
threads regulado pelo tamanho da matriz.
*/

/*
Medias de execução (ms)

SEQUENCIAL

1,954966667
386,9660889
7746,578211

CONCORRENTE

0,5227444444
40,14957778
1097,749622

*/

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var MATRIX_SIZE int = 1000 //PARAMETRO PARA O NUMERO DE THREADS
var SEED1 int64 = 81
var SEED2 int64 = 50
var M1 [][]int
var M2 [][]int
var M_OUT [][]int = make([][]int, MATRIX_SIZE)

func generateMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j] = rand.Intn(100) // Gera um número inteiro aleatório entre 0 e 99
		}
	}
	return matrix
}

func printMatrix(matrix [][]int) {
	for _, row := range matrix {
		fmt.Println(row)
	}
}

func matrixMultSequencial(m1 [][]int, m2 [][]int, mSize int) [][]int {

	mOut := make([][]int, mSize)

	for i := 0; i < mSize; i++ {
		mOut[i] = make([]int, mSize)
		for j := 0; j < mSize; j++ {
			for k := 0; k < mSize; k++ {
				mOut[i][j] += m1[i][k] * m2[k][j]
			}
		}
	}

	return mOut
}

func matrixMultThreads(m1 [][]int, m2 [][]int, mSize int, startTime int64) {
	wg := sync.WaitGroup{}

	for i := 0; i < mSize; i++ {
		M_OUT[i] = make([]int, mSize)
		wg.Add(1)
		go lineMultiply(m1[i], m2, i, mSize, &wg)
	}

	wg.Wait()
	//printMatrix(M_OUT)
	elapsedTime := time.Now().UnixNano() - startTime
	fmt.Printf("%d\n", elapsedTime)

}

func lineMultiply(line1 []int, m2 [][]int, row int, mSize int, wg *sync.WaitGroup) {
	defer wg.Done()
	lineOut := make([]int, mSize)
	for j := 0; j < mSize; j++ {
		for k := 0; k < mSize; k++ {
			lineOut[j] += line1[k] * m2[k][j]
		}
	}
	M_OUT[row] = lineOut
}

func Run1() {

	rand.Seed(SEED1)
	M1 = generateMatrix(MATRIX_SIZE, MATRIX_SIZE)
	//printMatrix(M1)

	rand.Seed(SEED2)
	M2 = generateMatrix(MATRIX_SIZE, MATRIX_SIZE)
	//printMatrix(M2)

	startTime := time.Now().UnixNano()

	//PASSO SEQUENCIAL
	//M_OUT = matrixMultSequencial(M1, M2, MATRIX_SIZE)
	//printMatrix(M_OUT)
	//elapsedTime := time.Now().UnixNano() - startTime
	//fmt.Printf("%d\n", elapsedTime)

	// CONCORRENTE
	matrixMultThreads(M1, M2, MATRIX_SIZE, startTime)

}
