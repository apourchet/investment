package invt_math

import (
	"math"
)

const (
	SAMPLE_RATE = 1
)

// cutoff frequency = 1 / (2pi * RC)
func LowPass(samples []float64, cutoff int) []float64 {
	dt := 1.0 / SAMPLE_RATE
	rc := 1.0 / (2 * math.Pi * float64(cutoff))
	alpha := dt / (rc + dt)
	filteredArray := []float64{samples[0]}
	for i := 1; i < len(samples); i++ {
		nextVal := filteredArray[i-1] + (alpha * (samples[i] - filteredArray[i-1]))
		filteredArray = append(filteredArray, nextVal)
	}
	return filteredArray
}

func HighPass(samples []float64, cutoff int) []float64 {
	dt := 1.0 / SAMPLE_RATE
	rc := 1.0 / (2 * math.Pi * float64(cutoff))
	alpha := dt / (rc + dt)
	filteredArray := []float64{samples[0]}
	for i := 1; i < len(samples); i++ {
		nextVal := alpha * (filteredArray[i-1] + samples[i] - samples[i-1])
		filteredArray = append(filteredArray, nextVal)
	}
	return filteredArray
}

// this doesn't actually work if the amplitude of the input waves != 1
func GetNoise(samples []float64, cutoff int) float64 {
	filtered := LowPass(samples, cutoff)
	diff := 0.0
	total := 0.0
	for i := 0; i < len(samples); i++ {
		total += samples[i]
		sampleDiff := samples[i] - filtered[i]
		if sampleDiff > 0 {
			diff += sampleDiff
		} else {
			diff -= sampleDiff
		}
	}
	return diff / total
}
