package BST

import (
    "testing"
    "fmt"
)

var printMap = false

func TestSimpleSequentialHash(t *testing.T) {
    inputPath := "input/simple.txt"
    
    Trees := constructTreeFromFile(inputPath)
    computeAndPrintHashSequential(Trees)
}

func TestCoarseSequentialHash(t *testing.T) {
    inputPath := "input/coarse.txt"
    
    Trees := constructTreeFromFile(inputPath)
    computeAndPrintHashSequential(Trees)
}

func TestFineSequentialHash(t *testing.T) {
    inputPath := "input/fine.txt"
    
    Trees := constructTreeFromFile(inputPath)
    computeAndPrintHashSequential(Trees)
}

// hash-workers = i
func TestSimpleParallelHash(t *testing.T) {
    inputPath:= "input/simple.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        computeAndPrintHashParallel(Trees, n)
    }
}

// hash-workers = i
func TestCoarseParallelHash(t *testing.T) {
    inputPath:= "input/coarse.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        computeAndPrintHashParallel(Trees, n)
    }
}

// hash-workers = i
func TestFineParallelHash(t *testing.T) {
    inputPath:= "input/fine.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        computeAndPrintHashParallel(Trees, n)
    }
}

func TestSimpleSequential(t *testing.T) {
    inputPath := "input/simple.txt"
    
    Trees := constructTreeFromFile(inputPath)
    hashmap := computeAndPrintHashmapSequential(Trees, printMap)
    computeAndPrintUniqueTreesSequential(hashmap, printMap)
}

func TestCoarseSequential(t *testing.T) {
    inputPath := "input/coarse.txt"

    Trees := constructTreeFromFile(inputPath)
    hashmap := computeAndPrintHashmapSequential(Trees, printMap)
    computeAndPrintUniqueTreesSequential(hashmap, printMap)
}

func TestFineSequential(t *testing.T) {
    inputPath := "input/fine.txt"

    Trees := constructTreeFromFile(inputPath)
    hashmap := computeAndPrintHashmapSequential(Trees, printMap)
    computeAndPrintUniqueTreesSequential(hashmap, printMap)
}

// hash-workers = i, data-workers = 1, comp-workers = i
func TestSimpleParallelManagerHashWorkers(t *testing.T) {
    inputPath:= "input/simple.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        hashmap := computeAndPrintHashmapParallel(Trees, n, 1, printMap)
        computeAndPrintUniqueTreesAdjacency(hashmap, n, printMap)
        computeAndPrintUniqueTreesParallelLock(hashmap, n, printMap)
    }
}

// hash-workers = i, data-workers = i, comp-workers = i
func TestSimpleParallelLock(t *testing.T) {
    inputPath:= "input/simple.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        hashmap := computeAndPrintHashmapParallel(Trees, n, n, printMap)
        fmt.Println("Adjaceny Implementation: ")
        computeAndPrintUniqueTreesAdjacency(hashmap, n, printMap)
        fmt.Println("Hashmap Implementation: ")
        computeAndPrintUniqueTreesParallelLock(hashmap, n, printMap)
    }
}

// hash-workers = i, data-workers = 1, comp-workers = i
func TestCoarseParallelManagerHashWorkers(t *testing.T) {
    inputPath:= "input/coarse.txt"

    Trees := constructTreeFromFile(inputPath)
    // numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    numThreads := [1]int{4}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        hashmap := computeAndPrintHashmapParallel(Trees, n, 1, printMap)
        computeAndPrintUniqueTreesAdjacency(hashmap, n, printMap)
        computeAndPrintUniqueTreesParallelLock(hashmap, n, printMap)
    }
}

// hash-workers = i, data-workers = i, comp-workers = i
func TestCoarseParallelLock(t *testing.T) {
    inputPath:= "input/coarse.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        hashmap := computeAndPrintHashmapParallel(Trees, n, n, printMap)
        fmt.Println("Adjaceny Implementation: ")
        computeAndPrintUniqueTreesAdjacency(hashmap, n, printMap)
        fmt.Println("Hashmap Implementation: ")
        computeAndPrintUniqueTreesParallelLock(hashmap, n, printMap)
    }
}

// hash-workers = i, data-workers = 1, comp-workers = i
func TestFineParallelManagerHashWorkers(t *testing.T) {
    inputPath:= "input/fine.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        hashmap := computeAndPrintHashmapParallel(Trees, n, 1, printMap)
        computeAndPrintUniqueTreesParallelLock(hashmap, n, printMap)
    }
}

// This hits an index out of bound error somehow on adjacency tree implmentation
// hash-workers = i, data-workers = i, comp-workers = i
func TestFineParallelLock(t *testing.T) {
    inputPath:= "input/fine.txt"

    Trees := constructTreeFromFile(inputPath)
    numThreads := [5]int{2, 4, 8, 16, len(Trees)}
    for _, n := range numThreads {
        fmt.Println("N = ", n)
        hashmap := computeAndPrintHashmapParallel(Trees, n, n, printMap)
        // fmt.Println("Adjaceny Implementation: ")
        // computeAndPrintUniqueTreesAdjacency(hashmap, n, printMap)
        fmt.Println("Hashmap Implementation: ")
        computeAndPrintUniqueTreesParallelLock(hashmap, n, printMap)
    }
}