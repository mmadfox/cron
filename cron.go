package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"syscall"
)

type Handle func() error

type Cron struct {
	lock     bool
	schedule string
	handle   Handle
	close    Handle
	sigs     chan os.Signal
	done     chan interface{}
	cr       *cron.Cron
}

func (c *Cron) Handle(fn Handle) {
	c.handle = fn
}

func (c *Cron) CloseHandle(fn Handle) {
	c.close = fn
}

func (c *Cron) process() {
	if c.lock {
		fmt.Println("Skip cron job")
		return
	}
	c.lock = true
	err := c.handle()
	if err != nil {
		c.close()
		c.cr.Stop()
		c.lock = false
		go func() {
			c.done <- true
		}()
		os.Exit(1)
		return
	}
	c.lock = false
}

func (c *Cron) Run() {
	c.process()
	c.cr.AddFunc(c.schedule, func() {
		c.process()
	})
	c.cr.Start()
	c.terminate()
}

func (c *Cron) terminate() {
	signal.Notify(c.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c.sigs
		c.close()
		c.cr.Stop()
		c.done <- true
	}()
	<-c.done
	os.Exit(0)
}

func New(schedule string) (*Cron, error) {
	sigs := make(chan os.Signal, 1)
	done := make(chan interface{})
	c := cron.New()
	return &Cron{
		lock:     false,
		schedule: schedule,
		sigs:     sigs,
		done:     done,
		cr:       c,
	}, nil
}
