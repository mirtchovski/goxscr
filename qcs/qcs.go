// QuasiCrystals
package main

import (
	"exp/gui"
	"exp/gui/x11"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

const Degree = math.Pi / 180

var palette = make([]image.Image, 256)

type point struct {
	x, y float64
}

var workers = flag.Int("w", 3, "workers")
var frames = flag.Int64("f", 30, "max framerate")
var randomize = flag.Bool("r", false, "randomize size, scale, degree and phi")

var phi = flag.Float64("phi", 5, "step phase change")
var size = flag.Int("size", 300, "crystal size")
var scale = flag.Float64("scale", 30, "scale")
var degree = flag.Int("degree", 5, "degree")

type Work struct {
	e    int        // frame number to compute.
	img  draw.Image // image to write to.
	done chan bool  // send on this channel when work is done.
}

func init() {
	for i := range palette {
		palette[i] = image.NewUniform(color.Gray{byte(i)})
	}
}

func pt(x, y int) point {
	denom := float64(*size) - 1
	X := *scale * ((float64(2*x) / denom) - 1)
	Y := *scale * ((float64(2*y) / denom) - 1)
	return point{X, Y}
}

func transform(θ float64, p point) point {
	sin, cos := math.Sincos(θ)
	p.x = p.x*cos - p.y*sin
	p.y = p.x*sin + p.y*cos
	return p
}

func worker(wc <-chan *Work) {
	buf := make([]byte, *size**size)
	sz := *size

	for w := range wc {
		r := w.img.Bounds()
		dx := r.Dx()
		dy := r.Dy()

		stridex := 1 + dx/sz // how big is each pixel from our crystal
		stridey := 1 + dy/sz

		ϕ := float64(w.e) * (*phi) * Degree
		quasicrystal(sz, *degree, ϕ, buf)

		for y := 0; y < sz; y++ {
			if y*stridey > dy {
				break
			}
			for x := 0; x < sz; x++ {
				if x*stridex > dx {
					break
				}
				nr := image.Rect(x*stridex, y*stridey, x*stridex+stridex, y*stridey+stridey)
				draw.Draw(w.img, nr, palette[buf[y*sz+x]], image.ZP, draw.Src)
			}
		}
		w.done <- true
	}
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*workers + 1)

	rand.Seed(time.Nanoseconds())

	if *randomize {
		*phi = rand.Float64() * 10
		*size = 100 + rand.Intn(200)
		*scale = 25 + rand.Float64()*10
		*degree = 3 + rand.Intn(5)
	}

	window, err := x11.NewWindow()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:", err.Error())
		return
	}
	quit := make(chan chan<- stats)
	go painter(window, quit)

loop:
	for e := range window.EventChan() {
		switch f := e.(type) {
		case gui.MouseEvent:
		case gui.KeyEvent:
			if f.Key == 65307 { // ESC
				break loop
			}
		case gui.ConfigEvent:
			// nothing for now
		case gui.ErrEvent:
			break loop
		}
	}

	c := make(chan stats)
	quit <- c
	st := <-c
	fmt.Printf("fps: %.1f, spf %.0fms, dev %.0fms\n", 1e9/st.mean, st.mean/1e6, st.stddev/1e6)
}

type stats struct {
	mean   float64
	stddev float64
}

func painter(win gui.Window, quit <-chan chan<- stats) {
	ticker := time.NewTicker(1e9 / *frames)
	screen := win.Screen()
	r := screen.Bounds()

	// make more work items than workers so that we
	// can keep a worker busy even when the last frame
	// that it has computed has not yet been retrieved by the
	// painter loop.
	work := make([]Work, *workers*2)
	workChan := make(chan *Work)

	for i := 0; i < *workers; i++ {
		go worker(workChan)
	}

	e := 0
	frames := 0

	now := time.Nanoseconds()
	start := now
	var sumdt2 float64

	// continuously cycle through the array of work items,
	// waiting for each to be done in turn.
	for {
		for i := range work {
			w := &work[i]
			if w.img == nil {
				// If this is the first time we've used a work item, so make the image
				// and the done channel. There's no image calculated yet.
				w.img = image.NewRGBA(screen.Bounds())
				w.done = make(chan bool, 1)
			} else {
				<-w.done
				draw.Draw(screen, r, w.img, image.ZP, draw.Src)
				win.FlushImage()
				frames++
				// wait for the next tick event or to be asked to quit.
				select {
				case t := <-ticker.C:
					dt := t - now
					sumdt2 += float64(dt * dt)
					now = t
				case c := <-quit:
					mean := float64((now - start) / int64(frames))
					c <- stats{
						mean:   mean,
						stddev: math.Sqrt((sumdt2 / float64(frames)) - mean*mean),
					}
					return
				}
			}

			// start the new work item running on any worker that's available.
			w.e = e
			e++
			workChan <- w
		}
	}
}

func wave(ϕ, θ float64, p point) float64 {
	sin, cos := math.Sincos(θ)
	return (math.Cos(cos*p.x+sin*p.y+ϕ) + 1.0) / 2.0
}

func wave1(ϕ, θ float64, p point) float64 {
	if θ != 0.0 {
		p = transform(θ, p)
	}
	sin, cos := math.Sincos(ϕ)
	return (math.Cos(cos*p.x+sin*p.y) + 1.0) / 2.0
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
