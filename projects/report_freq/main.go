package main

import (
	"machine"
	"math"
	"time"
)

var threshold uint16 = 4736

func errorOut() {
	for {
		machine.LED.Set(!machine.LED.Get())
		println("Error")
		time.Sleep(500 * time.Millisecond)
	}
}

func floor(x float64) float64 {
	if x == 0 || math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}
	if x < 0 {
		d, fract := math.Modf(-x)
		if fract != 0.0 {
			d = d + 1
		}
		return -d
	}
	d, _ := math.Modf(x)
	return d
}

func main() {
	machine.InitADC()
	machine.InitSerial()
	pin := machine.ADC{Pin: machine.ADC0}
	pin.Configure(machine.ADCConfig{})
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	println("Works!")

	samplingInterval := 1 * time.Millisecond
	var previousValue uint16 = 0
	lastCrossing := time.Now()

	var frequency float64 = 0.0

	var i = 0
	for {
		currentValue := pin.Get()

		// Detect a rising edge (simple approach: check if the signal crosses a threshold)
		if currentValue > threshold && previousValue <= threshold {
			// Calculate time difference between crossings
			timeNow := time.Now()
			period := timeNow.Sub(lastCrossing)
			lastCrossing = timeNow

			// Calculate frequency
			frequency = float64(time.Second) / float64(period)
		}

		if i%1000 == 0 {
			println("Frequency:", uint(floor(frequency)), "Hz")
		}

		machine.LED.Set(currentValue > threshold)

		previousValue = currentValue
		time.Sleep(samplingInterval)
		i++
	}
}
