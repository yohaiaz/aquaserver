package main

import (
	"fmt"
	"path"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
)

type StatisticsStore struct {
	totalFiles int64
	totalSize  int64
	maxSise    int64
	extList    map[string]int

	mutex sync.Mutex
}

var instance *StatisticsStore
var once sync.Once

func GetStatisticsStore() *StatisticsStore {
	once.Do(func() {
		instance = &StatisticsStore{
			mutex:   sync.Mutex{},
			extList: make(map[string]int),
		}
	})
	return instance
}

func (s *StatisticsStore) Add(fileName string, size int64) {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	atomic.AddInt64(&(s.totalFiles), 1)

	atomic.AddInt64(&(s.totalSize), size)

	if size > atomic.LoadInt64(&(s.maxSise)) {
		atomic.StoreInt64(&(s.maxSise), size)
	}

	s.extList[path.Ext(fileName)] += 1
}

func (s *StatisticsStore) Print() string {

	s.mutex.Lock()

	defer s.mutex.Unlock()

	res := "statistics: \n\n"

	res += fmt.Sprintf("Number of files received: %d \n", atomic.LoadInt64(&(s.totalFiles)))

	res += fmt.Sprintf("Maximum file size: %d \n", atomic.LoadInt64(&(s.maxSise)))

	if atomic.LoadInt64(&(s.totalFiles)) == 0 {
		res += fmt.Sprintf("Average file Size: NA \n")
	} else {
		avg := atomic.LoadInt64(&(s.totalSize)) / atomic.LoadInt64(&(s.totalFiles))
		res += fmt.Sprintf("Average file Size: %d \n", avg)
	}

	ext := reflect.ValueOf(s.extList).MapKeys()
	res += fmt.Sprintf("List of extensions: %s \n", ext)

	var es entries
	for k, v := range s.extList {
		es = append(es, entry{val: v, key: k})
	}

	sort.Sort(sort.Reverse(es))

	topExt := ""
	for _, e := range es {
		topExt += fmt.Sprintf("%s (%d), ", e.key, e.val)
	}

	res += fmt.Sprintf("Most frequent file extensions: %s\n", topExt)

	return res
}
