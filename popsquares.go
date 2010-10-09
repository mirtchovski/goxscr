/* Copyright (c) 2003 Levi Burton <donburton@sbcglobal.net>
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
 * Ported to Plan 9 by mirtchov@cpsc.ucalgary.ca, 10/03
 */

// Rewritten in Go by mirtchovski@gmail.com 10/10
package main

import (
	"exp/draw"
	"rand"
	"image"
	"sync"

	"./xscr"
)

type square struct {
	col int
	r   image.Rectangle
}

var sw, sh, gw, gh, nsquares int
var squares []*square
var subdivision int
var twitch = false

var colors []image.RGBAColor
var ncolors = 256

var once sync.Once

func popinit(screen draw.Image) {
	if rand.Intn(2) > 0 {
		twitch = true
	}
	colors = xscr.SmoothRandomCmap(ncolors)

	subdivision = rand.Intn(15) + 10
	sw = screen.Bounds().Dx() / subdivision
	sh = screen.Bounds().Dy() / subdivision
	gw = screen.Bounds().Dx() / sw
	gh = screen.Bounds().Dy() / sh
	nsquares = gw * gh

	squares = make([]*square, nsquares, nsquares)
	for y := 0; y < gh; y++ {
		for x := 0; x < gw; x++ {
			col := rand.Intn(ncolors)
			r := image.Rect(x*sw, y*sh, (x+1)*sw-1, (y+1)*sh-1)
			squares[gw*y+x] = &square{col, r}
		}
	}
}

func popsquares(screen draw.Image) {
	once.Do(func() { popinit(screen) })

	for y := 0; y < gh; y++ {
		for x := 0; x < gw; x++ {
			s := squares[gw*y+x]
			draw.Draw(screen, s.r, image.ColorImage{colors[s.col]}, image.ZP)
			s.col = s.col + 1
			if s.col >= ncolors {
				if twitch { //&& rand.Intn(3) == 0 {
					colors = xscr.SmoothRandomCmap(ncolors)
				}
				s.col = rand.Intn(ncolors)
			}
		}
	}
	xscr.Flush()
}

func main() {

	xscr.Init(popsquares, 10e6)
	xscr.Run()
}
