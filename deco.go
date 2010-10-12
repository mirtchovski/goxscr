/* xscreensaver, Copyright (c) 1997, 1998, 2002 Jamie Zawinski <jwz@jwz.org>
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
 * ported to Plan 9 by andrey@lanl.gov, 07/02
 */

// Rewritten in Go by mirtchovski@gmail.com 10/10
package main

import (
	"exp/draw"
	"rand"
	"image"

	"./xscr"
)

var ncolors int = 256
var maxDepth int = 10
var minHeight int = 0
var minWidth int = 0
var col int = 0

var colors []image.RGBAColor
var black = image.NewColorImage(image.RGBAColor{0, 0, 0, 0xff})

func deco1(screen draw.Image, x, y, w, h, depth int) {
	if rand.Intn(maxDepth+1) < depth || w <= minWidth || h <= minHeight {
		col = col + 1
		if col >= ncolors {
			col = 0
		}
		r := image.Rect(x, y, x+w, y+h)
		r = r.Add(screen.Bounds().Min)
		draw.Draw(screen, r, image.NewColorImage(colors[col]), image.ZP)
		draw.Border(screen, r, 1, black, image.ZP)
	} else {
		if rand.Intn(2) > 0 {
			deco1(screen, x, y, w/2, h, depth+1)
			deco1(screen, x+w/2, y, w/2, h, depth+1)
		} else {
			deco1(screen, x, y, w, h/2, depth+1)
			deco1(screen, x, y+h/2, w, h/2, depth+1)
		}
	}
}

func deco(screen draw.Image) {
	colors = xscr.RandomCmap(ncolors)

	draw.Draw(screen, screen.Bounds(), image.NewColorImage(colors[col]), image.ZP)
	deco1(screen, 0, 0, screen.Bounds().Dx(), screen.Bounds().Dy(), 0)
}

func main() {
	xscr.Init(deco, 10e9)
	xscr.Run()
}
