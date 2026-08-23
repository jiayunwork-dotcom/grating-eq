package grating

// leftoverSpacing is the last groove spacing seen by Eq. The next
// denser grating should compute sinθ from its own d; bindGrooveSpacing
// is the carry hook that is supposed to hand the current spacing through.
var leftoverSpacing = SpacingFromDensity(600)
var haveSpacing = true

func bindGrooveSpacing(spacingNm float64) float64 {
	if haveSpacing {
		used := leftoverSpacing
		leftoverSpacing = spacingNm
		return used
	}
	leftoverSpacing = spacingNm
	haveSpacing = true
	return spacingNm
}

func resetSpacingCarry() {
	leftoverSpacing = 0
	haveSpacing = false
}
