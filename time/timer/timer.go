package timer

import (
	"container/heap"
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// Forever
const Forever = -1

// precision represents timer precision
const precision = time.Millisecond

// ID represents ID of timer task
type ID int64

func (id ID) Valid() bool { return id > 0 }

// Task represents timer task
type Task interface {
	Exec(ID)
}

// TaskFunc wraps function as a task
type TaskFunc func(ID)

// Exec implements Task Exec method
func (fn TaskFunc) Exec(id ID) { fn(id) }

// Scheduler schedules timers
type Scheduler interface {
	// Start starts the Scheduler
	Start()
	// Shutdown shutdowns the Scheduler
	Shutdown()
	// Add adds a new timer task
	Add(next, duration time.Duration, task Task, times int) ID
	// Remove removes a timer task by ID
	Remove(id ID)
}

type timer struct {
	id       ID
	task     Task
	times    int
	next     int64
	duration int64
}

type timers struct {
	next   int64
	timers []timer
}

func (ts *timers) remove(id ID) bool {
	n := len(ts.timers)
	for i := 0; i < n; i++ {
		if ts.timers[i].id == id {
			copy(ts.timers[i:n-1], ts.timers[i+1:])
			ts.timers[n-1] = timer{}
			ts.timers = ts.timers[:n-1]
			return true
		}
	}
	return false
}

// memoryScheduler implements Scheduler in memory
type memoryScheduler struct {
	allTimers []timers
	indices   map[int64]int // next => indexof(groups)
	timers    map[ID]int64  // id => next

	update chan timer
	quit   chan struct{}
	wait   chan struct{}

	nextId  int64
	running int32

	queue  *list.List
	locker sync.Mutex
	cond   *sync.Cond
}

// NewMemoryScheduler creates in-memory Scheduler
func NewMemoryScheduler() Scheduler {
	s := &memoryScheduler{
		indices: make(map[int64]int),
		timers:  make(map[ID]int64),
		update:  make(chan timer, 128),
		quit:    make(chan struct{}),
		wait:    make(chan struct{}),
		queue:   list.New(),
	}
	s.cond = sync.NewCond(&s.locker)
	return s
}

// Start implements Scheduler Start method
func (s *memoryScheduler) Start() {
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return
	}
	go s.receive()
	go s.schedule()
}

// Shutdown implements Scheduler Shutdown method
func (s *memoryScheduler) Shutdown() {
	if atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		s.cond.Signal()
		close(s.quit)
		<-s.wait
	}
}

func (s *memoryScheduler) receive() {
	for atomic.LoadInt32(&s.running) == 1 {
		s.cond.L.Lock()
		for s.queue.Len() == 0 {
			s.cond.Wait()
		}
		front := s.queue.Front()
		task := front.Value.(timer)
		s.queue.Remove(front)
		s.cond.L.Unlock()
		s.update <- task
	}
}

func (s *memoryScheduler) schedule() {
	var timer *time.Timer
	for {
		if len(s.allTimers) == 0 {
			select {
			case x := <-s.update:
				if x.id < 0 {
					s.removeTimer(-x.id)
				} else {
					s.addTimer(x)
				}
			case <-s.quit:
				close(s.wait)
				return
			}
		}
		now := time.Duration(time.Now().UnixNano()) / time.Nanosecond / precision * precision

		first := heap.Pop(s).(timers)
		next := first.next
		dt := time.Duration(next)*precision - now

		if dt <= 0 {
			s.execGroup(first)
			continue
		}

		if timer == nil {
			timer = time.NewTimer(dt)
		} else {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(dt)
		}

	WAIT_DO_FIRST:
		for {
			select {
			case <-timer.C:
				s.execGroup(first)
				break WAIT_DO_FIRST
			case x := <-s.update:
				if x.id < 0 {
					next, ok := s.timers[-x.id]
					if ok {
						if next == first.next {
							delete(s.timers, -x.id)
							first.remove(-x.id)
						} else {
							s.removeTimerByNext(-x.id, next)
						}
					}
				} else {
					if x.next == first.next {
						first.timers = append(first.timers, x)
					} else {
						s.addTimer(x)
						if x.next < first.next {
							heap.Push(s, first)
							break WAIT_DO_FIRST
						}
					}
				}
			case <-s.quit:
				close(s.wait)
				return
			}
		}
	}
}

