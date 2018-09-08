package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ContinentalBreakfast17/seriard"
)

var RGB_PINS [3]int

func initArduino() *seriard.Arduino {
	arduino := &seriard.Arduino{
		ModelName: seriard.MODEL_UNO,
		Baud: seriard.BAUD_115200,
		InitWait: 2*time.Second,
	}
	err := arduino.Connect(os.Getenv("RGB_PORT"))
	errorHandler("Failed to connect to Arduino", err, false)

	RGB_PINS = [...]int{9, 10, 11}

	for _, pin := range RGB_PINS {
		_, err = arduino.SetPinMode(pin, seriard.MODE_OUTPUT)
		errorHandler("Failed to set pin mode", err, false)
	}

	return arduino
}

func writeColor(arduino *seriard.Arduino, channel, val int) {
	err := arduino.AnalogWrite(RGB_PINS[channel], uint8(val))
	errorHandler("Failed to write to Arduino", err, true)
}

func writeColors(arduino *seriard.Arduino, colors []int) {
	for i, _ := range RGB_PINS {
		err := arduino.AnalogWrite(RGB_PINS[i], uint8(colors[i]))
		errorHandler("Failed to write to Arduino", err, true)
	}
}

func writeMode(arduino *seriard.Arduino, mode string) {
	err := arduino.CustomCommand("set_rgb_mode", mode)
	errorHandler("Failed to set rgb mode", err, true)
}

func writeSpeed(arduino *seriard.Arduino, speed int) {
	err := arduino.CustomCommand("set_speed", fmt.Sprintf("%d", speed))
	errorHandler("Failed to set rgb speed", err, true)
}