// Server1 is a minimal "echo" server.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var (
	mu    sync.Mutex
	count int
)

func main() {
	http.HandleFunc("/", handle) // each request calls handler
	http.HandleFunc("/count", counter)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/lissajous", lissajous)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

// handle echoes the Path component of the request URL
func handle(w http.ResponseWriter, r *http.Request) {
	IncrCount()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func IncrCount() {
	mu.Lock()
	count++
	mu.Unlock()
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

// echo echoes the HTTP request.
func echo(w http.ResponseWriter, r *http.Request) {
	IncrCount()
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Head[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

// lissajous gen lissajous gif
func lissajous(w http.ResponseWriter, r *http.Request) {
	IncrCount()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	var cycle float64
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
		if k == "cycle" {
			temp, err := strconv.Atoi(v[0])
			if err != nil {
				return
			}
			cycle = float64(temp)
		}
	}
	lissajousGIF(w, cycle)
}

// Lissajous generates GIF animations of random Lissajous figures.
var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func lissajousGIF(out io.Writer, cycle float64) {
	const (
		cycles  = 5     // number of complex x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	if cycle == 0 {
		cycle = 5
	}
	freq := rand.Float64() * 3.0 // relative frequency of oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycle*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
