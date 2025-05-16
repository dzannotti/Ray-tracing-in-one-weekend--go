package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"raytracer/math3"
	"raytracer/raytracer"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	world := raytracer.World{}

	ground := raytracer.Lambertian{Albedo: math3.Vec3{0.5, 0.5, 0.5}}
	world.Add(raytracer.Sphere{Center: math3.Vec3{0, -1000, -1}, Radius: 1000, Material: ground})

	mat1 := raytracer.Dialectric{RefractionIndex: 1.5}
	mat2 := raytracer.Lambertian{Albedo: math3.Vec3{0.4, 0.2, 0.1}}
	mat3 := raytracer.Metal{Albedo: math3.Vec3{0.7, 0.6, 0.5}, Fuzz: 0.0}

	world.Add(raytracer.Sphere{Center: math3.Vec3{0, 1, 0}, Radius: 1, Material: mat1})
	world.Add(raytracer.Sphere{Center: math3.Vec3{-4, 1, 0}, Radius: 1, Material: mat2})
	world.Add(raytracer.Sphere{Center: math3.Vec3{4, 1, 0}, Radius: 1, Material: mat3})

	camera := raytracer.NewCamera(raytracer.CameraParams{
		Width:           320,
		AspectRatio:     16 / 9.0,
		SamplesPerPixel: 100,
		MaxDepth:        50,
		VFov:            20,
		LookFrom:        math3.Vec3{13, 2, 3},
		LookAt:          math3.Vec3{0, 0, 0},
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
	mf, err := os.Create("mem_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(mf)
	mf.Close()
}
