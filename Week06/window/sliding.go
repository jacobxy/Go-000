package window

import (
	"sync"
	"time"
)

type Elem struct {
	Cnt int
}

type Sliding struct {
	Elems    map[int64]*Elem
	Mutex    sync.RWMutex
	Interval int64
	AvgMax   float64
}

func NewSliding(interval int64, avgmax float64) *Sliding {
	return &Sliding{
		Elems:    make(map[int64]*Elem, 10),
		Interval: interval,
		AvgMax:   avgmax,
	}
}

func (sd *Sliding) getCurElem() *Elem {
	now := time.Now().Unix()
	elem, ok := sd.Elems[now]
	if !ok {
		elem = &Elem{}
		sd.Elems[now] = elem
	}
	return elem
}

func (sd *Sliding) Clear() {
	start := time.Now().Unix() - sd.Interval
	for k := range sd.Elems {
		if k < start {
			delete(sd.Elems, k)
		}
	}
}

func (sd *Sliding) increment() {
	sd.Mutex.Lock()
	defer sd.Mutex.Unlock()

	v := sd.getCurElem()
	v.Cnt += 1
	sd.Clear()
}

func (sd *Sliding) avg() float64 {
	sum := int(0)
	start := time.Now().Unix() - sd.Interval
	sd.Mutex.RLock()
	defer sd.Mutex.RUnlock()
	for k, v := range sd.Elems {
		if k > start {
			sum += v.Cnt
		}
	}
	return float64(sum) / float64(sd.Interval)
}

func (sd *Sliding) Allow() (float64, bool) {
	avg := sd.avg()
	if avg > sd.AvgMax {
		return avg, false
	}
	sd.increment()
	return avg, true
}
