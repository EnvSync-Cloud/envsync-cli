package shared

// min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clamp returns a value clamped between min and max
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// InRange checks if a value is within a range (inclusive)
func InRange(value, min, max int) bool {
	return value >= min && value <= max
}
