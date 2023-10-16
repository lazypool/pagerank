package utils

import "sort"

// sort the map by its value then return a slice of its key
// from small to large when 'asc' is 'true', contrary when 'false'
func SortMapByValue(m *map[int]float64, asc bool) []int {
	var res []int
	for k := range *m {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool {
		return (*m)[res[i]] < (*m)[res[j]] == asc
	})
	return res
}
