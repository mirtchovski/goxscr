// QuasiCrystals
// Following the Limbo implementation by Jeff Sickel
// Translated to go by Andrey Mirtchovski, mirtchovski@gmail.com
package main

import (
	"exp/gui/x11"
	"exp/gui"
	"flag"
	"fmt"
	"image/draw"
	"image"
	"image/color"
	"math"
	"os"
	"rand"
	"runtime"
	"time"
)

const Degree = math.Pi / 180

var palette []image.Image

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

func quasicrystal(size, degree int, ϕ float64, buf []byte) {
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

func emitter(e chan int) {
	for i := 0; ; i++ {
		e <- i
	}
}

func worker(c chan image.Image, e chan int, r image.Rectangle) {
	img := image.NewRGBA(r)
	buf := make([]byte, *size*(*size))
	sz, deg := *size, *degree

	stridex := r.Dx() / sz // how big is each pixel from our crystal
	if stridex == 0 {
		stridex = 1
	}
	stridey := r.Dy() / sz
	if stridey == 0 {
		stridey = 1
	}

	for {
		ϕ := float64(<-e) * (*phi) * Degree
		quasicrystal(sz, deg, ϕ, buf[:])

		for y := 0; y < sz; y++ {
			for x := 0; x < sz; x++ {
				nr := image.Rect(x*stridex, y*stridey, x*stridex+stridex, y*stridey+stridey)
				draw.Draw(img, nr, palette[buf[y*sz+x]], image.ZP, draw.Src)
			}
		}
		c <- img
	}
}

var doneframes int64

func painter(c chan image.Image, w gui.Window) {
	t := time.NewTicker(1e9 / (*frames))
	img := w.Screen()
	r := img.Bounds()

	// workers may complete in any order but they sync on 'c' if
	// they're faster than the drawing here
	for {
		draw.Draw(img, r, <-c, image.ZP, draw.Src)
		w.FlushImage()
		doneframes++
		<-t.C
	}
}

var frames = flag.Int64("f", 30, "max framerate")
var phi = flag.Float64("s", 5, "step phase change")
var size = flag.Int("size", 300, "crystal size")
var scale = flag.Float64("scale", 30, "scale")
var degree = flag.Int("degree", 5, "degree")
var workers = flag.Int("w", 3, "workers")

func init() {
	palette = make([]image.Image, 256)
	for i := 0; i < 256; i++ {
		palette[i] = image.NewUniform(color.Gray{byte(i)})
	}
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*workers + 1)

	rand.Seed(int64(os.Getpid()))

	window, err := x11.NewWindow()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:", err.Error())
		return
	}

	emit := make(chan int, *workers)
	go emitter(emit)

	imgchan := make(chan image.Image, *workers)
	for i := 0; i < *workers; i++ {
		go worker(imgchan, emit, window.Screen().Bounds())
	}
	go painter(imgchan, window)

	begin := time.Seconds()
loop:
	for {
		select {
		case e := <-window.EventChan():
			switch f := e.(type) {
			case gui.MouseEvent:
			case gui.KeyEvent:
				if f.Key == 65307 { // ESC
					fmt.Println("fps: ", doneframes/(time.Seconds()-begin))
					break loop
				}
			case gui.ConfigEvent:
				// nothing for now
			case gui.ErrEvent:
				break loop
			}
		}
	}
}
