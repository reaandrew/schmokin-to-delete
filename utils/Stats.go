package utils

func AverageFloat64(values []float64) (result float64) {
	for _, value := range values {
		result = result + value
	}
	result = result / float64(len(values))
	return
}

func Sum(values []int64) (result int64) {
	for _, value := range values {
		result += value
	}
	return
}

func Max(values []int64) (result int64) {
	for _, value := range values {
		if value > result {
			result = value
		}
	}
	return
}

func Min(values []int64) (result int64) {
	for _, value := range values {
		if result == 0 || value < result {
			result = value
		}
	}
	return
}
