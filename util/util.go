package util

import (
	"strings"
	"sync"
)

// Rotator
type CharRotator struct {
	chars   string
	current []int
	size    uint
	maxChar int
}

func NewCharRotator(chars string, size uint) CharRotator {
	current := make([]int, size)

	return CharRotator{
		chars:   chars,
		size:    size,
		current: current,
		maxChar: len(chars) - 1,
	}
}

func (cr CharRotator) HasNext() bool {
	// check if each char is maxed out
	for _, v := range cr.current {
		if v < cr.maxChar {
			return true
		}
	}
	// no value less than max
	return false
}

func (cr *CharRotator) Inc() bool {
	if cr.HasNext() {
		var i int = int(cr.size) - 1
		for i >= 0 {
			c := cr.current[i]
			if c == cr.maxChar {
				cr.current[i] = 0
				i -= 1
				continue
			} else {
				cr.current[i] += 1
				return true
			}
		}
	}

	return false
}

func (cr *CharRotator) Next() (string, bool) {
	res := cr.CurrentString()
	cont := cr.Inc()
	return res, cont
}

func (cr *CharRotator) CurrentString() string {
	chars := make([]string, cr.size)
	var i uint = 0
	for i < cr.size {
		chars[i] = cr.chars[cr.current[i] : cr.current[i]+1]
		i += 1
	}
	return strings.Join(chars, "")
}

func ForEachCharCombo(charset string, length uint, action func(string)) {
	cr := NewCharRotator(charset, length)
	for s, n := cr.Next(); n; s, n = cr.Next() {
		action(s)
	}
	action(cr.CurrentString())
}

// Url Generator
type UrlGenerator struct {
	base           string
	maxLength      uint
	currentLength  uint
	charset        string
	currentRotator CharRotator
}

func NewUrlGenerator(base string, size uint, charset string) UrlGenerator {
	return UrlGenerator{
		base:           base,
		maxLength:      size,
		currentLength:  1,
		charset:        charset,
		currentRotator: NewCharRotator(charset, 1),
	}
}

func (ug *UrlGenerator) IsDone() bool {
	// if rotator has no next and length is at max
	return !ug.currentRotator.HasNext() && ug.currentLength >= ug.maxLength
}

func (ug *UrlGenerator) Inc() bool {
	if !ug.IsDone() {
		if ug.currentRotator.HasNext() {
			ug.currentRotator.Inc()
		} else {
			ug.currentRotator = NewCharRotator(ug.charset, ug.currentLength+1)
			ug.currentLength += 1
		}
		return true
	}

	return false
}

func (ug *UrlGenerator) Next() (string, bool) {
	u := ug.CurrentUrl()
	n := ug.Inc()
	return u, n
}

func (ug *UrlGenerator) CurrentUrl() string {
	if strings.HasSuffix(ug.base, "/") {
		return ug.base + ug.currentRotator.CurrentString()
	}

	return ug.base + "/" + ug.currentRotator.CurrentString()
}

func ForEachUrlGen(base string, charset string, l uint, action func(string)) {
	ug := NewUrlGenerator(base, l, charset)
	for u, n := ug.Next(); n; u, n = ug.Next() {
		action(u)
	}
	action(ug.CurrentUrl())
}

// Task Scheduler
type WorkerPool struct {
	workers    uint
	taskchan   chan func()
	quit       chan struct{}
	bufferSize uint
	wg         *sync.WaitGroup
}

func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
}

func (wp *WorkerPool) Start() {
	runTask := func(task func(), quit chan struct{}) {
		// channel to track when done
		donechan := make(chan struct{})
		// execute then notify if done
		go func(task func()) {
			task()
			close(donechan)
		}(task)

		// wait for done or quit
		select {
		case <-quit:
			return
		case <-donechan:
			return
		}
	}

	for i := 0; i < int(wp.workers); i++ {
		go func(wp *WorkerPool) {
			for {
				select {
				case f := <-wp.taskchan:
					wp.wg.Add(1)
					runTask(f, wp.quit)
					wp.wg.Done()
				case <-wp.quit:
					return
				}
			}
		}(wp)
	}
}

func NewWorkerPool(bufferSize uint, workers uint) WorkerPool {
	var wg sync.WaitGroup
	return WorkerPool{
		workers:    workers,
		bufferSize: bufferSize,
		taskchan:   make(chan func(), bufferSize),
		quit:       make(chan struct{}),
		wg:         &wg,
	}
}

func (wp *WorkerPool) SubmitTask(f func()) {
	wp.taskchan <- f
}
