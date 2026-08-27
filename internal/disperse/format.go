package disperse

import (
	"fmt"
	"strings"
)

func formatValue(v float64) string {
	s := fmt.Sprintf("%.*g", 6, v)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	u := uint64(v)
	if neg {
		u = uint64(-v)
	}
	var buf [24]byte
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
