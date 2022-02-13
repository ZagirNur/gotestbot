package util

import "time"

func MonthDaysCount(year, month int) int {

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return lastOfMonth.Day()
}
