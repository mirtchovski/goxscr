// Palette useful for testing xscr.*Cmap
package main

import (
	"exp/draw"
	"image"
	"flag"

	"./xscr"
)

var subdivision int

var sw, sh, gw, gh int
var ncolors = 1024
var col = 0
var colors []image.RGBAColor
var cycle bool

func palette(screen draw.Image) {
	sw = screen.Bounds().Dx() / subdivision
	sh = screen.Bounds().Dy() / subdivision
	gw = screen.Bounds().Dx() / sw
	gh = screen.Bounds().Dy() / sh

	for y := 0; y < gh; y++ {
		for x := 0; x < gw; x++ {
			col = (col + 1) % ncolors
			r := image.Rect(x*sw, y*sh, (x+1)*sw-1, (y+1)*sh-1)
			draw.Draw(screen, r, image.NewColorImage(colors[col]), image.ZP)
		}
	}
	if cycle {
println(cycle)
		col = (col + 1) % ncolors
	}
	xscr.Flush()
}

var size = flag.Int("size", 16, "width of the palette")
var fcycle = flag.Bool("cycle", false, "cycle through colors")

func main() {
	flag.Parse()
	subdivision = *size
	cycle = *fcycle
	ncolors = subdivision * subdivision
	colors = xscr.SmoothRandomCmap(ncolors)
println(cycle, *fcycle)
	xscr.Init(palette, 10e6)
	xscr.Run()
}
