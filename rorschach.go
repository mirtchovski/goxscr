/* xscreensaver, Copyright (c) 1992, 1996, 1998, 2001
 *  Jamie Zawinski <jwz@jwz.org>
 *
 * Permission to use, copy, modify, distribute, and sell this software and its
 * documentation for any purpose is hereby granted without fee, provided that
 * the above copyright notice appear in all copies and that both that
 * copyright notice and this permission notice appear in supporting
 * documentation.  No representations are made about the suitability of this
 * software for any purpose.  It is provided "as is" without express or 
 * implied warranty.
 *
 * 19971004: Johannes Keukelaar <johannes@nada.kth.se>: Use helix screen
 *           eraser.
 */

/* 
 * ported to Plan 9 by andrey@lanl.gov, 08/02
 */

// Rewritten in Go by mirtchovski@gmail.com, 10/10

package main

import (
	"image"
	"rand"
	"exp/draw"
	"time"
	"./xscr"
)

var iterations = 40000
var offset = 4

func hurm(screen draw.Image) {
	var points = [4]image.Point{}

	xsym := rand.Intn(2)
	ysym := rand.Intn(2)

	color := xscr.RandomCmap(1)[0]
	black := image.ColorImage{image.RGBAColor{0, 0, 0, 0xff}}

	draw.Draw(screen, screen.Bounds(), black, image.ZP)
	xscr.Flush()

	xlim := screen.Bounds().Dx()
	ylim := screen.Bounds().Dy()

	x := xlim / 2
	y := ylim / 2

	for i := 0; i < iterations; i++ {
		j := 0

		x += rand.Intn(1+(offset<<1)) - offset
		y += rand.Intn(1+(offset<<1)) - offset

		points[j].X = x
		points[j].Y = y

		j++
		if xsym > 0 {
			points[j].X = xlim - x
			points[j].Y = y
			j++
		}
		if ysym > 0 {
			points[j].X = x
			points[j].Y = ylim - y
			j++
		}
		if xsym > 0 && ysym > 0 {
			points[j].X = xlim - x
			points[j].Y = ylim - y
			j++
		}
		for i2 := 0; i2 < j; i2++ {
			screen.Set(points[i2].X, points[i2].Y, color)
		}
		xscr.Flush()
		time.Sleep(10)
	}
}

func main() {

	xscr.Init(hurm, 1e9)
	xscr.Run()
}
