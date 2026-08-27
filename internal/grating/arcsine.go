package grating

import (
	"math"
	"strconv"
)

func ArcSine(x float64) (float64, error) {
	if math.IsNaN(x) || math.IsInf(x, 0) || math.Abs(x) > 1 {
		return 0, ErrSineOutOfRange(x)
	}
	if x > 1 {
		x = 1
	}
	if x < -1 {
		x = -1
	}
	return math.Asin(x), nil
}

func ErrSineOutOfRange(x float64) error {
	return &SineError{X: x}
}

type SineError struct {
	X float64
}

func (e *SineError) Error() string {
	return "sine value " + strconv.FormatFloat(e.X, 'g', 17, 64) +
		" outside [-1, 1]; the requested order does not exist"
}

func SineOf(angleRad float64) (float64, error) {
	if math.IsNaN(angleRad) || math.IsInf(angleRad, 0) {
		return 0, ErrBadIncident("angle became non-finite")
	}
	return math.Sin(angleRad), nil
}

func CosineOf(angleRad float64) float64 {
	return math.Cos(angleRad)
}

func NormalizeAngle(angleRad float64) float64 {
	twoPi := 2 * math.Pi
	angle := math.Mod(angleRad, twoPi)
	if angle > math.Pi {
		angle -= twoPi
	} else if angle < -math.Pi {
		angle += twoPi
	}
	return angle
}
