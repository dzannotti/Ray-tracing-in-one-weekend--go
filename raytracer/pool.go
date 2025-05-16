package raytracer

import (
	"fmt"
	"image"
	"image/color"
	"sync"
)

type WorkerJob struct {
	XStart int
	YStart int
	XEnd   int
	YEnd   int
	Chunk  int
}

type ComputeFunc func(x int, y int, world *World) color.Color

type WorkerPool struct {
	Workers int
	Jobs    chan WorkerJob
	Wg      *sync.WaitGroup
	World   *World
}

func NewWorkerPool(workers int, world *World) *WorkerPool {
	return &WorkerPool{
		Workers: workers,
		Jobs:    make(chan WorkerJob),
		World:   world,
		Wg:      &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Start(img *image.RGBA, compute ComputeFunc) {
	for i := 0; i < wp.Workers; i++ {
		wp.Wg.Add(1)
		go wp.worker(img, compute)
	}
}

func (wp *WorkerPool) worker(img *image.RGBA, compute ComputeFunc) {
	defer wp.Wg.Done()
	for job := range wp.Jobs {
		fmt.Printf("Worker picked up chunk %d\n", job.Chunk)
		for y := job.YStart; y < job.YEnd; y++ {
			for x := job.XStart; x < job.XEnd; x++ {
				color := compute(x, y, wp.World)
				img.Set(x, y, color)
			}
		}
		fmt.Printf("Worker finished chunk %d\n", job.Chunk)
	}
}

func (wp *WorkerPool) Wait() {
	close(wp.Jobs)
	wp.Wg.Wait()
}
