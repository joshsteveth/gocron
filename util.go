package cron

import (
	"fmt"
	"time"
)

//CalculateTimeDiff
//returns time difference between selected time location and server local time
//for instance if the server runs in UTC and loc is UTC+7
//then this function should yield time.Hour * 7 as result
func CalculateTimeDiff(serverTime time.Time, loc *time.Location) (time.Duration, error) {
	getSeconds := func(t time.Time) (float64, error) {
		tz := t.Format("-0700")
		dur, err := time.ParseDuration(fmt.Sprintf("%sh%sm",
			tz[1:3], tz[3:5]))
		if err != nil {
			return 0, err
		}

		sec := dur.Seconds()
		if tz[:1] == "-" {
			sec = sec * -1
		}

		return sec, nil
	}

	serverSec, err := getSeconds(serverTime)
	if err != nil {
		return time.Second * 0, err
	}

	localSec, err := getSeconds(serverTime.In(loc))
	if err != nil {
		return time.Second * 0, err
	}

	return time.ParseDuration(fmt.Sprintf("%.0fs",
		serverSec-localSec))
}
