package main

import "github.com/go-echarts/go-echarts/v2/opts"

func generatePieItems(fields []string, counts []int) []opts.PieData {
	items := make([]opts.PieData, 0, len(fields))
	for i, f := range fields {
		var v interface{} = 0
		if i < len(counts) {
			v = counts[i]
		}
		items = append(items, opts.PieData{Name: f, Value: v})
	}
	return items
}

func generateBarItemsInts(values []int) []opts.BarData {
	items := make([]opts.BarData, 0, len(values))
	for _, v := range values {
		items = append(items, opts.BarData{Value: v})
	}
	return items
}

func generateBarItemsFloats(values []float64) []opts.BarData {
	items := make([]opts.BarData, 0, len(values))
	for _, v := range values {
		items = append(items, opts.BarData{Value: v})
	}
	return items
}
