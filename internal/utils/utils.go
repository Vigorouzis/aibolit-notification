package utils

import "time"

func RoundToNearestQuarter(t time.Time) time.Time {
	minutes := t.Minute()
	remainder := minutes % 15

	if remainder == 0 {
		return t
	}

	return t.Add(time.Duration(15-remainder) * time.Minute)
}

func CalculateIntakeTimes(frequency int) []string {
	if frequency <= 0 {
		return nil
	}

	const totalHours = 14
	const minutesPerHour = 60
	totalMinutes := totalHours * minutesPerHour
	interval := totalMinutes / frequency

	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Local)
	endTime := startTime.Add(14 * time.Hour)

	var intakeTimes []string

	for i := 0; i < frequency; i++ {
		timeSlot := startTime.Add(time.Duration(i*interval) * time.Minute)

		timeSlot = RoundToNearestQuarter(timeSlot)

		if timeSlot.After(endTime) {
			timeSlot = endTime
		}

		intakeTimes = append(intakeTimes, timeSlot.Format("15:04"))
	}

	return intakeTimes
}
