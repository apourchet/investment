package invt_math

import (
	"math"
	"testing"
)

// assume 10k samples per one time unit
func sineWave(hz, samples int, amp, phase float64) []float64 {
	wave := []float64{}
	frequencyRadiansPerSample := float64(hz) * 2.0 * math.Pi / 10000.0
	for i := 0; i < samples; i++ {
		phase += frequencyRadiansPerSample
		sampleValue := amp * math.Sin(phase)
		wave = append(wave, sampleValue)
	}
	return wave
}

// TODO actually test
func TestLowPass(t *testing.T) {
	combined := []float64{}
	lowFreqWave := sineWave(100, 100, 1, 0)
	highFreqWave := sineWave(10000, 100, 1, 0)
	for i := 0; i < len(lowFreqWave); i++ {
		combined = append(combined, lowFreqWave[i]+highFreqWave[i])
	}
	//lp := LowPass(combined, 1000)
	//for i := 0; i < len(lowFreqWave); i++ {
	//	fmt.Printf("%f,", lp[i])
	//}
	//for i := 0; i < len(lowFreqWave); i++ {
	//	fmt.Printf("[%f, %d]\n", lowFreqWave[i], i)
	//}
	//fmt.Println()
	//for i := 0; i < len(lowFreqWave); i++ {
	//	fmt.Printf("[%f, %d]\n", lp[i], i)
	//}
}
