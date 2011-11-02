/* wander, by Rick Campbell <rick@campbellcentral.org>, 19 December 1998.
 *
 * Permission to use, copy, modify, distribute, and sell this software and its
 * documentation for any purpose is hereby granted without fee, provided that
 * the above copyright notice appear in all copies and that both that
 * copyright notice and this permission notice appear in supporting
 * documentation.  No representations are made about the suitability of this
 * software for any purpose.  It is provided "as is" without express or 
 * implied warranty.
 */

/* 
 * ported to Plan 9 by andrey@lanl.gov, 06/02
 */

// Rewritten in Go by mirtchovski@gmail.com 10/10
package main

import (
	"image/draw"
	"image/color"
	"rand"
	"image"
	"sync"

	"./xscr"
)

var colors []color.RGBA
var ncolors int = 16
var col int

var circles bool = false
var advance int = 1
var density int = 5
var length int = 0
var reset int = 0
var size int = 0
var width, height int
var lastx, lasty int

var once sync.Once

var black *image.Uniform

func wanderinit(screen draw.Image) {
	colors = xscr.SmoothRandomCmap(ncolors)
	col = rand.Intn(ncolors)

	reset = 2500000
	if rand.Intn(2) > 0 {
		circles = true
	}
	width = screen.Bounds().Dx()
	height = screen.Bounds().Dy()

	length = 25000

	lastx = rand.Intn(width)
	lasty = rand.Intn(height)

	black = image.NewUniform(color.RGBA{0, 0, 0, 0xff})
	draw.Draw(screen, screen.Bounds(), black, image.ZP, draw.Over)
}

func wander(screen draw.Image) {
	once.Do(func() { wanderinit(screen) })

	width1 := width - 1
	height1 := height - 1
	col = rand.Intn(ncolors)

	var x, y int

	if rand.Intn(density) > 0 {
		x = lastx
		y = lasty
	} else {
		x = (lastx + width1 + rand.Intn(3)) % width
		y = (lasty + height1 + rand.Intn(3)) % height
	}
	if rand.Intn(length) == 0 {
		if advance == 0 {
			col = rand.Intn(ncolors)
		} else {
			col = (col + advance) % ncolors
		}
	}
	/*if rand.Int() > reset {
		draw.Draw(screen, screen.Bounds(), black, image.ZP)
		col = rand.Intn(ncolors)
		x = rand.Intn(width)
		y = rand.Intn(height)
	}*/
	lastx, lasty = x, y
	screen.Set(x, y, colors[col])
}

func main() {
	xscr.Init(wander, 1e3)
	xscr.Run()
}
