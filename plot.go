package main

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)


// This is very AI because of it being very repetitive

func plotEverything(trafficAnal TrafficAnal, filename string) error {
	var endpoints []string
	var counts []int
	// Bar graph for endpoint popularity
	for endpoint, count := range trafficAnal.endpointCounts {
		endpoints = append(endpoints, endpoint)
		counts = append(counts, count)
	}

	// Bar graph for HTTP status code counts
	var statusCodes []string
	var statusCounts []int
	for code, count := range trafficAnal.statusCodes {
		statusCodes = append(statusCodes, code)
		statusCounts = append(statusCounts, count)
	}

	// Bar graph for max and avg response times per endpoint
	var respEndpoints []string
	var maxTimes []float64
	var avgTimes []float64
	for endpoint, stats := range trafficAnal.resTime {
		respEndpoints = append(respEndpoints, endpoint)
		maxTimes = append(maxTimes, stats[1])
		avgTimes = append(avgTimes, stats[0])
	}

	// Pie chart for unique IDs divided into different fields
	var pieBatch []string
	var pieBatchSpecificCounts []int
	for field, count := range trafficAnal.idCountMap {
		pieBatch = append(pieBatch, field)
		pieBatchSpecificCounts = append(pieBatchSpecificCounts, count)
	}

	// Plot bar chart for endpoint popularity
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Endpoint Popularity"}))
	bar.SetXAxis(endpoints).
		AddSeries("Hits", generateBarItemsInts(counts))

	// Plot donut chart for unique ID fields
	donut := charts.NewPie()
	donut.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Unique ID Fields"}))
	donut.AddSeries("Fields", generatePieItems(pieBatch, pieBatchSpecificCounts)).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Formatter: "{b}: {c}"}))

	// Plot bar chart for HTTP status codes
	statusBar := charts.NewBar()
	statusBar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "HTTP Status Codes"}))
	statusBar.SetXAxis(statusCodes).
		AddSeries("Counts", generateBarItemsInts(statusCounts))

	// Plot bar chart for max and avg response times per endpoint
	respBar := charts.NewBar()
	respBar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Response Times per Endpoint"}))
	respBar.SetXAxis(respEndpoints).
		AddSeries("Max Response Time", generateBarItemsFloats(maxTimes)).
		AddSeries("Avg Response Time", generateBarItemsFloats(avgTimes))

	// Render all charts in a single page
	page := components.NewPage()
	page.AddCharts(
		bar,
		statusBar,
		respBar,
		donut,
	)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := page.Render(f); err != nil {
		return err
	}
	return nil
}


