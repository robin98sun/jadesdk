package jadesdk

import (
	"math"
	"sync"
	"time"
)

type UnitStat struct {
	Count     int64 `json:"count,omitempty"`
	Mean      float64 `json:"mean,omitempty"`
	Variance  float64 `json:"variance,omitempty"`
	Deviation float64 `json:"deviation,omitempty"`
	Sum       int64 `json:"sum,omitempty"`
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
	if number > math.MaxInt64 - s.Sum {
		s.Count = 0
		s.Mean = float64(0)
		s.Variance = float64(0)
		s.Deviation = float64(0)
		s.Sum = 0
	}

	s.Sum += number

	var Xn float64
	Xn = float64(number)

	var N float64
	N = float64(s.Count + 1)

	var meanPre float64
	meanPre = s.Mean

	var varancePre float64
	varancePre = s.Variance

	var meanN float64
	meanN = (meanPre*(N-1) + Xn) / N

	var varanceN float64
	varanceN = (N - 1) / N * varancePre
	varanceN += (N - 1) / N * (meanN - meanPre) * (meanN - meanPre)
	varanceN += 1 / N * (Xn - meanN) * (Xn - meanN)

	s.Count = int64(N)
	s.Mean = meanN
	s.Variance = varanceN
	s.Deviation = math.Sqrt(varanceN)
}

func (s *UnitStat) AddDuration(duration time.Duration) {
	s.AddNumber(int64(duration))
}

type Stat struct {
	PreService 	   *UnitStat `json:"preService,omitempty"`
	Service        *UnitStat `json:"service,omitempty"`			// for aggregator, it is the time from enqueued to all subtasks finished
	ExecutionTime  *UnitStat `json:"executionTime,omitempty"`	// for aggregator, it is the cumulated time of processing each subtask result
	PostService    *UnitStat `json:"postService,omitempty"`
	Forwarding     *UnitStat `json:"forwarding,omitempty"`
	ReportToMaster *UnitStat `json:"reportToMaster,omitempty"`
	Communication  *UnitStat `json:"communication,omitempty"`
	Request        *UnitStat `json:"request,omitempty"`
	PackageSize    *UnitStat `json:"packageSize,omitempty"`
	QueueingTime   *UnitStat `json:"queueingTime,omitempty"`
	QueueLength    *UnitStat `json:"queueLength,omitempty"`
}

func NewStat() *Stat {
	return &Stat{
		PreService:     newUnitStat(),
		Service:        newUnitStat(),
		ExecutionTime:  newUnitStat(),
		PostService:    newUnitStat(),
		Forwarding:     newUnitStat(),
		ReportToMaster: newUnitStat(),
		Communication:  newUnitStat(),
		Request:        newUnitStat(),
		PackageSize:    newUnitStat(),
		QueueLength:    newUnitStat(),
		QueueingTime:   newUnitStat(),
	}
}

type StatItem struct {
	PreService  time.Duration `json:"preService,omitempty"`
	Service     time.Duration `json:"service,omitempty"` // for aggregator, it is the time from been enqueued to all subtasks finished
	Execution 	time.Duration `json:"execution,omitempty"` // for aggregator, it is the cumulated time of processing each subtask result
	PostService time.Duration `json:"postService,omitempty"`
	Forwarding  time.Duration `json:"forwarding,omitempty"`
	PackageSize int64         `json:"packageSize,omitempty"`
	RetryCountOfArrivalComm int64 `json:"retryCountOfArrivalComm,omitempty"`
}

func (s *Stat) ApplyItem(item *StatItem) {
	if item == nil {
		return
	}
	s.PreService.AddDuration(item.PreService)
	s.Service.AddDuration(item.Service)
	s.ExecutionTime.AddDuration(item.Execution)
	s.PostService.AddDuration(item.PostService)
	s.Forwarding.AddDuration(item.Forwarding)
	s.PackageSize.AddNumber(item.PackageSize)
}
