package main

import (
	"log"
	"runtime"
	"time"

	"github.com/chrisprobst/alloc"
)

// Execute within this directory with: GODEBUG=gctrace=1 go run .
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
		runtime.GC()
		time.Sleep(time.Second)
	}
}
