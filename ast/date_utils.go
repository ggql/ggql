package ast

import (
	"strconv"
	"strings"
	"time"
)

const (
	ChronoTimeFormat         = "15:04:05"
	ChronoDateFormat         = "2006-01-02"
	ChronoDateTimeFormat     = "2006-01-02 15:04:05"
	ChronoDateTimeFullFormat = "2006-01-02 15:04:05.000"
)

// nolint:gomnd
func GetUnixTimestampMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func TimeStampToDate(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)

	return t.Format(ChronoDateFormat)
}

func TimeStampToTime(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)

	return t.Format(ChronoTimeFormat)
}

func TimeStampToDateTime(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)

	return t.Format(ChronoDateTimeFormat)
}

func DateToTimeStamp(date string) int64 {
	t, _ := time.Parse(ChronoDateFormat, date)

	return t.Unix()
}

func DateTimeToTimeStamp(date string) int64 {
	data_format := ChronoDateTimeFormat

	if strings.Contains(date, ".") {
		data_format = ChronoDateTimeFullFormat
	}

	t, _ := time.Parse(data_format, date)

	return t.Unix()
}

func DateTimeToHour(date int64) int64 {
	t := time.Unix(date, 0)

	return int64(t.Hour())
}

func DateToDayNumberInMonth(date int64) int {
	t := time.Unix(date, 0)

	return t.Day()
}

func DateToDayName(date int64) string {
	t := time.Unix(date, 0)
	dayName := ""

	switch t.Weekday() {
	case time.Monday:
		dayName = "Monday"
	case time.Tuesday:
		dayName = "Tuesday"
	case time.Wednesday:
		dayName = "Wednesday"
	case time.Thursday:
		dayName = "Thursday"
	case time.Friday:
		dayName = "Friday"
	case time.Saturday:
		dayName = "Saturday"
	case time.Sunday:
		dayName = "Sunday"
	}

	return dayName
}

func DateToMonthName(date int64) string {
	t := time.Unix(date, 0)
	monthName := ""

	switch t.Month() {
	case time.January:
		monthName = "January"
	case time.February:
		monthName = "February"
	case time.March:
		monthName = "March"
	case time.April:
		monthName = "April"
	case time.May:
		monthName = "May"
	case time.June:
		monthName = "June"
	case time.July:
		monthName = "July"
	case time.August:
		monthName = "August"
	case time.September:
		monthName = "September"
	case time.October:
		monthName = "October"
	case time.November:
		monthName = "November"
	case time.December:
		monthName = "December"
	}

	return monthName
}

func TimeStampFromYearAndDay(year int, dayOfYear uint) int64 {
	t := time.Date(year, 1, 0, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, int(dayOfYear))

	return t.Unix()
}

// nolint:gomnd
func IsValidTimeFormat(timeStr string) bool {
	var milliseconds int

	// Check length of the string
	if len(timeStr) < 8 || len(timeStr) > 12 {
		return false
	}

	// Split the string into hours, minutes, seconds, and optional milliseconds
	parts := strings.Split(timeStr, ":")
	if len(parts) < 3 || len(parts) > 4 {
		return false
	}

	// Extract hours, minutes, seconds, and optionally milliseconds
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	secondsParts := strings.Split(parts[2], ".")
	seconds, _ := strconv.Atoi(secondsParts[0])

	if len(secondsParts) == 2 {
		milliseconds, _ = strconv.Atoi(secondsParts[1])
	} else {
		milliseconds = 0
	}

	// Validate the parsed values
	return hours >= 0 && hours < 24 &&
		minutes >= 0 && minutes < 60 &&
		seconds >= 0 && seconds < 60 &&
		milliseconds >= 0 && milliseconds < 1000
}

// nolint:gomnd
func IsValidDateFormat(dateStr string) bool {
	// Check length of the string
	if len(dateStr) != 10 {
		return false
	}

	// Split the string into year, month, and day
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return false
	}

	// Extract year, month, and day
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])

	// Validate the parsed values
	return year >= 1 && month >= 1 && month <= 12 && day >= 1 && day <= 31
}

// nolint:gomnd
func IsValidDateTimeFormat(datetimeStr string) bool {
	// Check length of the string
	if len(datetimeStr) < 19 || len(datetimeStr) > 23 {
		return false
	}

	// Split the string into date and time components
	parts := strings.Fields(datetimeStr)
	if len(parts) != 2 {
		return false
	}

	// Check the validity of date and time components
	return IsValidDateFormat(parts[0]) && IsValidTimeFormat(parts[1])
}
