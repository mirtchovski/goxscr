// Support functions for xscreensaver ports
package xscr

//TODO(aam): benchmarks? how do we allow gotest to benchmarks various hacks?
import (
	"exp/draw/x11"
	"exp/draw"
	"image"
	"time"
	"rand"
	"fmt"
	"os"
)
// Creates a random colormap
func RandomCmap(ncol int) (cmap []image.RGBAColor) {
	cmap = make([]image.RGBAColor, ncol, ncol)

	for i := 0; i < ncol; i++ {
		cmap[i] = image.RGBAColor{uint8(rand.Intn(0x100)), uint8(rand.Intn(0x100)), uint8(rand.Intn(0x100)), 0xff}
	}
	return
}

// Creates a color interpolation between c1 and c2 over num colors
func Interpolate(c1, c2 image.RGBAColor, ncol int) (cmap []image.RGBAColor) {
	cmap = make([]image.RGBAColor, ncol, ncol)

	for i := 0; i < ncol/2; i++ {
		r := int(c1.R) + 2*i*int(c2.R)/ncol
		g := int(c1.G) + 2*i*int(c2.G)/ncol
		b := int(c1.B) + 2*i*int(c2.B)/ncol
		if r > 0xff { r = 0xff }
		if g > 0xff { g = 0xff }
		if b > 0xff { b = 0xff }

		cmap[i] = image.RGBAColor{uint8(r), uint8(g), uint8(b), 0xff}
	}
	for i := 0; i < ncol/2; i++ {
		r := int(c2.R) + 2*i*int(c1.R)/ncol
		g := int(c2.G) + 2*i*int(c1.G)/ncol
		b := int(c2.B) + 2*i*int(c1.B)/ncol
		if r > 0xff { r = 0xff }
		if g > 0xff { g = 0xff }
		if b > 0xff { b = 0xff }
		cmap[ncol-i-1] = image.RGBAColor{uint8(r), uint8(g), uint8(b), 0xff}
	}
	return
}

// A randomized smooth color palette from a dark color to a brighter one
func SmoothRandomCmap(ncol int) []image.RGBAColor {
	var c [2]image.RGBAColor

	c[0].R = uint8(rand.Intn(128))
	c[0].G = uint8(rand.Intn(128))
	c[0].B = uint8(rand.Intn(128))
	c[0].A = 0xff

	c[1].R = uint8(rand.Intn(128)) + c[0].R 
	c[1].G = uint8(rand.Intn(128)) + c[0].G
	c[1].B = uint8(rand.Intn(128)) + c[0].B
	c[1].A = 0xff

	return Interpolate(c[0], c[1], ncol)
}

var window draw.Window
var hackchan chan bool
var hackfun func(draw.Image)
var hackdelay int64

func hack(s draw.Image) {
	for {
		if _, ok := <-hackchan; ok {
			return
		}
		hackfun(s)
		window.FlushImage()
		time.Sleep(hackdelay)
	}
}

func Flush() {
	window.FlushImage()
}

// Run the xscreensaver hack
// TODO(aam): make sure resize works once draw/x11 starts supporting it
func Run() {
	go hack(window.Screen())
loop:
	for {
		e := <-window.EventChan()
		switch f := e.(type) {
		case draw.MouseEvent:
		case draw.KeyEvent:
			if f.Key == 65307 { // ESC
				break loop
			}
		case draw.ConfigEvent:
			hackchan <- true
			go hack(window.Screen())
		case draw.ErrEvent:
			break loop
		}
	}
}

// Initialize all state
func Init(hack func(draw.Image), delay int64) bool {
	var err os.Error

	window, err = x11.NewWindow()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:", err.String())
		return false
	}

	rand.Seed(int64(os.Getpid()))

	hackchan = make(chan bool, 0)
	hackfun = hack
	hackdelay = delay

	return true
}
