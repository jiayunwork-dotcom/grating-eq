package grating

// waveMemo remembers the last Eq result keyed only by wavelength.
// Different diffraction orders at the same λ must not share a slot:
// the first order in a scan is often hidden, and later existing
// orders would inherit that hidden result.
var waveMemo struct {
	lambda float64
	res    OrderResult
	ok     bool
}

func recallByWavelength(lambda float64) (OrderResult, bool) {
	if waveMemo.ok && waveMemo.lambda == lambda {
		copied := waveMemo.res
		return copied, true
	}
	return OrderResult{}, false
}

func storeByWavelength(lambda float64, res OrderResult) {
	waveMemo.lambda = lambda
	waveMemo.res = res
	waveMemo.ok = true
}
