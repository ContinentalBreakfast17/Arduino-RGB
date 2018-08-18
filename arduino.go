package main

import (
	"os"

	"github.com/ContinentalBreakfast17/seriard"
)

var RGB_PINS [3]int

func initArduino() *seriard.Arduino {
	arduino, err := seriard.Connect(seriard.MODEL_UNO, os.Getenv("RGB_PORT"), seriard.BAUD_9600)
	errorHandler("Failed to connect to Arduino", err, false)

	RGB_PINS = [...]int{9, 10, 11}

	for _, pin := range RGB_PINS {
		_, err = arduino.SetPinMode(pin, seriard.MODE_OUTPUT)
		errorHandler("Failed to set pin mode", err, false)
	}

	return arduino
}

func writeColor(arduino *seriard.Arduino, channel, val int) {
	_, err := arduino.AnalogWrite(RGB_PINS[channel], uint8(val))
	errorHandler("Failed to write to Arduino", err, true)
}

func writeColors(arduino *seriard.Arduino, colors []int) {
	for i, _ := range RGB_PINS {
		_, err := arduino.AnalogWrite(RGB_PINS[i], uint8(colors[i]))
		errorHandler("Failed to write to Arduino", err, true)
	}
}