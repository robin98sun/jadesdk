package jadesdk

type AggregativeTaskCache struct {
	Cache map[string]*AggregativeTaskCacheItem // taskId: taskQueueItem
}

func NewAggregativeTaskCache() *AggregativeTaskCache {
	inst := &AggregativeTaskCache{
		Cache: make(map[string]*AggregativeTaskCacheItem),
	}
	return inst
}

type AggregativeTaskCacheItem struct {
	Subtasks   map[string]*AggregativeTaskCacheSubtaskItem // subtaskId: subtaskItem
	Cumulation interface{}                                 // cumulative result of so far responded subtasks
	ReportTo   []*Interface                                // upper layer aggregators
}

type AggregativeTaskCacheSubtaskItem struct {
	SubtaskKey string
	result     interface{}
}

func (q *AggregativeTaskCache) EnqueueAggregativeTask(msg *AggregatorEnqueuingMessage) {
	if q == nil || msg == nil || msg.TaskKey == "" || len(msg.Subtasks) == 0 {
		return
	}
	if q.Cache == nil {
		q.Cache = make(map[string]*AggregativeTaskCacheItem)
	}
	if _, e := q.Cache[msg.TaskKey]; !e {
		q.Cache[msg.TaskKey] = &AggregativeTaskCacheItem{
			Subtasks: make(map[string]*AggregativeTaskCacheSubtaskItem),
			ReportTo: msg.ReportTo,
		}
	}
	for _, subtaskKey := range msg.Subtasks {
		q.Cache[msg.TaskKey].Subtasks[subtaskKey] = &AggregativeTaskCacheSubtaskItem{
			SubtaskKey: subtaskKey,
		}
	}
}

func (q *AggregativeTaskCache) IsTaskDone(taskKey string) bool {
	if q == nil || len(q.Cache) == 0 {
		return false
	}
	if cacheItem, e := q.Cache[taskKey]; !e {
		return false
	} else if len(cacheItem.Subtasks) == 0 {
		return false
	} else {
		for _, subtaskItem := range cacheItem.Subtasks {
			if subtaskItem.result == nil {
				return false
			}
		}
	}
	return true
}

func (q *AggregativeTaskCache) DoesSubtaskExists(taskKey string, subtaskKey string) bool {
	if q == nil || len(q.Cache) == 0 {
		return false
	}
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
	if cacheItem, e := q.Cache[taskKey]; e {
		return cacheItem.Cumulation, q.GetResultsOfCompletedSubtasks(taskKey)
	}
	return nil, nil
}

func (q *AggregativeTaskCache) GetResultsOfCompletedSubtasks(taskKey string) []interface{} {
	if q == nil || len(q.Cache) == 0 {
		return nil
	}
	if cacheItem, e := q.Cache[taskKey]; e {
		result := []interface{}{}
		for _, subtaskItem := range cacheItem.Subtasks {
			result = append(result, subtaskItem.result)
		}
		return result
	}
	return nil
}

func (q *AggregativeTaskCache) SetSubtaskResult(taskKey string, subtaskKey string, result interface{}) {
	if q == nil || len(q.Cache) == 0 {
		return
	}
	if cacheItem, e := q.Cache[taskKey]; e {
		if subtaskItem, e := cacheItem.Subtasks[subtaskKey]; e {
			subtaskItem.result = result
		}
	}
}

func (q *AggregativeTaskCache) SetCumulation(taskKey string, cumulation interface{}) {
	if q == nil || len(q.Cache) == 0 {
		return
	}
	if cacheItem, e := q.Cache[taskKey]; e {
		cacheItem.Cumulation = cumulation
	}
}
