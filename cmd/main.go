package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"raytracer/math3"
	"raytracer/raytracer"
	"sync"
)

const width, height = 1920, 1080
const workers = 8

// processPixel simulates expensive pixel processing
func processPixel(x, y int) color.Color {
	r := uint8(128)
	g := uint8(0)
	b := uint8(64)
	return color.RGBA{r, g, b, 255}
}

func worker(id int, jobs <-chan [2]int, output chan<- [3]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for coord := range jobs {
		x, y := coord[0], coord[1]
		col := processPixel(x, y)
		output <- [3]interface{}{x, y, col}
	}
}

func main() {
	/*img := image.NewRGBA(image.Rect(0, 0, width, height))
	jobs := make(chan [2]int, workers)
	output := make(chan [3]interface{}, workers)
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go worker(w, jobs, output, &wg)
	}

	// Writer goroutine
	go func() {
		for result := range output {
			x := result[0].(int)
			y := result[1].(int)
			col := result[2].(color.Color)
			img.Set(x, y, col)
		}
	}()

	// Distribute work
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			jobs <- [2]int{x, y}
		}
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(output)*/
	world := raytracer.World{}

	ground := raytracer.Lambertian{Albedo: math3.Vec3{X: 0.5, Y: 0.5, Z: 0.5}}
	world.Add(raytracer.Sphere{Center: math3.Vec3{X: 0, Y: -1000, Z: -1}, Radius: 1000, Material: ground})

	mat1 := raytracer.Lambertian{Albedo: math3.Vec3{X: 0.93, Y: 0.74, Z: 0.74}}
	mat2 := raytracer.Lambertian{Albedo: math3.Vec3{X: 0.89, Y: 0.78, Z: 0.56}}
	mat3 := raytracer.Lambertian{Albedo: math3.Vec3{X: 0.65, Y: 0.81, Z: 0.53}}

	world.Add(raytracer.Sphere{Center: math3.Vec3{X: 0, Y: 1, Z: 0}, Radius: 1, Material: mat1})
	world.Add(raytracer.Sphere{Center: math3.Vec3{X: -4, Y: 1, Z: 0}, Radius: 1, Material: mat2})
	world.Add(raytracer.Sphere{Center: math3.Vec3{X: 4, Y: 1, Z: 0}, Radius: 1, Material: mat3})

	camera := raytracer.NewCamera(raytracer.CameraParams{
		Width:           320,
		AspectRatio:     16 / 9.0,
		SamplesPerPixel: 100,
		MaxDepth:        50,
		VFov:            20,
		LookFrom:        math3.Vec3{X: 13, Y: 2, Z: 3},
		LookAt:          math3.Vec3{X: 0, Y: 0, Z: 0},
		DefocusAngle:    0.6,
		FocusDist:       10,
	})
	img := camera.Render(world)
	file, err := os.Create("final-screenshot.png")
	if err != nil {
		panic("Could not open final-screenshot.png")
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic("Could not write to file")
	}

	fmt.Println("Image saved as final-screenshot.png")
}
