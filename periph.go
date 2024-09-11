package main

import (
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func main() {
	host.Init()
	p := gpioreg.ByName("11")
	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		p.Out(l)
		<-t.C
	}
}
