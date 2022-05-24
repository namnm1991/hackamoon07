// Package feed provide data
package feed

import (
	"fmt"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
)

type Feed struct {
	mutex sync.Mutex
	data  []float64
	done  chan bool
}

func NewFeed() *Feed {
	f := &Feed{
		mutex: sync.Mutex{},
		data:  fetchData(),
		done:  make(chan bool),
	}

	go f.update()

	return f
}

func (f *Feed) FetchData() []float64 {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.data
}

func (f *Feed) update() {
	t := time.NewTicker(5 * time.Second).C
	for {
		select {
		case <-t:
			f.mutex.Lock()
			f.data = fetchData()
			f.mutex.Unlock()
		case <-f.done:
			return
		}
	}
}

func (f *Feed) Close() {
	f.done <- true
}

func fetchData() []float64 {
	// data := []float64{-1000, 1, 3, 4, 4, 6, 6, 6, 6, 7, 8, 15, 18, 100}
	// return data

	// data := []float64{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}

	// generate random
	// var data []float64
	// for i := 0; i < 100; i++ {

	// 	data = append(data, float64(rand.Intn(100)))
	// }
	// return data

	// data := stats.NormBoxMullerRvs(10, 5, 20)
	// fmt.Println(data)
	// return data

	data := stats.NormBoxMullerRvs(10, 50, 10)
	fmt.Println(data)
	return data
}
