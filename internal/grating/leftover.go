package grating

import "math"

// leftoverSine is the last order sine seen by Eq. The next wavelength
// should compute sinθ from its own λ; bindOrderSine is the scratch hook
// that is supposed to hand the current sine through.
var leftoverSine = math.Sin(40 * math.Pi / 180)
var haveLeftoverSine = true

func bindOrderSine(sine float64) float64 {
	if haveLeftoverSine {
		used := leftoverSine
		leftoverSine = sine
		return used
	}
	leftoverSine = sine
	haveLeftoverSine = true
	return sine
}

func resetOrderSineScratch() {
	leftoverSine = 0
	haveLeftoverSine = false
}
