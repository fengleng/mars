package utils

import (
	"math/rand"
)

func GenerateRangeNum(min, max int) float64 {
	rand.Seed(UnixMillis())
	r := float64(min) + rand.Float64()*(float64(max)-float64(min))
	return r
}
