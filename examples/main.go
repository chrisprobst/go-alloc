package main

import (
	"log"
	"runtime"
	"time"

	"github.com/chrisprobst/alloc"
)

func main() {
	go func() {
		N := uint32(1024 * 1024 * 1024)
		AllocOwner := uint8(0)

		type Slice = alloc.Slice[uint32, uint8]
		allocator := alloc.RegisterGlobal[uint32, uint8, byte](N, AllocOwner)

		datas := make([]Slice, N+1)

		for {
			log.Print("Allocating...")
			s := "Hello"
			for i := range N {
				datas[i] = allocator.StoreSlice([]byte(s[:2]))

			}
			log.Print("Done...")

			var m runtime.MemStats
			for {
				runtime.ReadMemStats(&m)
				log.Printf("Current live objects %v but allocated objects %v", m.HeapObjects, allocator.Len())

				agg := 0
				for _, d := range datas[:len(datas)-1] {
					agg += int(allocator.DerefSlice(d)[0])
				}
				log.Printf("Agg: %v", agg)
			}

		}
	}()

	for {
		// runtime.GC()
		time.Sleep(time.Second)
	}
}

// func main() {
// 	s := strconv.Itoa(42731643872)
// 	N := 1024 * 1024 * 1024
// 	alloc := NewAlloc[Data](uint(N + 1))
// 	datas := make([]Ptr, 0, N+1)
// 	for {
// 		alloc.Reset()
// 		datas = datas[:0]

// 		log.Print("Allocating...")
// 		for range N {
// 			datas = append(datas, alloc.Alloc())
// 		}

// 		log.Print("Filling...")
// 		for i, p := range datas {
// 			d := alloc.Get(p)
// 			d.age = i * 10
// 			d.salary = float32(i * 1000)
// 			copy(d.name[:], s)
// 		}

// 		log.Print("Freeing...")
// 		for _, p := range datas {
// 			alloc.Free(p)
// 		}

// 		log.Print("Done")

// 		runtime.GC()
// 	}
// }
