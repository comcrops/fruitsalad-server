package model

import (
	"math/rand"
)

type Color struct {
	Red uint8
	Green uint8
	Blue uint8
}

func newColor() Color {
	return Color{
		0,
		0,
		0,
	}
}

func GetRandomRgbValue() Color {
	red := uint8(rand.Int())
	green := uint8(rand.Int())
	blue := uint8(rand.Int())

	return Color{
		red,
		green,
		blue,
	}
}

