package main

import (
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
)

func main() {
	// Initialize the periph library
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Replace this with the correct pin name for your Rock 3A board
	pinName := "GPIO12" // Example pin name

	// Get the GPIO pin
	pin, err := host.Pins(pinName)
	if err != nil {
		log.Fatal(err)
	}

	// Set the pin as an output
	if err := pin.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	// Toggle the pin
	for {
		pin.Out(gpio.High)
		time.Sleep(time.Second)
		pin.Out(gpio.Low)
		time.Sleep(time.Second)
	}
}
