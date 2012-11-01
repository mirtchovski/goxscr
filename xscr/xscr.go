// Support functions for xscreensaver ports
package xscr

import (
	"code.google.com/p/x-go-binding/ui"
	"code.google.com/p/x-go-binding/ui/x11"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"time"
)

func Border(dst draw.Image, r image.Rectangle, w int, src image.Image, sp image.Point) {
	i := w
	if i > 0 {
		// inside r
		draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+i), src, sp, draw.Over)                                // top
		draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y+i, r.Min.X+i, r.Max.Y-i), src, sp.Add(image.Pt(0, i)), draw.Over)        // left
		draw.Draw(dst, image.Rect(r.Max.X-i, r.Min.Y+i, r.Max.X, r.Max.Y-i), src, sp.Add(image.Pt(r.Dx()-i, i)), draw.Over) // right
		draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-i, r.Max.X, r.Max.Y), src, sp.Add(image.Pt(0, r.Dy()-i)), draw.Over)     // bottom
		return
	}

	// outside r;
	i = -i
	draw.Draw(dst, image.Rect(r.Min.X-i, r.Min.Y-i, r.Max.X+i, r.Min.Y), src, sp.Add(image.Pt(-i, -i)), draw.Over) // top
	draw.Draw(dst, image.Rect(r.Min.X-i, r.Min.Y, r.Min.X, r.Max.Y), src, sp.Add(image.Pt(-i, 0)), draw.Over)      // left
	draw.Draw(dst, image.Rect(r.Max.X, r.Min.Y, r.Max.X+i, r.Max.Y), src, sp.Add(image.Pt(r.Dx(), 0)), draw.Over)  // right
	draw.Draw(dst, image.Rect(r.Min.X-i, r.Max.Y, r.Max.X+i, r.Max.Y+i), src, sp.Add(image.Pt(-i, 0)), draw.Over)  // bottom
}

// Creates a random colormap
func RandomCmap(ncol int) (cmap []color.RGBA) {
	cmap = make([]color.RGBA, ncol, ncol)

	for i := 0; i < ncol; i++ {
		cmap[i] = color.RGBA{uint8(rand.Intn(0x100)), uint8(rand.Intn(0x100)), uint8(rand.Intn(0x100)), 0xff}
	}
	return
}

// Creates a color interpolation between c1 and c2 over num colors
func Interpolate(c1, c2 color.RGBA, ncol int) (cmap []color.RGBA) {
	cmap = make([]color.RGBA, ncol, ncol)

	for i := 0; i < ncol/2; i++ {
		r := int(c1.R) + 2*i*int(c2.R)/ncol
		g := int(c1.G) + 2*i*int(c2.G)/ncol
		b := int(c1.B) + 2*i*int(c2.B)/ncol
		if r > 0xff {
			r = 0xff
		}
		if g > 0xff {
			g = 0xff
		}
		if b > 0xff {
			b = 0xff
		}

		cmap[i] = color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
	}
	for i := 0; i < ncol/2; i++ {
		r := int(c2.R) + 2*i*int(c1.R)/ncol
		g := int(c2.G) + 2*i*int(c1.G)/ncol
		b := int(c2.B) + 2*i*int(c1.B)/ncol
		if r > 0xff {
			r = 0xff
		}
		if g > 0xff {
			g = 0xff
		}
		if b > 0xff {
			b = 0xff
		}
		cmap[ncol-i-1] = color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
	}
	return
}

// A randomized smooth color palette from a dark color to a brighter one
func SmoothRandomCmap(ncol int) []color.RGBA {
	var c [2]color.RGBA

	c[0].R = uint8(rand.Intn(0xff))
	c[0].G = uint8(rand.Intn(0xff))
	c[0].B = uint8(rand.Intn(0xff))
	c[0].A = 0xff

	if rand.Intn(2) > 0 {
		c[1].R = c[0].R + uint8(rand.Intn(0xff-int(c[0].R))) // make sure we don't overflow
	} else {
		c[1].R = c[0].R - uint8(rand.Intn(int(c[0].R))) // or underflow
	}
	if rand.Intn(2) > 0 {
		c[1].G = c[0].G + uint8(rand.Intn(0xff-int(c[0].G))) // make sure we don't overflow
	} else {
		c[1].G = c[0].G - uint8(rand.Intn(int(c[0].G))) // or underflow
	}
	if rand.Intn(2) > 0 {
		c[1].B = c[0].B + uint8(rand.Intn(0xff-int(c[0].B))) // make sure we don't overflow
	} else {
		c[1].B = c[0].B - uint8(rand.Intn(int(c[0].B))) // or underflow
	}
	c[1].A = 0xff

	return Interpolate(c[0], c[1], ncol)
}

var window ui.Window
var hackchan chan bool
var hackfun func(draw.Image)
var hackdelay time.Duration

func Flush() {
	window.FlushImage()
}

// Run the xscreensaver hack
// TODO(aam): resize
func Run() {
	hackfun(window.Screen())
	c := time.Tick(hackdelay)
loop:
	for {
		select {
		case e := <-window.EventChan():
			switch f := e.(type) {
			case ui.MouseEvent:
			case ui.KeyEvent:
				if f.Key == 65307 { // ESC
					break loop
				}
			case ui.ConfigEvent:
				hackchan <- true
				hackfun(window.Screen())
			case ui.ErrEvent:
				break loop
			}
		case <-c:
			hackfun(window.Screen())
			window.FlushImage()
		}
	}
}

func Init(hack func(draw.Image), delay time.Duration) bool {
	var err error

	window, err = x11.NewWindow()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:", err.Error())
		return false
	}

	rand.Seed(int64(os.Getpid()))

	hackchan = make(chan bool, 0)
	hackfun = hack
	hackdelay = delay

	return true
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

// Brezenham's line drawing algorithm (no width)
func Line(dst draw.Image, src color.Color, p0, p1 image.Point) {
	steep := abs(p1.Y-p0.Y) > abs(p1.X-p0.X)
	if steep {
		p0.X, p0.Y = p0.Y, p0.X
		p1.X, p1.Y = p1.Y, p1.X

	}
	if p0.X > p1.X {
		p0, p1 = p1, p0
	}
	deltax := p1.X - p0.X
	deltay := abs(p1.Y - p0.Y)
	error := deltax / 2
	y := p0.Y

	var ystep int
	if p0.Y < p1.Y {
		ystep = 1
	} else {
		ystep = -1
	}

	for x := p0.X; x < p1.X; x++ {
		if steep {
			dst.Set(y, x, src)
		} else {
			dst.Set(x, y, src)
		}
		error = error - deltay
		if error < 0 {
			y = y + ystep
			error = error + deltax
		}
	}
}
