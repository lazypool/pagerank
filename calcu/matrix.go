package calcu

import (
	"encoding/json"
	"os"
)

// the parse COO matrix
type Matrix struct {
	Row []int
	Col []int
	Val []float64
}

// append an element to the matrix
func (m *Matrix) Set(row int, col int, val float64) {
	m.Row = append(m.Row, row)
	m.Col = append(m.Col, col)
	m.Val = append(m.Val, val)
}

// save the matrix to a json file
func (m *Matrix) Dump(path string) {
	data, _ := json.Marshal(m)
	os.WriteFile(path, data, 0644)
}

// read the matrix from a json file
func (m *Matrix) Load(path string) {
	data, _ := os.ReadFile(path)
	json.Unmarshal(data, m)
}

// multiply a matrix and a vector, transpose matrix if 'trans' is 'true'
func (m *Matrix) Mult(vec *map[int]float64, trans bool) map[int]float64 {
	dst := make(map[int]float64)
	if !trans {
		for idx, r := range m.Row {
			dst[r] += m.Val[idx] * (*vec)[m.Col[idx]]
		}
	} else {
		for idx, c := range m.Col {
			dst[c] += m.Val[idx] * (*vec)[m.Row[idx]]
		}
	}
	return dst
}

// do a normalize to the matrix, the 'byRow' indicated how to normalize it
func (m *Matrix) Normal(byRow bool) {
	sums := make(map[int]float64)
	if byRow {
		for idx, r := range m.Row {
			sums[r] += m.Val[idx]
		}
		for idx, r := range m.Row {
			m.Val[idx] /= sums[r]
		}
	} else {
		for idx, c := range m.Col {
			sums[c] += m.Val[idx]
		}
		for idx, c := range m.Col {
			m.Val[idx] /= sums[c]
		}
	}
}
