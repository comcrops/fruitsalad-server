package model

import (
	"math/rand"
)

type RgbValue struct {
	Red uint8
	Green uint8
	Blue uint8
}

func GetRandomRgbValue() RgbValue {
	red := uint8(rand.Int())
	green := uint8(rand.Int())
	blue := uint8(rand.Int())

	return RgbValue{
		red,
		green,
		blue,
	}
}

