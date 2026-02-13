package calendar

import (
	"strconv"
	"time"
)

// calendarDay represents a single day cell in the calendar grid.
type calendarDay struct {
	Date    time.Time
	InMonth bool
	IsToday bool
}

// DateString returns the date in "2006-01-02" format for signal values.
func (d calendarDay) DateString() string {
	return d.Date.Format("2006-01-02")
}

// DayLabel returns the day number as a string.
func (d calendarDay) DayLabel() string {
	return strconv.Itoa(d.Date.Day())
}

var weekdayLabels = []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}

// buildGrid returns a 6×7 (42-cell) grid of days for the given month.
// The week starts on Monday. Days from the previous/next month fill
// the edges, with InMonth set to false.
func buildGrid(year int, month time.Month) []calendarDay {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Offset from Monday: Monday=0 … Sunday=6.
	wd := int(first.Weekday())
	if wd == 0 {
		wd = 6
	} else {
		wd--
	}

	start := first.AddDate(0, 0, -wd)

	days := make([]calendarDay, 42)
	for i := range days {
		d := start.AddDate(0, 0, i)
		days[i] = calendarDay{
			Date:    d,
			InMonth: d.Month() == month,
			IsToday: d.Equal(today),
		}
	}
	return days
}
