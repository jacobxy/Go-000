package main

import (
	"fmt"
	"math/rand"
	"sliding/window"
	"time"
)

func main() {
	sld := window.NewSliding(10, 19)
	rand.Seed(time.Now().Unix())
	cnt := 0
	for cnt < 3 {
		r := rand.Int63() % 200
		time.Sleep(time.Duration(r) * time.Millisecond)
		avg, ok := sld.Allow()
		fmt.Println(avg, ok)
		if !ok {
			time.Sleep(1 * time.Second)
			cnt++
		}
	}
}
