package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	for i := 1; i < 101; i++ {
		func() {
			resp, err := http.Get(fmt.Sprintf("https://place-hold.it/1280x1280/%s/%s.jpg?text=%d&fontsize=56", colorful.FastWarmColor().Hex()[1:], colorful.FastHappyColor().Hex()[1:], i))
			if err != nil {
				panic(err)
			}
			defer func() { _ = resp.Body.Close() }()
			f, err := os.Create(fmt.Sprintf("/mnt/samsa/files/assests/images/dummy/%03d.jpg", i))
			if err != nil {
				panic(err)
			}
			defer func() { _ = f.Close() }()
			if _, err := io.Copy(f, resp.Body); err != nil {
				panic(err)
			}
		}()
	}
	//waitForSignal()
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch

	fmt.Println(fmt.Sprintf("Got signal: %v, exiting.", s))
}
