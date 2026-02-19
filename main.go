package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var MULTI = true

func main() {
	start := time.Now()
	if MULTI {
		cores := runtime.NumCPU()
		BARRIER = NewBarrier(cores)
		var wg sync.WaitGroup
		wg.Add(cores)
		for i := range cores {
			go func(core, numCores int) {
				defer wg.Done()
				actual_main(core, numCores, os.Args)
			}(i, cores)
		}
		wg.Wait()
	} else {
		BARRIER = NewBarrier(1)
		actual_main(0, 1, os.Args)
	}
	elapsed := time.Since(start)
	fmt.Println("Elapsed:", elapsed)
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

var NUMBERS []string
var FILE_SIZE = atomic.Int64{}
var RESULT = atomic.Int64{}
var BARRIER *Barrier

type Barrier struct {
	barrierId     int
	finishedCount int
	numCores      int
	cond          *sync.Cond
	mutex         sync.Mutex
}

func NewBarrier(numCores int) *Barrier {
	b := &Barrier{
		numCores: numCores,
	}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

func (b *Barrier) Wait() {
	if b.numCores == 1 {
		return
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	id := b.barrierId
	b.finishedCount++

	if b.finishedCount == b.numCores {
		b.barrierId++
		b.finishedCount = 0
		b.cond.Broadcast()
		return
	}

	for id == b.barrierId { //this for is specifically for false-positives
		b.cond.Wait() //wait might return even if there was no broadcast
	}
}

func calculateTaskBounds(core, numCores, taskSize int) (_start, _end int) {
	remainder := taskSize % numCores
	sliceSize := taskSize / numCores

	hasLeftovers := core < remainder
	leftovers := 0
	if hasLeftovers {
		leftovers = core
	} else {
		leftovers = remainder
	}
	startIdx := (core * sliceSize) + leftovers
	endIdx := startIdx + sliceSize
	if hasLeftovers {
		endIdx += 1
	}
	return startIdx, endIdx
}

func actual_main(core, numCores int, args []string) {
	if core == 0 {
		f, err := os.Open("input.txt")
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			NUMBERS = append(NUMBERS, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		//todo: read parts of file with each core
		// info, err := f.Stat()
		// if err != nil {
		// 	panic(err)
		// }

		// FILE_SIZE.Store(info.Size())
		f.Close()
	}
	BARRIER.Wait()

	startIdx, endIdx := calculateTaskBounds(core, numCores, len(NUMBERS))
	sum := int64(0)
	for _, num := range NUMBERS[startIdx:endIdx] {
		val, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			panic(err)
		}
		sum += int64(val)
	}

	RESULT.Add(sum)
	BARRIER.Wait()
	if core == 0 {
		corePrintlf(core, numCores, "sum result = %v\n", RESULT.Load())
	}
}

func corePrintlf(core, numCores int, format string, a ...any) {
	fmt.Printf("%d/%d -> %s\n", core, numCores, fmt.Sprintf(format, a...))
}