// Add implements Scheduler Add method
func (s *memoryScheduler) Add(next, duration time.Duration, task Task, times int) ID {
	id := ID(atomic.AddInt64(&s.nextId, 1))

	s.locker.Lock()
	t := timer{
		id:       id,
		times:    times,
		task:     task,
		duration: int64(duration / time.Millisecond),
		next:     int64(next / time.Millisecond),
	}
	s.queue.PushBack(t)
	l := s.queue.Len()
	s.locker.Unlock()

	if l == 1 {
		s.cond.Signal()
	}

	return id
}

// Remove implements Scheduler Remove method
func (s *memoryScheduler) Remove(id ID) {
	s.locker.Lock()
	t := timer{
		id: -id,
	}
	s.queue.PushBack(t)
	l := s.queue.Len()
	s.locker.Unlock()

	if l == 1 {
		s.cond.Signal()
	}
}

func (s *memoryScheduler) addTimer(x timer) {
	s.timers[x.id] = x.next
	if i, ok := s.indices[x.next]; ok {
		s.allTimers[i].timers = append(s.allTimers[i].timers, x)
	} else {
		g := timers{
			next:   x.next,
			timers: make([]timer, 0, 8),
		}
		g.timers = append(g.timers, x)
		heap.Push(s, g)
	}
}

func (s *memoryScheduler) removeTimer(id ID) {
	next, ok := s.timers[id]
	if !ok {
		return
	}
	s.removeTimerByNext(id, next)
}

func (s *memoryScheduler) removeTimerByNext(id ID, next int64) {
	delete(s.timers, id)

	i, ok := s.indices[next]
	if !ok {
		return
	}

	if s.allTimers[i].remove(id) {
		if len(s.allTimers[i].timers) == 0 {
			heap.Remove(s, i)
		}
	}
}

func (s *memoryScheduler) execGroup(g timers) {
	n := 0
	for i := range g.timers {
		g.timers[i].task.Exec(g.timers[i].id)
		if g.timers[i].times > 0 {
			g.timers[i].times--
		}
		if g.timers[i].times != 0 {
			g.timers[i].next += g.timers[i].duration
			if i != n {
				g.timers[n] = g.timers[i]
				n++
			}
		} else {
			delete(s.timers, g.timers[i].id)
		}
	}
	g.timers = g.timers[:n]
	for i := range g.timers {
		s.addTimer(g.timers[i])
	}
}

// Len implements heap.Interface Len method
func (s *memoryScheduler) Len() int { return len(s.allTimers) }

// Less implements heap.Interface Less method
func (s *memoryScheduler) Less(i, j int) bool { return s.allTimers[i].next < s.allTimers[j].next }

// Swap implements heap.Interface Swap method
func (s *memoryScheduler) Swap(i, j int) {
	s.allTimers[i], s.allTimers[j] = s.allTimers[j], s.allTimers[i]
	s.indices[s.allTimers[i].next] = i
	s.indices[s.allTimers[j].next] = j
}

// Push implements heap.Interface Push method
func (s *memoryScheduler) Push(x any) {
	g := x.(timers)
	l := len(s.allTimers)
	s.allTimers = append(s.allTimers, g)
	s.indices[g.next] = l
}

// Pop implements heap.Interface Pop method
func (s *memoryScheduler) Pop() any {
	l := len(s.allTimers)
	x := s.allTimers[l-1]
	s.allTimers[l-1] = timers{}
	s.allTimers = s.allTimers[:l-1]
	delete(s.indices, x.next)
	return x
}

// SetTimeout add a timeout timer
func SetTimeout(scheduler Scheduler, d time.Duration, task Task) ID {
	return scheduler.Add(time.Duration(time.Now().Add(d).UnixNano()), d, task, 1)
}

// SetTimeoutFunc add a timeout timer func
func SetTimeoutFunc(scheduler Scheduler, d time.Duration, fn TaskFunc) ID {
	return SetTimeout(scheduler, d, fn)
}

// ClearTimeout removes the timeout timer by ID
func ClearTimeout(scheduler Scheduler, id ID) {
	scheduler.Remove(id)
}

// SetInterval add interval timer
func SetInterval(scheduler Scheduler, d time.Duration, task Task) ID {
	return scheduler.Add(time.Duration(time.Now().Add(d).UnixNano()), d, task, Forever)
}

// SetIntervalFunc add a interval timer func
func SetIntervalFunc(scheduler Scheduler, d time.Duration, fn TaskFunc) ID {
	return SetInterval(scheduler, d, fn)
}

// ClearTimeout removes the interval timer by ID
func ClearInterval(scheduler Scheduler, id ID) {
	scheduler.Remove(id)
}
