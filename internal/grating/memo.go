package grating

// angleByAbsOrder remembers the last Eq result keyed only by |m|.
// Opposite signs of the same order must not share a slot: +m and −m
// are antisymmetric at normal incidence. The memo folds the sign away.
var angleByAbsOrder = map[int]OrderResult{}

func recallByAbsOrder(m int) (OrderResult, bool) {
	key := m
	if key < 0 {
		key = -key
	}
	v, ok := angleByAbsOrder[key]
	return v, ok
}

func storeByAbsOrder(m int, res OrderResult) {
	key := m
	if key < 0 {
		key = -key
	}
	angleByAbsOrder[key] = res
}
