package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

var MULTI = true

func main() {
	if MULTI {
		cores := runtime.NumCPU()
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
		actual_main(0, 1, os.Args)
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

var NUMBERS = []int64{73490, 70760, 37401, 2336, 17842, 96998, 52582, 29535, 73234, 79890, 3701, 27748, 99582, 87670, 29294, 5959, 79684, 20625, 35237, 94387, 3277, 6154, 74222, 59031, 67957, 77899, 59759, 68265, 76286, 91649, 434, 41238, 76650, 99484, 97021, 20275, 34814, 96385, 9555, 53314, 65808, 32339, 74984, 67499, 90586, 12380, 64975, 84346, 87131, 12446, 46062, 76368, 40129, 82238, 60367, 36958, 43041, 17776, 67768, 94943, 82284, 34872, 54147, 45757, 85319, 23191, 8065, 37207, 36583, 879, 92206, 69356, 45591, 20676, 45346, 94936, 92648, 66213, 43495, 78014, 37555, 24778, 90825, 90635, 58522, 16317, 20898, 81731, 49322, 37625, 38405, 4516, 99065, 3194, 12615, 65435, 19621, 34985, 19434, 5922, 8274, 9324, 99504, 52215, 66474, 131, 33231, 83555, 24543, 4305, 88154, 80181, 47840, 18891, 2197, 50466, 28194, 50961, 48902, 34937, 96172, 62534, 57747, 39759, 7610, 30565, 5296, 3957, 35569, 77666, 29042, 15645, 35587, 39707, 98138, 22572, 80923, 50296, 18619, 47588, 50647, 40289, 9907, 82888, 90531, 45316, 27311, 44708, 37616, 78387, 57733, 50999, 90574, 4981, 26764, 29208, 72028, 28186, 63896, 54632, 1540, 45370, 47145, 15949, 79508, 44132, 1500, 45409, 98709, 29711, 77230, 92750, 38581, 46834, 18033, 76842, 62332, 21290, 78161, 4880, 28627, 68510, 97592, 81525, 9922, 800, 5913, 15452, 46388, 34983, 57258, 99539, 88966, 33078, 64015, 58132, 17628, 34090, 70397, 90632, 96167, 45090, 57620, 2137, 94875, 8704, 86784, 65178, 15487, 77314, 76977, 20636, 67379, 8461, 47277, 9867, 74302, 8880, 80645, 12135, 64921, 81295, 70941, 88641, 77192, 65306, 49453, 26901, 61064, 69014, 43941, 23937, 24770, 12373, 52920, 82238, 72948, 92698, 44094, 18544, 2649, 69329, 9153, 14033, 23215, 55892, 92330, 42686, 62940, 25686, 13622, 30404, 66909, 54514, 96127}
var RESULT = atomic.Int64{}
var BARRIER Barrier

type Barrier struct {
	Flag atomic.Bool
	Wg   sync.WaitGroup
}

func (b *Barrier) Wait(numCores int) {
	if !b.Flag.Load() {
		b.Flag.Store(true)
		b.Wg.Add(numCores)
	}
	b.Wg.Done()
	b.Wg.Wait()
	if b.Flag.Load() {
		b.Flag.Store(false)
	}
}

func actual_main(core, numCores int, args []string) {
	sum := int64(0)
	remainder := len(NUMBERS) % numCores
	sliceSize := len(NUMBERS) / numCores

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
	for _, num := range NUMBERS[startIdx:endIdx] {
		sum += int64(num)
	}

	RESULT.Add(sum)
	BARRIER.Wait(numCores)
	if core == 0 {
		corePrintlf(core, numCores, "sum result = %v\n", RESULT.Load())
	}
}

func corePrintlf(core, numCores int, format string, a ...any) {
	fmt.Printf("%d/%d -> %s\n", core, numCores, fmt.Sprintf(format, a...))
}
