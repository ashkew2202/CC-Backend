package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func main() {

	chartBoolptr := flag.Bool("chart", false, "a bool")
	fileNameptr := flag.String("chartFileName", "default.html", "The default name of the file where the charts are plotted")
	flag.Parse()
	// endpointsBoolptr := flag.Bool("endpoints", false, "analyze endpoints")
	// perfBoolptr := flag.Bool("performance", false, "collect performance metrics")
	// appInsightsBoolptr := flag.Bool("appSpecificInsights", false, "gather app-specific insights")
	// uidAnalBoolptr := flag.Bool("uniqueIDanalysis", false, "perform unique ID analysis")
	f, err := os.OpenFile("timetable.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var logs [][]string
	for scanner.Scan() {
		line := scanner.Text()
		list := strings.Fields(line)
		logs = append(logs, list)
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		trafficAnal(logs)
	}()
	go func() {
		defer wg.Done()
		appSpecificInsights(logs)
	}()
	go func() {
		defer wg.Done()
		idAnal(logs)
	}()
	wg.Wait()
	trafficAnalChan := <-logParser(logs)
	trafficAnalResult := trafficAnalChan
	if *chartBoolptr {
		plotEverything(trafficAnalResult, *fileNameptr)
	}
}

func UnitConv(Val string) float64 {
	re := regexp.MustCompile(`([\d.]+)\s*([a-zA-Zµ]+)`)
	var resTimeFloat float64
	numPart := ""
	unitPart := ""
	matches := re.FindStringSubmatch(Val)
	if len(matches) == 3 {
		numPart = matches[1]
		unitPart = matches[2]
	}
	if numPart != "" && unitPart != "" {
		unit := strings.ToLower(unitPart)
		val, err := strconv.ParseFloat(numPart, 64)
		if err == nil {
			// The target is to convert all the response times to ms only
			switch unit {
			case "ms":
				resTimeFloat = val
			case "s":
				resTimeFloat = val * 1000
			case "us", "µs":
				resTimeFloat = val / 1000
			default:
				resTimeFloat = val
			}
		}
	}
	return resTimeFloat
}
