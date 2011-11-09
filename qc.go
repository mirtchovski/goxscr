// QuasiCrystals
// Following the Limbo implementation by Jeff Sickel
// Translated to go by Andrey Mirtchovski, mirtchovski@gmail.com
package main

import (
	"./xscr"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"os"
)
// todo(aam): double-buffer

var offset = 50
var ncolors = 256

const Degree = math.Pi / 180

type point struct {
	x, y float64
}

func pt(x, y int) point {
	denom := float64(*size) - 1
	X := *scale * ((float64(2*x) / denom) - 1)
	Y := *scale * ((float64(2*y) / denom) - 1)
	return point{X, Y}
}

func transform(θ float64, p point) point {
	p.x = p.x*math.Cos(θ) - p.y*math.Sin(θ)
	p.y = p.x*math.Sin(θ) + p.y*math.Cos(θ)
	return p
}

func wave(ϕ, θ float64, p point) float64 {
	return (math.Cos(math.Cos(θ)*p.x+math.Sin(θ)*p.y+ϕ) + 1.0) / 2.0
}

func wave1(ϕ, θ float64, p point) float64 {
	if θ != 0.0 {
		p = transform(θ, p)
	}
	return (math.Cos(math.Cos(ϕ)*p.x+math.Sin(ϕ)*p.y) + 1.0) / 2.0
}

func wave2(ϕ, θ float64, p point) float64 {
	if θ != 0.0 {
		p = transform(θ, p)
	}
	return (math.Cos(ϕ+p.y) + 1.) / 2.0
}

var buf []byte

func quasicrystal(size, degree int, ϕ float64) {
	buf = make([]byte, size*size)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			θ := 0 * Degree
			//θ := atan2(real x, real y) + ϕ;
			//θ := atan2(real x + ϕ, real y + ϕ) + ϕ;
			p := pt(x, y)
			acc := wave(ϕ, θ, p)
			for d := 1; d < degree; d++ {
				θ += 180 * Degree / float64(degree)
				if d%2 == 1 {
					acc += 1 - wave(ϕ, θ, p)
				} else {
					acc += wave(ϕ, θ, p)
				}
			}
			buf[y*size+x] = byte(acc * 255.0)
		}
	}
	return
}

func frame(sz, z, deg, time int, img draw.Image) {
	ϕ := float64(time) * (*phi) * Degree
	quasicrystal(sz, deg, ϕ)

	stridex := img.Bounds().Dx() / sz // how big is each pixel from our crystal
	stridey := img.Bounds().Dy() / sz
	for y := 0; y < sz; y++ {
		for x := 0; x < sz; x++ {
			//img.Set(x, y, image.NewUniform(color.Gray{buf[y*sz+x]}))
			r := image.Rect(x*stridex, y*stridey, x*stridex+stridex, y*stridey+stridey)
			draw.Draw(img, r, image.NewUniform(color.Gray{buf[y*sz+x]}), image.ZP, draw.Over)
		}
	}
	return
}

var t = 0

func hack(img draw.Image) {
	frame(*size, *zoom, *degree, t, img)
	t++
}

var frate = flag.Int64("f", 60, "framerate")
var phi = flag.Float64("s", 5, "step phase change")
var size = flag.Int("size", 300, "crystal size")
var zoom = flag.Int("zoom", 3, "zoom")
var scale = flag.Float64("scale", 30, "scale")
var degree = flag.Int("degree", 5, "degree")

func main() {
	flag.Parse()

	rand.Seed(int64(os.Getpid()))

	*phi = 1 + rand.Float64()*10
	*size = 100 + rand.Intn(300)
	*zoom = 1 + rand.Intn(3)
	*scale = 10 + rand.Float64()*30
	*degree = 1 + rand.Intn(10)

	buf = make([]byte, *size*(*size))

	xscr.Init(hack, 1e9/(*frate))
	xscr.Run()
}

var procs = flag.Int("p", 1, "workers")
