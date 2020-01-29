package main

import (
	"fmt"
	"image/color"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
)

type Size struct {
	Width  int
	Height int
}

func download(size Size, color color.RGBA) {
	var fileName = fmt.Sprintf("images/%04dx%04d", size.Width, size.Height)
	if _, err := os.Stat(fileName); err != nil && !os.IsNotExist(err) {
		panic(err)
	} else if err == nil {
		return
	}
	fmt.Println("downloading", fileName)
	resp, err := http.Get(fmt.Sprintf("https://place-hold.it/%dx%d/%02x%02x%02x.jpg?text=", size.Width, size.Height, color.R, color.G, color.B))
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()
	if _, err := io.Copy(file, resp.Body); err != nil {
		panic(err)
	}
}

func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
	}
}

func main() {
	if err := os.MkdirAll("images", os.ModePerm); err != nil {
		panic(err)
	}
	var sizes = []Size{
		{Width: 100, Height: 100},
		{Width: 1920, Height: 1080},
	}
	var wg sync.WaitGroup
	for _, s := range sizes {
		wg.Add(1)
		var size = s
		go func() {
			download(size, randomColor())
			wg.Done()
		}()
	}
	for i := 0; i < 100; i++ {
		var base = rand.Intn(1901) + 100
		var k = float64(rand.Intn(11)+5) / 10
		var size = Size{
			Width:  base,
			Height: int(float64(base) * k),
		}
		wg.Add(1)
		go func() {
			download(size, randomColor())
			wg.Done()
		}()
	}
	wg.Wait()
}
