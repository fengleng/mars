package utils

import "sort"

func SortUint64(list []uint64) {
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
}
