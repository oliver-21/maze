package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/teacat/noire"
)

func pythagDist(a, b coor[float64]) float64 {
	dx := (a.x - b.x)
	dy := (a.y - b.y)
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func darkenOutside(original *ebiten.Image, pos coor[float64], max coor[int]) {
	bounds := original.Bounds()
	img := image.NewRGBA(bounds)

	var wg sync.WaitGroup
	fmt.Println(bounds.Max)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; x < bounds.Max.Y; y++ {
			col := original.At(x, y)
			wg.Add(1)
			func(x, y int) {
				defer wg.Done()
				if pythagDist(pos, coor[float64]{float64(x), float64(y)}) > 20 {
					r, g, b, a := col.RGBA()
					col = darken(color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
				}
				img.Set(x, y, col)
			}(x, y)
		}
	}
	wg.Wait()
	original.DrawImage(ebiten.NewImageFromImage(img), &ebiten.DrawImageOptions{})
}

func darken(c color.RGBA) color.RGBA {
	noi := noire.NewRGB(float64(c.R), float64(c.G), float64(c.A))
	noi.Darken(0.15)

	r, g, b := noi.RGB()
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
