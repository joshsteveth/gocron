[![Go Report Card](https://goreportcard.com/badge/github.com/joshsteveth/gocron)](https://goreportcard.com/report/github.com/joshsteveth/gocron)

### Go Cron
Cronjob like scheduler. With current features:
- custom interval 
- hourly
- daily

#### Scheduler interface
Scheduler interface requires 2 methods:
- `GetInterval() time.Duration` : interval between jobs
- `GetSleepDuration(time.Time) (time.Duration, error)`: how long do we have to wait from a certain time before executing the first job


Available build in Scheduler:
- `CustomInterval` : define your own interval and sleep duration
- `Hourly` : runs every hour `h := NewHourly("30")` 
- `Daily` : runs every day `d := NewDaily("1530")`
- `Weekly` : runs every week `w := NewWeekly("Monday@1530")`

`*time.Location` can be added as extra argument for `NewHourly`, `NewDaily`, and `NewWeekly`

#### CronJob
Initialize new cronjob. For example:
```
berlin, _ := time.LoadLocation("Europe/Berlin")
scheduler, err := NewDaily("1530", berlin)
if err != nil{
	log.Fatalf("Failed to initiate a Scheduler: %v", err)
}

testJob := func(ctx context.Context){fmt.Println("test")}

c, err := NewCronJob(testJob, scheduler)
if err != nil{
	log.Fatalf("Failed to create new cron job: %v", err)
}

//starting the cronjob
//it will return an error if it has been started before
if err := c.Start(); err != nil{
	log.Fatalf("Failed to start the cron job: %v", err)
}

//checking the cronjob status
//return true if it's currently active
isActive := c.IsActive()

//stopping the cronjob
//it will return an error if it is not active already
if isActive{
	c.Stop()
}
```
