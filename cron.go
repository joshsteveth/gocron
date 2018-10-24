package cron

import (
	"context"
	"fmt"
	"log"
	"time"
)

var (
	ErrorAlreadyActive   = fmt.Errorf("Failed to start: CronJob is already active")
	ErrorAlreadyInactive = fmt.Errorf("Failed to stop: CronJob is already inactive")
)

//CronJob is an entity to run
//a certain func in a certain interval
type CronJob struct {
	job func(context.Context)
	//scheduler represents the schedule of the job
	//it should not be able to be modified
	scheduler Scheduler
	//quit channel signals the cronjob to stop doing its job
	quitChan chan interface{}
	jobChan  chan interface{}
	//active bool
}

func NewCronJob(job func(context.Context), s Scheduler) (*CronJob, error) {
	if _, err := s.GetSleepDuration(time.Now()); err != nil {
		return nil, err
	}

	return &CronJob{
		job:       job,
		scheduler: s,
		quitChan:  make(chan interface{}, 1),
		jobChan:   make(chan interface{}, 1),
	}, nil
}

func (c *CronJob) GetSchedule() Scheduler {
	return c.scheduler
}

//Start starts the cronjob in new goroutine
//make sure that this one is not active already
func (c *CronJob) Start() error {
	if len(c.jobChan) == 1 {
		return ErrorAlreadyActive
	}

	go c.start()

	return nil
}

func (c *CronJob) start() {
	sleepDur, err := c.scheduler.GetSleepDuration(time.Now())
	if err != nil {
		log.Printf("Unable to get sleep duration: %v", err)
		return
	}

	c.jobChan <- "start!"

	select {
	case <-time.After(sleepDur):
		//time to start our first job
	case signal := <-c.quitChan:
		//cancel this cronjob and set status into inactive
		log.Printf("CronJob is stopped with signal: %v", signal)
		return
	}

	jobInterval := c.scheduler.GetInterval()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), jobInterval)

		go c.job(ctx)

		select {
		case <-ctx.Done():
			//time to do another job
			continue
		case signal := <-c.quitChan:
			log.Printf("CronJob is stopped with signal: %v", signal)
			cancel()
			return
		}

	}
}

func (c *CronJob) Stop(msg string) error {
	if !c.IsActive() {
		return ErrorAlreadyInactive
	}

	<-c.jobChan
	c.quitChan <- msg

	return nil
}

func (c *CronJob) IsActive() bool {
	if len(c.jobChan) == 1 {
		return true
	}
	return false
}
