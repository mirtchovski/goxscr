/* xscreensaver, Copyright (c) 1997, 1998, 2001 Jamie Zawinski <jwz@jwz.org>
 *
 * Permission to use, copy, modify, distribute, and sell this software and its
 * documentation for any purpose is hereby granted without fee, provided that
 * the above copyright notice appear in all copies and that both that
 * copyright notice and this permission notice appear in supporting
 * documentation.  No representations are made about the suitability of this
 * software for any purpose.  It is provided "as is" without express or
 * implied warranty.
 *
 * Concept snarfed from Michael D. Bayne in
 * http://www.go2net.com/internet/deep/1997/04/16/body.html
 */

/*
 * Ported to Plan 9 by mirtchovski@gmail.com, 09/03
 */

// Rewritten in Go by mirtchovski@gmail.com, 10/10
package main

import (
	"image/draw"

	"code.google.com/p/goxscr/xscr"
	"math/rand"
)

var offset = 50
var ncolors = 256

func moire(screen draw.Image) {
	factor := rand.Intn(offset) + 1
	r := screen.Bounds()

	colors := xscr.SmoothRandomCmap(ncolors)
	ncolors = len(colors)
	xo := rand.Intn(r.Dx()/2) + r.Dx()/2
	yo := rand.Intn(r.Dy()/2) + r.Dy()/2

	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x++ {
			xx := x + xo
			yy := y + yo
			i := ((xx * xx) + (yy * yy)) / factor
			screen.Set(x, y, colors[i%ncolors])
		}
	}
}

func main() {
	xscr.Init(moire, 1e9)
	xscr.Run()
}
