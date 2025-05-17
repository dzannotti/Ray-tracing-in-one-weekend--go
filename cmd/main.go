package main

import (
	"fmt"
	"image/png"
	"math/rand/v2"
	"os"
	"raytracer/math3"
	"raytracer/raytracer"
	"time"
)

func main() {
	start := time.Now()

	world := raytracer.World{}

	ground := raytracer.Lambertian{Albedo: math3.Vec3{0.5, 0.5, 0.5}}
	world.Add(&raytracer.Sphere{Center: math3.Vec3{0, -1000, -1}, Radius: 1000, Material: ground})

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMaterial := rand.Float64()
			center := math3.Vec3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			if center.Sub(math3.Vec3{4, 0.2, 0}).Length() > 0.9 {
				albedo := math3.Random().Multiply(math3.Random())
				var material raytracer.Material
				switch {
				case chooseMaterial < 0.8:
					material = raytracer.Lambertian{Albedo: albedo}
				case chooseMaterial < 0.95:
					albedo = math3.RandomBetween(0.5, 1)
					fuzz := math3.RandomBetween(0, 0.5)[0]
					material = raytracer.Metal{Albedo: albedo, Fuzz: fuzz}
				default:
					material = raytracer.Dialectric{RefractionIndex: 1.5}
				}
				world.Add(&raytracer.Sphere{Center: center, Radius: 0.2, Material: material})
			}
		}
	}

	mat1 := raytracer.Dialectric{RefractionIndex: 1.5}
	mat2 := raytracer.Lambertian{Albedo: math3.Vec3{0.4, 0.2, 0.1}}
	mat3 := raytracer.Metal{Albedo: math3.Vec3{0.7, 0.6, 0.5}, Fuzz: 0.0}

	world.Add(&raytracer.Sphere{Center: math3.Vec3{0, 1, 0}, Radius: 1, Material: mat1})
	world.Add(&raytracer.Sphere{Center: math3.Vec3{-4, 1, 0}, Radius: 1, Material: mat2})
	world.Add(&raytracer.Sphere{Center: math3.Vec3{4, 1, 0}, Radius: 1, Material: mat3})

	world.Prepare()
	camera := raytracer.NewCamera(raytracer.CameraParams{
		Width:           1200,
		AspectRatio:     16 / 9.0,
		SamplesPerPixel: 500,
		MaxDepth:        50,
		VFov:            20,
		LookFrom:        math3.Vec3{13, 2, 3},
		LookAt:          math3.Vec3{0, 0, 0},
		DefocusAngle:    0.6,
		FocusDist:       10,
	})
	render := camera.Render(&world, true)
	fmt.Println("denoising....")
	img := raytracer.BilateralFilter(render, 3.0, 0.2)
	file, err := os.Create("final-screenshot.png")
	if err != nil {
		panic("Could not open final-screenshot.png")
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic("Could not write to file")
	}
	total := time.Since(start)
	fmt.Printf("Took: %s\n", total.String())
}
