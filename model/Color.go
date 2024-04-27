package model

import (
	"crypto/rand"
	"encoding/binary"
)

type Color struct {
	Red uint8
	Green uint8
	Blue uint8
}

func newColor(red, green, blue uint8) Color {
	return Color{
		red,
		green,
		blue,
	}
}

func GetRandomRgbValue() Color {
	var red, green, blue uint8

	binary.Read(rand.Reader, binary.LittleEndian, &red)
	binary.Read(rand.Reader, binary.LittleEndian, &green)
	binary.Read(rand.Reader, binary.LittleEndian, &blue)

	return newColor(red, green, blue)
}

