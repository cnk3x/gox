package cron

import (
	"log/slog"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cnk3x/gox/chans"
)

type Cron struct {
	entries []*Entry
	parser  Parser
	local   *time.Location

	addch    chan *Entry
	closeAdd func()

	delch    chan *Entry
	closeDel func()

	stopch    chan struct{}
	closeStop func()

	id      atomic.Uint64
	running atomic.Bool

	wg sync.WaitGroup
	mu sync.Mutex
}

func New() *Cron {
	c := &Cron{}
	c.local = time.Local
	c.parser = optionalParser
	c.addch, c.closeAdd = chans.MakeChan[*Entry](5)
	c.delch, c.closeDel = chans.MakeChan[*Entry](5)

	c.stopch, c.closeStop = chans.StructChan()
	return c
}

func (c *Cron) Add(name, spec string, jobFn func()) (err error) {
	shedule, e := c.parser.Parse(spec)
	if err = e; err != nil {
		return
	}

	entry := &Entry{ID: c.id.Add(1), Name: name, Cron: spec, Fn: jobFn, Schedule: shedule}

	if c.running.Load() {
		c.addch <- entry
	} else {
		c.entries = append(c.entries, entry)
	}
	return
}

func (c *Cron) Run() {
	if c.running.CompareAndSwap(false, true) {
		c.run()
	}
}

func (c *Cron) now() time.Time {
	now := time.Now()
	if c.local != nil {
		now = now.In(c.local)
	}
	return now
}

func (c *Cron) run() {
	slog.Info("[cron] starting")
	defer func() {
		c.closeStop()
		c.closeAdd()
		slog.Info("[cron] stopped")
	}()

	now := c.now()
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
		slog.Info("[cron] schedule", "name", entry.Name, "next", entry.Next)
	}

	for {
		// Determine the next entry to run.
		slices.SortFunc(c.entries, func(a, b *Entry) int { return a.Next.Compare(b.Next) })

		var timer *time.Timer
		if len(c.entries) == 0 || c.entries[0].Next.IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries and stop requests.
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(c.entries[0].Next.Sub(now))
		}

		select {
		case now = <-timer.C:
			now = now.In(c.local)
			slog.Info("[cron] wake", "now", now)
			// Run every entry whose next time was less than now
			for _, e := range c.entries {
				if e.Next.After(now) || e.Next.IsZero() {
					break
				}
				c.jonRun(e.Fn)
				e.Prev = e.Next
				e.Next = e.Schedule.Next(now)
				if e.Next.IsZero() {
					slog.Info("[cron] job unsatisfiable, remove it", "name", e.Name, "id", e.ID)
					c.delch <- e
				}
				slog.Info("[cron] job run", "name", e.Name, "id", e.ID, "next", e.Next)
			}
		case e := <-c.addch:
			timer.Stop()
			now = c.now()
			e.Next = e.Schedule.Next(now)
			c.entries = append(c.entries, e)
			slog.Info("[cron] job added", "name", e.Name, "id", e.ID, "next", e.Next)
		case e := <-c.delch:
			timer.Stop()
			now = c.now()
			c.entries = slices.DeleteFunc(c.entries, func(item *Entry) bool { return item.ID == e.ID })
			slog.Info("[cron] job deled", "name", e.Name, "id", e.ID, "next", e.Next)
		case <-c.stopch:
			timer.Stop()
			slog.Info("[cron] stop")
			return
		}
	}
}

func (c *Cron) jonRun(j func()) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		func() {
			defer recover()
			j()
		}()
	}()
}

func (c *Cron) Stop() { c.closeStop() }

type Entry struct {
	Name string
	ID   uint64
	Cron string
	Fn   func()

	// Schedule on which this job should be run.
	Schedule Schedule

	// Next time the job will run, or the zero time if Cron has not been
	// started or this entry's schedule is unsatisfiable
	Next time.Time

	// Prev is the last time this job was run, or the zero time if never.
	Prev time.Time
}
