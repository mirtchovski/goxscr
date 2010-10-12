// Palette useful for testing xscr.*Cmap
package main

import (
	"exp/draw"
	"image"

	"./xscr"
)

func palette(screen draw.Image) {
	const subdivision = 16

	var sw, sh, gw, gh int
	var ncolors = 256
	var col = 0

	colors := xscr.SmoothRandomCmap(ncolors)

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
}

func main() {

	xscr.Init(palette, 1e12)
	xscr.Run()
}
