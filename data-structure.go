package main

type TrafficAnal struct {
	totalRequests  int
	endpointCounts map[string]int
	statusCodes    map[string]int
	resTime map[string][2]float64
	stratCounts map[string]int
	timetables []int
	idCountMap map[string]int
	unique_ids []string
}

func NewTrafficAnal() TrafficAnal {
	return TrafficAnal{
		totalRequests:  0,
	    endpointCounts : make(map[string]int),
		statusCodes : make(map[string]int),
		resTime : make(map[string][2]float64),
		stratCounts : make(map[string]int),
		timetables : []int{},
		idCountMap : make(map[string]int),
		unique_ids : []string{},
	}
}
