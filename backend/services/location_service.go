package services

import "math"

// IsWithinRadius checks if a given point (lat1, long1) is within a specified radius of another point (lat2, long2).
func IsWithinRadius(lat1, long1, lat2, long2, radius float64) bool {
	const earthRadius = 6371 // Radius of the Earth in kilometers

	// Convert latitude and longitude from degrees to radians
	lat1Rad := degreesToRadians(lat1)
	long1Rad := degreesToRadians(long1)
	lat2Rad := degreesToRadians(lat2)
	long2Rad := degreesToRadians(long2)

	// Haversine formula
	diffLong := long2Rad - long1Rad
	diffLat := lat2Rad - lat1Rad

	a := math.Sin(diffLat/2)*math.Sin(diffLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(diffLong/2)*math.Sin(diffLong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c

	return distance <= radius
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
