package raytracer

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"
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
	Workers       int
	totalJobs     int
	remainingJobs int
	startTime     time.Time
	mu            sync.Mutex
	Jobs          chan WorkerJob
	Wg            *sync.WaitGroup
	World         *World
}

func NewWorkerPool(workers int, world *World) *WorkerPool {
	return &WorkerPool{
		Workers: workers,
		Jobs:    make(chan WorkerJob),
		World:   world,
		Wg:      &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Start(totalJobs int, img *image.RGBA, compute ComputeFunc) {
	wp.totalJobs = totalJobs
	wp.remainingJobs = totalJobs
	wp.startTime = time.Now()
	for i := 0; i < wp.Workers; i++ {
		wp.Wg.Add(1)
		go wp.worker(img, compute)
	}
}

func (wp *WorkerPool) worker(img *image.RGBA, compute ComputeFunc) {
	defer wp.Wg.Done()
	for job := range wp.Jobs {
		start := time.Now()
		fmt.Printf("Worker picked up job %d\n", job.Chunk)
		for y := job.YStart; y < job.YEnd; y++ {
			for x := job.XStart; x < job.XEnd; x++ {
				color := compute(x, y, wp.World)
				img.Set(x, y, color)
			}
		}
		wp.mu.Lock()
		wp.remainingJobs--
		wp.mu.Unlock()
		total := time.Since(start)
		elapsed := time.Since(wp.startTime)
		fmt.Printf("Worker finished job %d of %d in %s (remaining %d - time so far %s)\n", job.Chunk, wp.totalJobs, formatDuration(total), wp.remainingJobs, elapsed.String())
	}
}

func (wp *WorkerPool) Wait() {
	close(wp.Jobs)
	wp.Wg.Wait()
}

func formatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second
	d -= s * time.Second

	ms := d / time.Millisecond
	d -= ms * time.Millisecond

	us := d / time.Microsecond

	result := ""
	if h > 0 {
		result += fmt.Sprintf("%dh", h)
	}
	if m > 0 {
		result += fmt.Sprintf("%dm", m)
	}
	if s > 0 {
		result += fmt.Sprintf("%ds", s)
	}
	if ms > 0 {
		result += fmt.Sprintf("%dms", ms)
	}
	if us > 0 && result == "" { // only show microseconds if it's the only value
		result += fmt.Sprintf("%dµs", us)
	}
	return result
}
