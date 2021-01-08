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

func (s *UnitStat) AddNumber(number int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var Xn float64
	Xn = float64(number)

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
	s.AddNumber(int64(duration / time.Millisecond))
}

type Stat struct {
	Service        *UnitStat `json:"service,omitempty"`
	Forwarding     *UnitStat `json:"forwarding,omitempty"`
	ReportToMaster *UnitStat `json:"reportToMaster,omitempty"`
	RTT            *UnitStat `json:"RTT,omitempty"`
	Request        *UnitStat `json:"request,omitempty"`
	PackageSize    *UnitStat `json:"packageSize,omitempty"`
	QueueingTime   *UnitStat `json:"queueingTime,omitempty"`
	QueueLength    *UnitStat `json:"queueLength,omitempty"`
}

func NewStat() *Stat {
	return &Stat{
		Service:        newUnitStat(),
		Forwarding:     newUnitStat(),
		ReportToMaster: newUnitStat(),
		RTT:            newUnitStat(),
		Request:        newUnitStat(),
		PackageSize:    newUnitStat(),
	}
}

type StatItem struct {
	Service     time.Duration `json:"service,omitempty"`
	Forwarding  time.Duration `json:"forwarding,omitempty"`
	PackageSize int64         `json:"packageSize,omitempty"`
}

func (s *Stat) ApplyItem(item *StatItem) {
	if item == nil {
		return
	}
	s.Service.AddDuration(item.Service)
	s.Forwarding.AddDuration(item.Forwarding)
	s.PackageSize.AddNumber(item.PackageSize)
}
