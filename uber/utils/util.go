package utils

import "math"

func CalculateDistance(fromLat, fromLong, toLat, toLong float64) float64 {
	
	return math.Sqrt(math.Pow(toLat-fromLat, 2) + math.Pow(toLong-fromLong, 2))
}
