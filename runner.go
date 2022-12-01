package llamatask

import (
	"sync"
	"time"
)

// Task is the simple interface used in the Runner.
// the Runner calls Run on each tick
type Task interface {
	Run()
}

// InitilizableTask is a Task that has a initializer (to setup vars, etc)
type InitilizableTask interface {
	Task
	Initialize()
}

// Runner is the main struct used to hold runner's configuration
type Runner struct {
	mut                   sync.Mutex
	ticker                *time.Ticker
	tasks                 []interface{}
	shouldRunOnGoroutines bool
}

// Run simply runs all the tasks.
// NOTE: it blocks the current thread forever if you don't want this
//
//	consider using RunAsync instead
func (r *Runner) Run() { // main runner thread
	for range r.ticker.C { // Run on each tick
		r.mut.Lock()
		for _, task := range r.tasks {
			if r.shouldRunOnGoroutines {
				go task.(Task).Run()
			} else {
				task.(Task).Run()
			}
		}
		r.mut.Unlock()
	}
}

func (r *Runner) RunOnce() {
	r.mut.Lock()
	defer r.mut.Unlock()
	for _, task := range r.tasks {
		if r.shouldRunOnGoroutines {
			go task.(Task).Run()
		} else {
			task.(Task).Run()
		}
	}
}

// RunOnceAsync runs RunOnce in a goroutine
func (r *Runner) RunOnceAsync() {
	go r.RunOnce()
}

// RunAsync runs Run in a goroutine
func (r *Runner) RunAsync() {
	go r.Run()
}

// AddTask adds a task to the Runner. and panics if t is neither Task or InitilizableTask
// NOTE: it blocks until the current iteration of the loop is complete
//
//	if you don't want this use AddTaskAsync instead
func (r *Runner) AddTask(t interface{}) {
	if initilizableTask, ok := t.(InitilizableTask); ok {
		initilizableTask.Initialize()
	} else if _, ok := t.(Task); !ok {
		panic("called AddTask on a task that doesn't implement Task")
	}
	r.mut.Lock()
	defer r.mut.Unlock()
	r.tasks = append(r.tasks, t)
}

// AddTaskAsync runs AddTask in a goroutine
func (r *Runner) AddTaskAsync(t interface{}) {
	go r.AddTask(r)
}

// NewRunner initializes a new Runner
func NewRunner(interval time.Duration, shouldRunOnGoroutines bool) Runner {
	return Runner{
		ticker:                time.NewTicker(interval),
		shouldRunOnGoroutines: shouldRunOnGoroutines,
	}
}
