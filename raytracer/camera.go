package raytracer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"raytracer/math3"
	"runtime"
)

type Camera struct {
	Width            int
	Height           int
	AspectRatio      float64
	SamplesPerPixel  int
	MaxDepth         int
	VFov             float64
	PixelSampleScale float64
	DefocusAngle     float64
	FocusDist        float64
	Center           math3.Vec3
	LookFrom         math3.Vec3
	LookAt           math3.Vec3
	VUp              math3.Vec3
	PixelDeltaU      math3.Vec3
	PixelDeltaV      math3.Vec3
	Pixel00Loc       math3.Vec3
	DefocusDiskU     math3.Vec3
	DefocusDiskV     math3.Vec3
}

type CameraParams struct {
	Width           int
	AspectRatio     float64
	SamplesPerPixel int
	MaxDepth        int
	VFov            float64
	LookFrom        math3.Vec3
	LookAt          math3.Vec3
	DefocusAngle    float64
	FocusDist       float64
}

func NewCamera(params CameraParams) *Camera {
	cam := &Camera{
		Width:            params.Width,
		Height:           int(math.Floor(float64(params.Width) / params.AspectRatio)),
		AspectRatio:      params.AspectRatio,
		SamplesPerPixel:  params.SamplesPerPixel,
		MaxDepth:         params.MaxDepth,
		PixelSampleScale: 1 / float64(params.SamplesPerPixel),
		VFov:             params.VFov,
		DefocusAngle:     params.DefocusAngle,
		FocusDist:        params.FocusDist,
		LookFrom:         params.LookFrom,
		LookAt:           params.LookAt,
		VUp:              math3.Vec3{0, 1, 0},
	}

	cam.Center = cam.LookFrom
	theta := math3.Deg2Rad(cam.VFov)
	h := math.Tan(theta / 2.0)
	viewportHeight := 2 * h * cam.FocusDist
	viewportWidth := viewportHeight * (float64(cam.Width) / float64(cam.Height))
	w := cam.LookFrom.Sub(cam.LookAt).Normalize()
	u := math3.Cross(cam.VUp, w).Normalize()
	v := math3.Cross(w, u)
	viewportU := u.Scale(viewportWidth)
	viewportV := v.Scale(-viewportHeight)
	cam.PixelDeltaU = viewportU.Div(float64(cam.Width))
	cam.PixelDeltaV = viewportV.Div(float64(cam.Height))
	viewportUpperLeft := cam.Center.Sub(w.Scale(cam.FocusDist)).Sub(viewportU.Div(2)).Sub(viewportV.Div(2))
	cam.Pixel00Loc = viewportUpperLeft.Add(cam.PixelDeltaU.Add(cam.PixelDeltaV).Div(2))
	defocusRadius := cam.FocusDist * math.Tan(math3.Deg2Rad(cam.DefocusAngle/2.0))
	cam.DefocusDiskU = u.Scale(defocusRadius)
	cam.DefocusDiskV = v.Scale(defocusRadius)
	return cam
}

func (cam *Camera) Render(world *World, usePool bool) *image.RGBA {
	if !usePool {
		img := image.NewRGBA(image.Rect(0, 0, cam.Width, cam.Height))
		for y := 0; y < cam.Height; y++ {
			fmt.Printf("Rendering scanline %d\n", y)
			for x := 0; x < cam.Width; x++ {
				img.Set(x, y, cam.RenderPixel(x, y, world))
			}
		}
		return img
	}
	return cam.RenderAsync(world)
}

func (cam *Camera) RenderAsync(world *World) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, cam.Width, cam.Height))
	numWorkers := runtime.NumCPU() - 1
	chunkSize := 32
	wp := NewWorkerPool(numWorkers, world)
	totalChunks := int(math.Ceil(float64(cam.Width)/float64(chunkSize))) * int(math.Ceil(float64(cam.Height)/float64(chunkSize)))
	chunks := make([]WorkerJob, 0, totalChunks)
	chunkID := 0

	for y := 0; y < cam.Height; y += chunkSize {
		for x := 0; x < cam.Width; x += chunkSize {
			chunks = append(chunks, WorkerJob{
				XStart: x,
				YStart: y,
				XEnd:   min(x+chunkSize, cam.Width),
				YEnd:   min(y+chunkSize, cam.Height),
				Chunk:  chunkID,
			})
			chunkID++
		}
	}

	rand.Shuffle(len(chunks), func(i, j int) {
		chunks[i], chunks[j] = chunks[j], chunks[i]
	})

	wp.Start(chunkID, img, cam.RenderPixel)
	fmt.Printf("Total jobs: %d\n", chunkID+1)
	for _, job := range chunks {
		wp.Jobs <- job
	}

	wp.Wait()
	return img
}

func (cam *Camera) RenderPixel(x int, y int, world *World) color.Color {
	pixelColor := math3.Vec3{}
	for sample := 0; sample < cam.SamplesPerPixel; sample++ {
		r := cam.GetRay(x, y)
		pixelColor = pixelColor.Add(cam.RayColor(r, cam.MaxDepth, world))
	}
	return convertPixel(pixelColor.Scale(cam.PixelSampleScale))
}

func (cam *Camera) GetRay(x, y int) math3.Ray {
	offsetX, offsetY := rand.Float64()-0.5, rand.Float64()-0.5
	pixelSample := cam.Pixel00Loc.Add(cam.PixelDeltaU.Scale(float64(x) + offsetX)).Add(cam.PixelDeltaV.Scale(float64(y) + offsetY))
	rayOrigin := cam.DefocusDiskSample()
	if cam.DefocusAngle <= 0 {
		rayOrigin = cam.Center
	}
	rayDirection := pixelSample.Sub(rayOrigin)
	return math3.Ray{Origin: rayOrigin, Direction: rayDirection}
}

func (cam *Camera) DefocusDiskSample() math3.Vec3 {
	p := math3.RandomInUnitDisk()
	return cam.Center.Add(cam.DefocusDiskU.Scale(p.X())).Add(cam.DefocusDiskV.Scale(p.Y()))
}

func (cam *Camera) RayColor(r math3.Ray, depth int, world *World) math3.Vec3 {
	if depth <= 0 {
		return math3.Vec3{0.0, 0.0, 0.0}
	}
	if result, hasHit := world.Hit(r, Interval{Min: 0.001, Max: math.MaxFloat64}); hasHit {
		if attenuation, scattered, ok := result.Material.Scatter(r, result); ok {
			survivalScale, shouldTerminate := cam.ShouldTerminateRay(&attenuation, depth)
			if shouldTerminate {
				return math3.Vec3{0.0, 0.0, 0.0}
			}
			if survivalScale > 0 {
				attenuation = attenuation.Scale(1 / survivalScale)
			}
			return cam.RayColor(scattered, depth-1, world).Multiply(attenuation)
		}
		return math3.Vec3{}
	}
	d := r.Direction.Normalize()
	a := 0.5 * (d.Y() + 1.0)
	return math3.Vec3{1.0, 1.0, 1.0}.Scale(1.0 - a).Add(math3.Vec3{0.5, 0.7, 1.0}.Scale(a))
}

func (cam *Camera) ShouldTerminateRay(attenuation *math3.Vec3, depth int) (float64, bool) {
	energy := attenuation.MaxComponent()
	var survivalProb float64
	// Start using Russian Roulette after a few bounces
	if depth < cam.MaxDepth-2 {
		terminationProb := math.Max(0.0, 1.0-energy)
		if rand.Float64() < terminationProb {
			return 0, true
		}

		survivalProb = 1.0 - terminationProb
	}
	return survivalProb, false
}
