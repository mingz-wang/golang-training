// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

var palette = []color.Color{color.White, color.RGBA{0, 0xff, 0, 0xff},
	color.RGBA{0, 0, 0xff, 0xff}, color.RGBA{0xff, 0, 0, 0xff},
	color.RGBA{0x22, 0x22, 0x22, 0xff}, color.RGBA{0xff, 0x44, 0x15, 0xff},
	color.RGBA{0x44, 0x62, 0x12, 0xff}, color.RGBA{0xff, 0xff, 0x15, 0xff}}

const (
	whiteIndex = 0 // first color in palette
	greenIndex = 1 // next color in palette
	blueIndex  = 2
	redIndex   = 3
)

func main() {
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.

	// rand.Seed(time.Now().UTC().UnixNano())
	// rand.Seed(seed) go 1.20 已弃用
	// https://tip.golang.org/doc/go1.20
	outFilename := os.Args[1]
	f, err := os.Create(outFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	lissajous(f)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complex x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	freq := rand.Float64() * 3.0 // relative frequency of oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	colorIndex := 1
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(colorIndex))
		}
		colorIndex = (colorIndex % 7) + 1
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
