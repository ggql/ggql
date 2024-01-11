package ast

import (
	"testing"
)

const (
	TEST_TIME           = "08:30:00"
	TEST_DATE           = "2024-01-02"
	TEST_DATE_TIME      = "2024-01-02 08:30:00"
	TEST_DATE_TIME_FULL = "2024-01-02 08:30:00.000"
)

func Test_TimeStampFromYearAndDay(t *testing.T) {
	type args struct {
		year      int
		dayOfYear uint
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: "", args: args{year: 2024, dayOfYear: 2}, want: 1704153600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeStampFromYearAndDay(tt.args.year, tt.args.dayOfYear); got != tt.want {
				t.Errorf("timeStampFromYearAndDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStampToDate(t *testing.T) {
	type args struct {
		timeStamp int64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{timeStamp: 1704153600}, want: TEST_DATE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeStampToDate(tt.args.timeStamp); got != tt.want {
				t.Errorf("TimeStampToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStampToDateTime(t *testing.T) {
	type args struct {
		timeStamp int64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{timeStamp: 1704155400}, want: TEST_DATE_TIME},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeStampToDateTime(tt.args.timeStamp); got != tt.want {
				t.Errorf("TimeStampToDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStampToFullDateTime(t *testing.T) {
	type args struct {
		timeStamp int64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{timeStamp: 1704155400}, want: TEST_DATE_TIME_FULL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeStampToFullDateTime(tt.args.timeStamp); got != tt.want {
				t.Errorf("TimeStampToFullDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateToTimeStamp(t *testing.T) {
	type args struct {
		date string
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{date: TEST_DATE}, want: 1704153600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateToTimeStamp(tt.args.date); got != tt.want {
				t.Errorf("DateToTimeStamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTimeToTimeStamp(t *testing.T) {
	type args struct {
		date string
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{date: TEST_DATE_TIME}, want: 1704184200},
		{name: "test2", args: args{date: TEST_DATE_TIME_FULL}, want: 1704184200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateTimeToTimeStamp(tt.args.date); got != tt.want {
				t.Errorf("DateTimeToTimeStamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStampToTime(t *testing.T) {
	type args struct {
		timeStamp int64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{timeStamp: 1704155400}, want: "08:30:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeStampToTime(tt.args.timeStamp); got != tt.want {
				t.Errorf("TimeStampToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeStampFromYearAndDay(t *testing.T) {
	type args struct {
		year      int
		dayOfYear uint
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{year: 2024, dayOfYear: 2}, want: 1704153600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeStampFromYearAndDay(tt.args.year, tt.args.dayOfYear); got != tt.want {
				t.Errorf("TimeStampFromYearAndDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidTimeFormat(t *testing.T) {
	type args struct {
		timeStr string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "test-ok", args: args{timeStr: TEST_TIME}, want: true},
		{name: "test-err1", args: args{timeStr: "8:30"}, want: false},
		{name: "test-err2", args: args{timeStr: "8:70"}, want: false},
		{name: "test-err3", args: args{timeStr: "33:20:00"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTimeFormat(tt.args.timeStr); got != tt.want {
				t.Errorf("IsValidTimeFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDateFormat(t *testing.T) {
	type args struct {
		dateStr string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "test-ok", args: args{dateStr: TEST_DATE}, want: true},
		{name: "test-err1", args: args{dateStr: "2024-33-04"}, want: false},
		{name: "test-err2", args: args{dateStr: "2024-01-60"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDateFormat(tt.args.dateStr); got != tt.want {
				t.Errorf("IsValidDateFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDateTimeFormat(t *testing.T) {
	type args struct {
		datetimeStr string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "test-ok1", args: args{datetimeStr: TEST_DATE_TIME}, want: true},
		{name: "test-ok2", args: args{datetimeStr: TEST_DATE_TIME_FULL}, want: true},
		{name: "test-err1", args: args{datetimeStr: "2024-33-02 08:30:00"}, want: false},
		{name: "test-err2", args: args{datetimeStr: "2024-01-53 08:30:00"}, want: false},
		{name: "test-err3", args: args{datetimeStr: "2024-01-02 35:30:00"}, want: false},
		{name: "test-err4", args: args{datetimeStr: "2024-01-02 08:87:00"}, want: false},
		{name: "test-err5", args: args{datetimeStr: "2024-01-02 08:30:99"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDateTimeFormat(tt.args.datetimeStr); got != tt.want {
				t.Errorf("IsValidDateTimeFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
