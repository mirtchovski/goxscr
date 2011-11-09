/* The Spiral Generator, Copyright (c) 2000 
 * by Rohit Singh <rohit_singh@hotmail.com>
 * 
 * Contains code from / To be used with:
 * xscreensaver, Copyright (c) 1992, 1995, 1996, 1997
 * Jamie Zawinski <jwz@jwz.org>
 *
 * Permission to use, copy, modify, distribute, and sell this software and its
 * documentation for any purpose is hereby granted without fee, provided that
 * the above copyright notices appear in all copies and that both that
 * copyright notices and this permission notice appear in supporting
 * documentation.  No representations are made about the suitability of this
 * software for any purpose.  It is provided "as is" without express or 
 * implied warranty.
 *
 * Modified (Dec 2001) by Matthew Strait <straitm@mathcs.carleton.edu>
 * Added -subdelay and -alwaysfinish
 * Prevented redrawing over existing lines
 */

/* 
 * ported to Plan 9 by andrey@lanl.gov, 06/02
 */

// Rewritten in go by mirtchovski@gmail.com 10/10
package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"./xscr"
)

var black *image.Uniform
var colors []color.RGBA
var ncolors int = 16
var clr int

func do(screen draw.Image, radius1, radius2, d, clr int) {
	var width, height, xmid, ymid, x1, y1, x2, y2 int
	var theta, delta int

	var firstx, firsty, tmpx, tmpy int
	firstx, firsty = 0, 0

	width = screen.Bounds().Dx()
	height = screen.Bounds().Dy()
	delta = 1
	xmid = width / 2
	ymid = height / 2

	x1 = xmid + radius1 - radius2 + d
	y1 = ymid

loop:
	for theta = 1; theta < 360*1000; theta++ {
		tmpx = xmid + int(float64(radius1-radius2)*math.Cos((float64(theta)*math.Pi)/180)) + int(float64(d)*math.Cos((float64(radius1*theta-delta)/float64(radius2))*math.Pi/180))

		tmpy = ymid + int(float64(radius1-radius2)*math.Sin((float64(theta)*math.Pi)/180)) + int(float64(d)*math.Sin(((float64((radius1*theta)-delta)/float64(radius2))*math.Pi/180)))

		/*makes integers from the calculated values to do the drawing*/
		x2 = tmpx
		y2 = tmpy

		/*stores the first values for later reference*/
		if theta == 1 {
			firstx = tmpx
			firsty = tmpy
		}
		npt1 := screen.Bounds().Min.Add(image.Pt(x1, y1))
		npt2 := screen.Bounds().Min.Add(image.Pt(x2, y2))

		xscr.Line(screen, colors[clr], npt1, npt2)

		xscr.Flush()
		time.Sleep(1e3)

		x1 = x2
		y1 = y2

		if tmpx == firstx && tmpy == firsty && theta != 1 {
			break loop
		}
	}
}

func min(f, s int) int {
	if f < s {
		return f
	}
	return s
}

func spirograph(screen draw.Image) {
	var width, height, radius, radius1, radius2 int
	var divisor float64
	var distance int

	width = screen.Bounds().Dx()
	height = screen.Bounds().Dy()

	radius = min(width, height) / 2

	draw.Draw(screen, screen.Bounds(), black, image.ZP, draw.Over)

	divisor = ((rand.Float64()*3.0 + 1) * float64((((rand.Intn(2) & 1) * 2) - 1)))

	radius1 = radius
	radius2 = int(float64(radius)/divisor) + 5
	distance = 100 + rand.Intn(200)

	clr = rand.Intn(ncolors)
	do(screen, radius1, -radius2, distance, clr)
	clr = rand.Intn(ncolors)
	do(screen, radius1, radius2, distance, clr)
}

func main() {
	black = image.NewUniform(color.RGBA{0, 0, 0, 0xff})
	colors = xscr.RandomCmap(ncolors)

	xscr.Init(spirograph, 1e9)
	xscr.Run()
}
