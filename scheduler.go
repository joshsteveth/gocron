package cron

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

var (
	ErrorInvalidStartingPoint = fmt.Errorf("Starting point string is not valid")
	ErrorInvalidTimeLocation  = fmt.Errorf("Time location is not valid")
)

//Scheduler interface is used as part of CronJob property
type Scheduler interface {
	//GetInterval retrieves the interval set
	GetInterval() time.Duration
	//GetSleepDuration returns the time duration needed
	//before executing the first job
	GetSleepDuration(time.Time) (time.Duration, error)
}

//CustomInterval
//runs immediately with a custom interval
type CustomInterval struct {
	Interval time.Duration
}

func (ci CustomInterval) GetInterval() time.Duration {
	return ci.Interval
}

func (ci CustomInterval) GetSleepDuration(t time.Time) (time.Duration, error) {
	return time.Second * 0, nil
}

//Hourly
//this scheduler runs every hour
//starting point should be a string between "00" and "59" (minute)
//use time location to make sure it runs correctly in your time zone
type Hourly struct {
	StartingPoint string
	Location      *time.Location
}

func NewHourly(start string, loc *time.Location) (*Hourly, error) {
	h := Hourly{
		StartingPoint: start,
		Location:      loc,
	}

	if err := h.validate(); err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *Hourly) validate() error {
	if _, err := time.Parse("04", h.StartingPoint); err != nil {
		return ErrorInvalidStartingPoint
	}

	if h.Location == nil {
		return ErrorInvalidTimeLocation
	}

	return nil
}

func (h Hourly) GetInterval() time.Duration {
	return time.Hour
}

func (h Hourly) GetSleepDuration(t time.Time) (time.Duration, error) {
	var result time.Duration

	if err := h.validate(); err != nil {
		return result, err
	}

	t = t.In(h.Location)

	//wait precision is in seconds
	//get both number of minute and second in float
	tMinute, tSecond := t.Format("04"), t.Format("05")
	min, err := strconv.ParseFloat(tMinute, 64)
	if err != nil {
		return result, err
	}
	sec, err := strconv.ParseFloat(tSecond, 64)
	if err != nil {
		return result, err
	}

	totalSeconds := min*60 + sec

	startMin, err := strconv.ParseFloat(h.StartingPoint, 64)
	if err != nil {
		return result, err
	}
	startingSeconds := startMin * 60
	if startingSeconds < totalSeconds {
		startingSeconds += 60 * 60
	}

	//calculate the difference between
	//the starting point and total seconds in seconds
	duration := math.Abs(startingSeconds - totalSeconds)

	return time.ParseDuration(fmt.Sprintf("%.0fs", duration))
}

//Daily
//this scheduler runs every other day
//starting point should be parseable into a valid hhmm format
//for instance 1530 for 3:30 pm
type Daily struct {
	StartingPoint string
	Location      *time.Location
}

func NewDaily(start string, loc *time.Location) (*Daily, error) {
	d := Daily{
		StartingPoint: start,
		Location:      loc,
	}

	if err := d.validate(); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *Daily) validate() error {
	if _, err := time.Parse("1504", d.StartingPoint); err != nil {
		return ErrorInvalidStartingPoint
	}

	if d.Location == nil {
		return ErrorInvalidTimeLocation
	}

	return nil
}

func (d Daily) GetInterval() time.Duration {
	return time.Hour * 24
}

func (d Daily) GetSleepDuration(t time.Time) (time.Duration, error) {
	if err := d.validate(); err != nil {
		return time.Second * 0, err
	}

	timeThen, err := time.Parse("20060102 -0700",
		t.Format("20060102 -0700"))
	if err != nil {
		return time.Second * 0, err
	}

	dur, err := time.ParseDuration(fmt.Sprintf("%sh%sm",
		d.StartingPoint[:2], d.StartingPoint[2:]))
	if err != nil {
		return time.Second * 0, err
	}

	timeDiff, err := CalculateTimeDiff(t, d.Location)
	if err != nil {
		return time.Second * 0, err
	}

	timeThen = timeThen.Add(dur).In(d.Location).Add(timeDiff)

	//we need to add 1 day if timeThen is before our start time
	if t.After(timeThen) {
		timeThen = timeThen.AddDate(0, 0, 1)
	}

	return timeThen.Sub(t), nil
}
