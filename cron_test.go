package cron

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestJob() func(context.Context) {
	return func(ctx context.Context) {}
}

func TestNewCronJob(t *testing.T) {
	j := getTestJob()
	invalidScheduler := Hourly{}

	berlin, _ := time.LoadLocation("Europe/Berlin")
	validScheduler, err := NewHourly("15", berlin)

	cj, err := NewCronJob(j, invalidScheduler)
	assert.Error(t, err)

	cj, err = NewCronJob(j, validScheduler)
	assert.NoError(t, err)

	assert.Equal(t, time.Hour, cj.GetSchedule().GetInterval())
}

func TestStart(t *testing.T) {
	j := getTestJob()
	s := CustomInterval{Interval: time.Second * 2}

	cj, err := NewCronJob(j, s)
	assert.NoError(t, err)

	cj.jobChan <- "foo"
	err = cj.Start()
	assert.Error(t, err)

	<-cj.jobChan
	err = cj.Start()
	assert.NoError(t, err)
}

func TestStop(t *testing.T) {
	j := getTestJob()
	s := CustomInterval{Interval: time.Second * 60}

	cj, err := NewCronJob(j, s)
	assert.NoError(t, err)

	msg := "foo"

	err = cj.Start()
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, true, cj.IsActive())

	err = cj.Stop(msg)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, false, cj.IsActive())

	assert.Error(t, cj.Stop(msg))
}
