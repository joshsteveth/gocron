package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTimeDiff(t *testing.T) {
	//let's say our server is in berlin
	//and our location is in jakarta

	berlin, _ := time.LoadLocation("Europe/Berlin")
	jkt, _ := time.LoadLocation("Asia/Jakarta")

	timeDiff, err := CalculateTimeDiff(time.Now().In(berlin), jkt)
	assert.NoError(t, err)

	assert.Equal(t, time.Hour*-5, timeDiff)
}

func TestGetWeekday(t *testing.T) {
	invalidString := "Monnday"
	validString1 := "saturday"
	validString2 := "Tuesday"

	_, err := getWeekday(invalidString)
	assert.Error(t, err)

	w, err := getWeekday(validString1)
	assert.NoError(t, err)
	assert.Equal(t, 6, int(w))

	w, err = getWeekday(validString2)
	assert.NoError(t, err)
	assert.Equal(t, 2, int(w))
}
