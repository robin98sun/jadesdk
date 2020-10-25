package jadesdk

import (
	"math"
	"sync"
	"time"
)

type UnitStat struct {
	Count       int64 `json:"count,omitempty"`
	Mean        int64 `json:"mean,omitempty"`
	MeanN       int64 `json:"meanN,omitempty"`
	VarianceN   int64 `json:"varianceN,omitempty"`
	DeviationN  int64 `json:"deviationN,omitempty"`
	latestItems []int64
	mutex       *sync.Mutex
	N           int64 `json:"N,omitempty"`
}

func newUnitStat() *UnitStat {
	return &UnitStat{
		mutex:       &sync.Mutex{},
		N:           10,
		latestItems: []int64{},
	}
}

func (s *UnitStat) AddMicroseconds(durationMicroseconds int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.latestItems = append(s.latestItems, durationMicroseconds)
	s.Mean = (s.Mean*s.Count + durationMicroseconds) / (s.Count + 1)
	s.Count++

	n := int64(len(s.latestItems))
	if n > s.N {
		s.latestItems = s.latestItems[n-s.N:]
	}
	var sum int64 = 0
	for _, d := range s.latestItems {
		sum += d
	}
	s.MeanN = sum / n

	s.VarianceN = 0
	for _, d := range s.latestItems {
		s.VarianceN += (d - s.MeanN) * (d - s.MeanN)
	}
	s.VarianceN /= n
	s.DeviationN = int64(math.Sqrt(float64(s.VarianceN)))
}

func (s *UnitStat) AddDuration(duration time.Duration) {
	s.AddMicroseconds(int64(duration / time.Microsecond))
}

type Stat struct {
	Decoding           *UnitStat
	Task               *UnitStat
	Forwarding         *UnitStat
	ReportToMaster     *UnitStat
	ReportToAggregater *UnitStat
}

func newStat() *Stat {
	return &Stat{
		Decoding:           newUnitStat(),
		Task:               newUnitStat(),
		Forwarding:         newUnitStat(),
		ReportToAggregater: newUnitStat(),
		ReportToMaster:     newUnitStat(),
	}
}
