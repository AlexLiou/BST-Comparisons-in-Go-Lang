package main

import (
	"fmt"
	"sync"
)

type ConcurrentBuffer struct {
	data     []interface{}
	maxSize  int
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

func NewConcurrentBuffer(maxSize int) *ConcurrentBuffer {
	buffer := &ConcurrentBuffer{
		data:    make([]interface{}, 0, maxSize),
		maxSize: maxSize,
	}
	buffer.notFull = sync.NewCond(&buffer.mu)
	buffer.notEmpty = sync.NewCond(&buffer.mu)
	return buffer
}

func (b *ConcurrentBuffer) Insert(data interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.data) == b.maxSize {
		b.notFull.Wait()
	}
	b.data = append(b.data, data)
	b.notEmpty.Signal()
}

func (b *ConcurrentBuffer) Remove() interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	for len(b.data) == 0 {
		b.notEmpty.Wait()
	}
	data := b.data[0]
	b.data = b.data[1:]
	b.notFull.Signal()
	return data
}

func (b *ConcurrentBuffer) IsEmpty() bool {
    b.mu.Lock()
    defer b.mu.Unlock()
    return len(b.data) == 0
}

func main() {
	buffer := NewConcurrentBuffer(10)

	// Number of threads
	n := 5

	var wg sync.WaitGroup
    exit := make(chan bool)

    // Create n remove goroutines
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			var step int = 10*10 / n
            var start int = i*step
            var end int = i*step + step
            if i == n - 1 {
                end = (10*10)
            }
			defer wg.Done()
			for i := start; i < end; i++ {
                data := buffer.Remove()
                fmt.Printf("Removed by goroutine %d: %d\n", id, data)
			}
		}(i)
	}

    go func() {
        wg.Add(1)
        defer wg.Done()
        for i := 0; i < 10; i++ {
            for j := 0; j < 10; j++ {
                data := i*j
                buffer.Insert(data)
                fmt.Printf("Inserted %d\n", data)
            }
        }	
        close(exit)
    }()

    wg.Wait()
}
