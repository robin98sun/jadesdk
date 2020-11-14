package jadesdk

import (
	"sync"
)

type AggregativeTaskCache struct {
	Cache map[string]*AggregativeTaskCacheItem // taskId: taskQueueItem
	mutex *sync.Mutex
}

func NewAggregativeTaskCache() *AggregativeTaskCache {
	inst := &AggregativeTaskCache{
		Cache: make(map[string]*AggregativeTaskCacheItem),
		mutex: &sync.Mutex{},
	}
	return inst
}

func (q *AggregativeTaskCache) Lock() {
	q.mutex.Lock()
}
func (q *AggregativeTaskCache) Unlock() {
	q.mutex.Unlock()
}

type AggregativeTaskCacheItem struct {
	Subtasks   map[string]*AggregativeTaskCacheSubtaskItem // subtaskId: subtaskItem
	Cumulation interface{}                                 // cumulative Result of so far responded subtasks
	ReportTo   []*Interface                                // upper layer aggregators
	SubtaskKey string                                      // the id of the aggregator subtask
}

type AggregativeTaskCacheSubtaskItem struct {
	SubtaskKey string
	Result     interface{}
}

func (q *AggregativeTaskCache) EnqueueAggregativeTask(msg *AggregatorEnqueuingMessage) {
	if q == nil || msg == nil || msg.TaskKey == "" || len(msg.Subtasks) == 0 {
		return
	}
	q.Lock()
	defer q.Unlock()
	if q.Cache == nil || len(q.Cache) == 0 {
		q.Cache = make(map[string]*AggregativeTaskCacheItem)
	}
	if _, e := q.Cache[msg.TaskKey]; !e {
		q.Cache[msg.TaskKey] = &AggregativeTaskCacheItem{
			Subtasks:   make(map[string]*AggregativeTaskCacheSubtaskItem),
			ReportTo:   msg.ReportTo,
			SubtaskKey: msg.SubtaskKey,
		}
	}
	for _, subtaskKey := range msg.Subtasks {
		q.Cache[msg.TaskKey].Subtasks[subtaskKey] = &AggregativeTaskCacheSubtaskItem{
			SubtaskKey: subtaskKey,
			Result:     nil,
		}
	}
}

func (q *AggregativeTaskCache) IsTaskDone(taskKey string) bool {
	if q == nil || len(q.Cache) == 0 {
		return false
	}
	q.Lock()
	defer q.Unlock()
	if cacheItem, e := q.Cache[taskKey]; !e {
		return false
	} else if len(cacheItem.Subtasks) == 0 {
		return false
	} else {
		for _, subtaskItem := range cacheItem.Subtasks {
			if subtaskItem.Result == nil {
				return false
			}
		}
	}
	return true
}

func (q *AggregativeTaskCache) DoesSubtaskExist(taskKey string, subtaskKey string) bool {
	if q == nil || len(q.Cache) == 0 {
		return false
	}
	q.Lock()
	defer q.Unlock()
	if cacheItem, e := q.Cache[taskKey]; e {
		if len(cacheItem.Subtasks) == 0 {
			return false
		}
		if _, e := cacheItem.Subtasks[subtaskKey]; e {
			return true
		}
	}

	return false
}

func (q *AggregativeTaskCache) GetCumulation(taskKey string) (interface{}, []interface{}) {
	if q == nil || len(q.Cache) == 0 {
		return nil, nil
	}
	q.Lock()
	defer q.Unlock()
	if cacheItem, e := q.Cache[taskKey]; e {
		return cacheItem.Cumulation, q.getResultsOfCompletedSubtasks(taskKey)
	}
	return nil, nil
}

func (q *AggregativeTaskCache) getResultsOfCompletedSubtasks(taskKey string) []interface{} {
	if q == nil || len(q.Cache) == 0 {
		return nil
	}
	if cacheItem, e := q.Cache[taskKey]; e {
		result := []interface{}{}
		for _, subtaskItem := range cacheItem.Subtasks {
			if subtaskItem.Result != nil {
				result = append(result, subtaskItem.Result)
			}
		}
		return result
	}
	return nil
}

func (q *AggregativeTaskCache) SetSubtaskResult(taskKey string, subtaskKey string, result interface{}) {
	if q == nil || len(q.Cache) == 0 {
		return
	}
	q.Lock()
	defer q.Unlock()
	if cacheItem, e := q.Cache[taskKey]; e {
		if subtaskItem, e := cacheItem.Subtasks[subtaskKey]; e {
			subtaskItem.Result = result
		}
	}
}

func (q *AggregativeTaskCache) SetCumulation(taskKey string, cumulation interface{}) {
	if q == nil || len(q.Cache) == 0 {
		return
	}
	q.Lock()
	defer q.Unlock()
	if cacheItem, e := q.Cache[taskKey]; e {
		cacheItem.Cumulation = cumulation
	}
}

func (q *AggregativeTaskCache) CleanTask(taskKey string) {
	if q == nil || len(q.Cache) == 0 {
		return
	}
	q.Lock()
	defer q.Unlock()
	if _, e := q.Cache[taskKey]; e {
		delete(q.Cache, taskKey)
	}
}

func (q *AggregativeTaskCache) GetAggregatorSubtaskKey(taskKey string) string {
	if q == nil || len(q.Cache) == 0 {
		return ""
	}
	q.Lock()
	defer q.Unlock()
	if taskItem, e := q.Cache[taskKey]; e {
		return taskItem.SubtaskKey
	}
	return ""
}

func (q *AggregativeTaskCache) GetReportTo(taskKey string) []*Interface {
	if q == nil || len(q.Cache) == 0 {
		return nil
	}
	q.Lock()
	defer q.Unlock()
	if taskItem, e := q.Cache[taskKey]; e {
		return taskItem.ReportTo
	}
	return nil
}
