package utils

import "math/rand"

func GetRandomColor() string {
	colors := []string{
		"grayModern",
		"error",
		"warning",
		"success",
		"grayWarm",
		"moss",
		"blueLight",
		"indigo",
		"violet",
		"pink",
	}
	return colors[rand.Intn(len(colors))]
}
