package raytracer

import (
	"fmt"
	"image"
	"math"
	"raytracer/math3"
)

type Camera struct {
	Width            int
	Height           int
	AspectRatio      float64
	SamplesPerPixel  int
	MaxDepth         float64
	VFov             float64
	PixelSampleScale float64
	DefocusAngle     float64
	FocusDist        float64
	CameraCenter     math3.Vec3
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
	MaxDepth        float64
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
		VUp:              math3.Vec3{X: 0, Y: 1, Z: 0},
	}

	cam.CameraCenter = cam.LookFrom
	theta := math3.Ded2Rad * cam.VFov
	h := math.Tan(theta / 2.0)
	viewportHeight := 2 * h * cam.FocusDist
	viewportWidth := viewportHeight * (float64(cam.Width) / float64(cam.Height))
	w := cam.LookFrom.Sub(cam.LookAt).Normalize()
	u := math3.Cross(cam.VUp, w).Normalize()
	v := math3.Cross(w, u)
	viewportU := u.K(viewportWidth)
	viewportV := v.K(-viewportHeight)
	cam.PixelDeltaU = viewportU.Div(float64(cam.Width))
	cam.PixelDeltaV = viewportV.Div(float64(cam.Height))
	viewportUpperLeft := cam.CameraCenter.Sub(w.K(cam.FocusDist)).Sub(viewportU.Div(2)).Sub(viewportV.Div(2))
	cam.Pixel00Loc = viewportUpperLeft.Add(cam.PixelDeltaU.Add(cam.PixelDeltaV).Div(2))
	defocusRadius := cam.FocusDist * math.Tan(math3.Ded2Rad*(cam.DefocusAngle/2.0))
	cam.DefocusDiskU = u.K(defocusRadius)
	cam.DefocusDiskV = v.K(defocusRadius)
	return cam
}

func (cam *Camera) Render(world World) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, cam.Width, cam.Height))
	for y := 0; y < cam.Height; y++ {
		fmt.Printf("Rendering scanline %d\n", y)
		for x := 0; x < cam.Width; x++ {
			pixelColor := math3.Vec3{}
			for sample := 0; sample < cam.SamplesPerPixel; sample++ {
				r := cam.GetRay(x, y)
				pixelColor = pixelColor.Add(cam.RayColor(&r, cam.MaxDepth, world))
			}
			img.Set(x, y, convertPixel(pixelColor))
		}
	}
	return img
}

func (cam *Camera) GetRay(x int, y int) math3.Ray {
	//offset := math3.RandomBetween(-0.5, 0.5)
	offset = math3.Vec3{}
	pixelSample := cam.Pixel00Loc.Add(cam.PixelDeltaU.K(float64(x) + offset.X)).Add(cam.PixelDeltaV.K(float64(y) + offset.Y))
	rayOrigin := cam.CameraCenter
	if cam.DefocusAngle > 0 {
		rayOrigin = cam.DefocusDiskSample()
	}
	rayDirection := pixelSample.Sub(rayOrigin)
	return math3.Ray{Origin: rayOrigin, Direction: rayDirection}
}

func (cam *Camera) DefocusDiskSample() math3.Vec3 {
	p := math3.RandomInUnitDisk()
	return cam.CameraCenter.Add(cam.DefocusDiskU.K(p.X)).Add(cam.DefocusDiskV.K(p.Y))
}

func (cam *Camera) RayColor(r *math3.Ray, depth float64, world World) math3.Vec3 {
	if depth <= 0 {
		return math3.Vec3{}
	}
	rec := HitRecord{}
	attenuation := &math3.Vec3{}
	scattered := &math3.Ray{}
	hasHit, resultRec := world.Hit(*r, Interval{Min: 0.0001, Max: math.MaxFloat64}, rec)
	if hasHit {
		if resultRec.Material.Scatter(&r, resultRec, &attenuation, &scattered) {
			return cam.RayColor(scattered, depth-1, world).VectorMultiply(*attenuation)
		}
		return math3.Vec3{}
	}
	unitDir := r.Direction.Normalize()
	a := 0.5 * (unitDir.Y + 1)
	return math3.Vec3{X: 1, Y: 1, Z: 1}.K(1 - a).Add(math3.Vec3{X: 0.5, Y: 0.7, Z: 1.0}.K(a))
}
