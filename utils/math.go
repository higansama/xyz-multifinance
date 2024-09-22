package utils

import "math"

func RoundToNearest500k(amount float64) float64 {
	return math.Ceil(amount/500000) * 500000
}
