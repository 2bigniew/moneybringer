package TimeUtils

import (
	"fmt"
	"time"
)

func GetCurrentTime() time.Time {
	currentTime := time.Now()
	return currentTime
}

func GetPreviousMonthTime() time.Time {
	currentTime := time.Now()
	previousMonthTime := currentTime.AddDate(0, -1, 0)
	return previousMonthTime
}

func FormatToDdMmYyyy(timeObj time.Time) string {
	formated := timeObj.Format("02-01-2006")
	return formated
}

func SetDayOfMonth(timeObj time.Time, dayNumber int) time.Time {
	updatedTime := time.Date(
		timeObj.Year(),
		timeObj.Month(),
		dayNumber,
		timeObj.Hour(),
		timeObj.Minute(),
		timeObj.Second(),
		timeObj.Nanosecond(),
		timeObj.Location())

	return updatedTime
}

func GetDataFromDdMmYyyyFormat(dateString string) (time.Time, bool) {
	layout := "02-01-2006"

	date, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Now(), false
	}

	return date, true
}
