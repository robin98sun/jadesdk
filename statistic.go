package jadesdk

import (
	"math"
	"sync"
	"time"
)

type UnitStat struct {
	Count     int64 `json:"count,omitempty"`
	Mean      int64 `json:"mean,omitempty"`
	Variance  int64 `json:"variance,omitempty"`
	Deviation int64 `json:"deviation,omitempty"`
	mutex     *sync.Mutex
}

func newUnitStat() *UnitStat {
	return &UnitStat{
		mutex: &sync.Mutex{},
	}
}

func (s *UnitStat) AddMilliseconds(durationMilliseconds int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var Xn float64
	Xn = float64(durationMilliseconds)

	var N float64
	N = float64(s.Count + 1)

	var meanPre float64
	meanPre = float64(s.Mean)

	var varancePre float64
	varancePre = float64(s.Variance)

	var meanN float64
	meanN = (meanPre*(N-1) + Xn) / N

	var varanceN float64
	varanceN = (N - 1) / N * varancePre
	varanceN += (N - 1) / N * (meanN - meanPre) * (meanN - meanPre)
	varanceN += 1 / N * (Xn - meanN) * (Xn - meanN)

	s.Count = int64(N)
	s.Mean = int64(meanN)
	s.Variance = int64(varanceN)
	s.Deviation = int64(math.Sqrt(varanceN))
}

func (s *UnitStat) AddDuration(duration time.Duration) {
	s.AddMilliseconds(int64(duration / time.Millisecond))
}

type Stat struct {
	Decoding       *UnitStat `json:"decoding,omitempty"`
	Task           *UnitStat `json:"task,omitempty"`
	Forwarding     *UnitStat `json:"forwarding,omitempty"`
	ReportToMaster *UnitStat `json:"reportToMaster,omitempty"`
	OnFly          *UnitStat `json:"onFly,omitempty"`
	Total          *UnitStat `json:"total,omitempty"`
}

func NewStat() *Stat {
	return &Stat{
		Decoding:       newUnitStat(),
		Task:           newUnitStat(),
		Forwarding:     newUnitStat(),
		ReportToMaster: newUnitStat(),
		OnFly:          newUnitStat(),
		Total:          newUnitStat(),
	}
}

type StatItem struct {
	Decoding   time.Duration `json:"decoding,omitempty"`
	Task       time.Duration `json:"task,omitempty"`
	Forwarding time.Duration `json:"forwarding,omitempty"`
}

func (s *Stat) ApplyItem(item *StatItem) {
	if item == nil {
		return
	}
	s.Decoding.AddDuration(item.Decoding)
	s.Task.AddDuration(item.Task)
	s.Forwarding.AddDuration(item.Forwarding)
}
