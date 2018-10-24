package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomInterval(t *testing.T) {
	interval := time.Second * 3
	sleepDuration := time.Millisecond * 500
	c := CustomInterval{
		Interval:      interval,
		SleepDuration: sleepDuration,
	}

	sleepDur, err := c.GetSleepDuration(time.Now())
	assert.NoError(t, err)

	assert.Equal(t, interval, c.GetInterval())
	assert.Equal(t, sleepDuration, sleepDur)
}

func TestHourly(t *testing.T) {
	invalidStartingPoint := "60"
	validStartingPoint := "30"

	loc, _ := time.LoadLocation("Europe/Berlin")

	_, err := NewHourly(invalidStartingPoint, loc)
	assert.Error(t, err)

	_, err = NewHourly(validStartingPoint, nil)
	assert.Error(t, err)

	h, err := NewHourly(validStartingPoint)
	assert.NoError(t, err)

	h, err = NewHourly(validStartingPoint, loc)
	assert.NoError(t, err)

	assert.Equal(t, time.Second*3600, h.GetInterval())

	startTime, _ := time.Parse("20060102 1504", "20181010 0010")
	timeDur, err := h.GetSleepDuration(startTime)
	assert.NoError(t, err)
	assert.Equal(t, float64(20*60), timeDur.Seconds())

	startTime, _ = time.Parse("20060102 1504", "20181010 1550")
	timeDur, err = h.GetSleepDuration(startTime)
	assert.NoError(t, err)
	assert.Equal(t, float64(40*60), timeDur.Seconds())
}

func TestDaily(t *testing.T) {
	invalidStartingPoint := "2530"
	validStartingPoint := "1530"

	berlin, _ := time.LoadLocation("Europe/Berlin")

	_, err := NewDaily(invalidStartingPoint, berlin)
	assert.Error(t, err)

	_, err = NewDaily(validStartingPoint, nil)
	assert.Error(t, err)

	d, err := NewDaily(validStartingPoint)
	assert.NoError(t, err)

	d, err = NewDaily(validStartingPoint, berlin)
	assert.NoError(t, err)

	assert.Equal(t, time.Hour*24, d.GetInterval())

	//this runs in TZ UTC+7
	startTime, _ := time.Parse("20060102 1504 -0700", "20181010 1750 +0700")
	timeDur, err := d.GetSleepDuration(startTime)
	assert.NoError(t, err)
	assert.Equal(t, float64(2*3600+40*60), timeDur.Seconds())

	//already pass a day
	//2210 in UTC+7 means 1710 un UTC+2
	startTime, _ = time.Parse("20060102 1504 -0700", "20181010 2210 +0700")
	timeDur, err = d.GetSleepDuration(startTime)
	assert.NoError(t, err)
	assert.Equal(t, float64(6*3600+50*60+15*3600+30*60), timeDur.Seconds())
}

func TestWeekly(t *testing.T) {
	invalidWeek := "monnday"
	validWeek := "wednesday"
	invalidStartingPoint := "2530"
	validStartingPoint := "1530"

	correctFormat := func(w, sp string) string {
		return fmt.Sprintf("%s@%s", w, sp)
	}

	falseFormat := func(w, sp string) string {
		return fmt.Sprintf("%s-%s", w, sp)
	}

	berlin, _ := time.LoadLocation("Europe/Berlin")

	//invalid week correct format
	_, err := NewWeekly(correctFormat(invalidWeek, validStartingPoint), berlin)
	assert.Error(t, err)

	//invalid starting point correct format
	_, err = NewWeekly(correctFormat(validWeek, invalidStartingPoint), berlin)
	assert.Error(t, err)

	//false format
	_, err = NewWeekly(falseFormat(validWeek, validStartingPoint), berlin)
	assert.Error(t, err)

	//correct format nil location
	_, err = NewWeekly(correctFormat(validWeek, validStartingPoint), nil)
	assert.Error(t, err)

	//no location is fine
	w, err := NewWeekly(correctFormat(validWeek, validStartingPoint))
	assert.NoError(t, err)

	//correct format valid location
	w, err = NewWeekly(correctFormat(validWeek, validStartingPoint), berlin)
	assert.NoError(t, err)

	//this runs in TZ UTC+7 on Tuesday at 17:50
	startTime, _ := time.Parse("20060102 1504 -0700", "20181023 1750 +0700")
	timeDur, err := w.GetSleepDuration(startTime)
	assert.NoError(t, err)
	//job should run every Wednesday at 15:30 in berlin
	//that means Wednesday 20:30 in UTC+7
	//so duration is 1 day 2 hours and 40 minutes
	assert.Equal(t, float64(26*3600+40*60), timeDur.Seconds())

	//this runs in TZ UTC+7 on Saturday at 23:35
	startTime, _ = time.Parse("20060102 1504 -0700", "20181020 2335 +0700")
	timeDur, err = w.GetSleepDuration(startTime)
	assert.NoError(t, err)
	//job should run every Wednesday at 15:30 in berlin
	//that means Wednesday 20:30 in UTC+7
	//so duration is 3 day 20 hours and 55 minutes
	assert.Equal(t, float64(3*24*3600+20*3600+55*60), timeDur.Seconds())
}
