package main

import (
	"os"

	"github.com/ContinentalBreakfast17/seriard"
)

func initArduino() *seriard.Arduino {
	arduino, err := seriard.Connect(seriard.MODEL_UNO, os.Getenv("RGB_PORT"), seriard.BAUD_9600)
	if err != nil {
		panic(err)
	}
	return arduino
}