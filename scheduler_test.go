package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomInterval(t *testing.T) {
	interval := time.Second * 3
	c := CustomInterval{
		Interval: interval,
	}

	sleepDur, err := c.GetSleepDuration(time.Now())
	assert.NoError(t, err)

	assert.Equal(t, interval, c.GetInterval())
	assert.Equal(t, time.Second*0, sleepDur)
}

func TestHourly(t *testing.T) {
	invalidStartingPoint := "60"
	validStartingPoint := "30"

	loc, _ := time.LoadLocation("Europe/Berlin")

	_, err := NewHourly(invalidStartingPoint, loc)
	assert.Error(t, err)

	_, err = NewHourly(validStartingPoint, nil)
	assert.Error(t, err)

	h, err := NewHourly(validStartingPoint, loc)
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

	d, err := NewDaily(validStartingPoint, berlin)
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
