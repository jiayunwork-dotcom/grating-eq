package report

import (
	"fmt"
	"strings"
)

type cell struct {
	text  string
	align string
}

type table struct {
	header []string
	align  []string
	rows   [][]cell
}

func newTable(header []string, align []string) *table {
	return &table{header: header, align: align}
}

func (t *table) addRow(values []string) {
	row := make([]cell, len(values))
	for i, v := range values {
		al := "left"
		if i < len(t.align) {
			al = t.align[i]
		}
		row[i] = cell{text: v, align: al}
	}
	t.rows = append(t.rows, row)
}

func (t *table) widths() []int {
	n := len(t.header)
	w := make([]int, n)
	for i, h := range t.header {
		w[i] = displayWidth(h)
	}
	for _, r := range t.rows {
		for i, c := range r {
			if i < n && displayWidth(c.text) > w[i] {
				w[i] = displayWidth(c.text)
			}
		}
	}
	return w
}

func (t *table) String() string {
	w := t.widths()
	var b strings.Builder
	writeRow := func(cells []string, align []string) {
		parts := make([]string, len(cells))
		for i, c := range cells {
			pad := w[i] - displayWidth(c)
			if i < len(align) && align[i] == "right" {
				parts[i] = strings.Repeat(" ", pad) + c
			} else {
				parts[i] = c + strings.Repeat(" ", pad)
			}
		}
		b.WriteString(strings.Join(parts, "  "))
		b.WriteString("\n")
	}
	writeRow(t.header, t.align)
	b.WriteString(strings.Repeat("-", totalWidth(w)) + "\n")
	for _, r := range t.rows {
		cells := make([]string, len(r))
		al := make([]string, len(r))
		for i, c := range r {
			cells[i] = c.text
			al[i] = c.align
		}
		writeRow(cells, al)
	}
	return b.String()
}

func totalWidth(w []int) int {
	total := 0
	for _, v := range w {
		total += v
	}
	if len(w) > 1 {
		total += 2 * (len(w) - 1)
	}
	return total
}

func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r > 0x2E7F {
			w += 2
		} else {
			w++
		}
	}
	return w
}

func pad(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
