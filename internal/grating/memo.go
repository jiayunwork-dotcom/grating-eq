package grating

// waveMemo remembers the last Eq result keyed by wavelength, signed
// order, and the rest of the grating geometry. Different orders at the
// same λ must not share a slot: the first order in a scan is often
// hidden, and later existing orders would inherit that hidden result.
type waveKey struct {
	lambda float64
	m      int
	d      float64
	si     float64
}

var waveMemo = map[waveKey]OrderResult{}

func waveMemoKey(g Grating, m int) waveKey {
	return waveKey{
		lambda: g.WavelengthNm,
		m:      m,
		d:      g.GrooveSpacingNm,
		si:     g.IncidentSine(),
	}
}

func recallByWavelength(g Grating, m int) (OrderResult, bool) {
	v, ok := waveMemo[waveMemoKey(g, m)]
	return v, ok
}

func storeByWavelength(g Grating, m int, res OrderResult) {
	waveMemo[waveMemoKey(g, m)] = res
}
