package ast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestTimestamp    = 1705117592
	TestDate         = "2024-01-10"
	TestDateTime     = "2024-01-10 12:36:31"
	TestDateTimeFull = "2024-01-10 12:36:31.000"
)

func TestGetUnixTimestampMs(t *testing.T) {
	ret := GetUnixTimestampMs()
	fmt.Println("GetUnixTimestampMs:", ret)
	assert.Greater(t, ret, int64(0))
}

func TestTimeStampToDate(t *testing.T) {
	ret := TimeStampToDate(TestTimestamp)
	fmt.Println("TimeStampToDate:", ret)
	assert.NotEqual(t, "", ret)
}

func TestTimeStampToTime(t *testing.T) {
	ret := TimeStampToTime(TestTimestamp)
	fmt.Println("TimeStampToTime:", ret)
	assert.NotEqual(t, "", ret)
}

func TestTimeStampToDateTime(t *testing.T) {
	ret := TimeStampToDateTime(TestTimestamp)
	fmt.Println("TimeStampToDateTime:", ret)
	assert.NotEqual(t, "", ret)
}

func TestDateToTimeStamp(t *testing.T) {
	ret := DateToTimeStamp(TestDateTime)
	fmt.Println("DateToTimeStamp:", ret)
	assert.LessOrEqual(t, ret, int64(0))

	ret = DateToTimeStamp(TestDate)
	fmt.Println("DateToTimeStamp:", ret)
	assert.Greater(t, ret, int64(0))
}

func TestDateTimeToTimeStamp(t *testing.T) {
	ret := DateTimeToTimeStamp(TestDateTimeFull)
	fmt.Println("DateTimeToTimeStamp:", ret)
	assert.Greater(t, ret, int64(0))

	ret = DateTimeToTimeStamp(TestDateTime)
	fmt.Println("DateTimeToTimeStamp:", ret)
	assert.Greater(t, ret, int64(0))

	ret = DateTimeToTimeStamp("invalid")
	fmt.Println("DateTimeToTimeStamp:", ret)
	assert.LessOrEqual(t, ret, int64(0))
}

func TestDateTimeToHour(t *testing.T) {
	ret := DateTimeToHour(TestTimestamp)
	fmt.Println("DateTimeToHour:", ret)
	assert.Greater(t, ret, int64(0))
}

func TestDateToDayNumberInMonth(t *testing.T) {
	ret := DateToDayNumberInMonth(TestTimestamp)
	fmt.Println("DateToDayNumberInMonth:", ret)
	assert.Greater(t, ret, 0)
}

func TestDateToDayName(t *testing.T) {
	ret := DateToDayName(TestTimestamp)
	fmt.Println("DateToDayName:", ret)
	assert.NotEqual(t, "", ret)
}

func TestDateToMonthName(t *testing.T) {
	ret := DateToMonthName(TestTimestamp)
	fmt.Println("DateToMonthName:", ret)
	assert.NotEqual(t, "", ret)
}

func TestTimeStampFromYearAndDay(t *testing.T) {
	ret := TimeStampFromYearAndDay(2024, 1)
	fmt.Println("TimeStampFromYearAndDay:", ret)
	assert.Greater(t, ret, int64(0))
}

func TestIsValidTimeFormat(t *testing.T) {
	ret := IsValidTimeFormat("")
	assert.Equal(t, false, ret)

	ret = IsValidTimeFormat("12:36:3")
	assert.Equal(t, false, ret)

	ret = IsValidTimeFormat("12:36:31.0000")
	assert.Equal(t, false, ret)

	ret = IsValidTimeFormat("12:36:31.000.000")
	assert.Equal(t, false, ret)

	ret = IsValidTimeFormat("12:36:61.000")
	assert.Equal(t, false, ret)

	ret = IsValidTimeFormat("12:36:31.000")
	assert.Equal(t, true, ret)
}

func TestIsValidDateFormat(t *testing.T) {
	ret := IsValidDateFormat("")
	assert.Equal(t, false, ret)

	ret = IsValidDateFormat("2024-01-100")
	assert.Equal(t, false, ret)

	ret = IsValidDateFormat("2024-01-10-00")
	assert.Equal(t, false, ret)

	ret = IsValidDateFormat("2024-01-40")
	assert.Equal(t, false, ret)

	ret = IsValidDateFormat("2024-01-10")
	assert.Equal(t, true, ret)
}

func TestIsValidDateTimeFormat(t *testing.T) {
	ret := IsValidDateTimeFormat("")
	assert.Equal(t, false, ret)

	ret = IsValidDateTimeFormat("2024-01-10 12:36:31.0000")
	assert.Equal(t, false, ret)

	ret = IsValidDateTimeFormat("2024-01-10 12:36:31.000 000")
	assert.Equal(t, false, ret)

	ret = IsValidDateTimeFormat("2024-01-10 12:36:71.000")
	assert.Equal(t, false, ret)

	ret = IsValidDateTimeFormat("2024-01-10 12:36:31.000")
	assert.Equal(t, true, ret)
}
