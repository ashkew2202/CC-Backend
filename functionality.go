package main
import (
	"fmt"
	"strconv"
)
func logParser(logs [][]string) chan TrafficAnal {

	trafficAnalChan := make(chan TrafficAnal)
	go func() {
		traffic_anal := NewTrafficAnal()

		for res := 0; res < len(logs); res++ {
			if len(logs[res]) < 7 {
				if len(logs[res]) == 6 && logs[res][3] == "router:" {
					id := logs[res][5]
					idStr := id[1:len(id)-2]
					found := false
					for _, uid := range traffic_anal.unique_ids {
						if uid == idStr {
							found = true
							break
						}
					}
					if !found {
						traffic_anal.unique_ids = append(traffic_anal.unique_ids, idStr)
					}
				}
				continue
			}
			if logs[res][3] == "POST" || logs[res][3] == "GET" {
				traffic_anal.totalRequests += 1

				route := logs[res][4]
				traffic_anal.endpointCounts[route]++

				statusCode := logs[res][5]
				traffic_anal.statusCodes[statusCode]++

				resTimeVal := logs[res][6]
				resTimeFloat := UnitConv(resTimeVal)
				resVal := traffic_anal.resTime[route]
				resVal[0] += resTimeFloat
				if resTimeFloat > resVal[1] {
					resVal[1] = resTimeFloat
				}
				traffic_anal.resTime[route] = resVal
			}
			if len(logs[res]) > 5 {
				if logs[res][5] == "Iterative" {
					traffic_anal.stratCounts["Iterative Random Sampling"]++
				}

				if logs[res][5] == "Heuristic" {
					traffic_anal.stratCounts["Heuristic Backtracking Strategy"]++
				}
				if logs[res][4] == "Generation" && len(logs[res]) > 7 {
					val, err := strconv.Atoi(logs[res][7])
					if err == nil {
						traffic_anal.timetables = append(traffic_anal.timetables, val)
					}
				}
			}
			if len(logs[res]) == 6 && logs[res][3] == "router:" {
				id := logs[res][5]
				idStr := id[1:len(id)-2]
				found := false
				for _, uid := range traffic_anal.unique_ids {
					if uid == idStr {
						found = true
						break
					}
				}
				if !found {
					traffic_anal.unique_ids = append(traffic_anal.unique_ids, idStr)
				}
			}
		}
		for i := range traffic_anal.unique_ids {
			year := ""
			if len(traffic_anal.unique_ids[i]) >= 4 {
				year = traffic_anal.unique_ids[i][:4]
			}
			traffic_anal.idCountMap[year]++
		}
		trafficAnalChan <- traffic_anal
	}()
	return trafficAnalChan
}


func trafficAnal(logs [][]string)  {

	trafficAnalChan := logParser(logs)
	trafficAnal := <-trafficAnalChan
	totalRequests := trafficAnal.totalRequests
	endpointCounts := trafficAnal.endpointCounts
	statusCodes := trafficAnal.statusCodes
	resTime := trafficAnal.resTime
	var output string = fmt.Sprintf("Traffic & Usage Analysis\n----------------------\nTotal API Requests Logged: %d\n\nEndpoints Popularity:", totalRequests)
	for endpoint, count := range endpointCounts {
		output += fmt.Sprintf("\n\t- %s: %d (%f%%)", endpoint, count, 100*float64(count)/float64(totalRequests))
	}
	output += fmt.Sprintf("\n\nHTTP Status Codes:")
	for code, count := range statusCodes {
		output += fmt.Sprintf("\n\t- %s: %d times", code, count)
	}
	output += fmt.Sprintf("\n----------------------\nPerformance Metrics: \n----------------------")
	for endpoint, resTimes := range resTime {
		count := endpointCounts[endpoint]
		avg := 0.0
		if count > 0 {
			avg = float64(resTimes[0]) / float64(count)
		}
		output += fmt.Sprintf("\nEndpoint: %s\n\t - Average Response Time: %f ms\n\t - Max Response Time: %f ms", endpoint, avg, float64(resTimes[1]))
	}
	fmt.Println(output)
}

func appSpecificInsights(logs [][]string) {
	trafficAnalChan := logParser(logs)
	trafficAnal := <-trafficAnalChan
	stratCounts := trafficAnal.stratCounts
	timetables := trafficAnal.timetables
	total_timetables := 0
	for _, t := range timetables {
		total_timetables += t
	}
	var output string = "\n-------------------\nApp-Specific Insights\n-------------------"
	output += fmt.Sprintf("\nStrategy Usage:")
	for strat, count := range stratCounts {
		output += fmt.Sprintf("\n\t- %s: %d times", strat, count)
	}
	if len(timetables) > 0 {
		avg := float64(total_timetables) / float64(len(timetables))
		max := timetables[0]
		for _, t := range timetables {
			if t > max {
				max = t
			}
		}
		output += fmt.Sprintf("\n\nAverage Timetable Generation Time: %.2f ms", avg)
		output += fmt.Sprintf("\nMax Timetable Generation Time: %d ms", max)
	}
	fmt.Println(output)
}

func idAnal(logs [][]string) {
	trafficAnalChan := logParser(logs)
	trafficAnal := <-trafficAnalChan
	unique_ids := trafficAnal.unique_ids
	idCountMap := trafficAnal.idCountMap
	totalUniqueIDs := len(unique_ids)
	var output string = fmt.Sprintf("\n-------------------\nUnique ID Analysis\n-------------------\nTotal Unique IDs Found: %d", totalUniqueIDs)
	for year, count := range idCountMap {
		output += fmt.Sprintf("\n\tBatch of %s: %d unique IDs", year, count)
	}
	fmt.Println(output)
}