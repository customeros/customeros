package helper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/graph/model"
	"time"
)

func TimeFilterFromValue(timeFilter model.TimeFilter) time.Time {
	switch timeFilter.TimePeriod {
	case model.TimePeriodCustom:
		return defaultIfNil(timeFilter.From)
	case model.TimePeriodLastHour:
		dur, _ := time.ParseDuration("-1h")
		return time.Now().UTC().Add(dur)
	case model.TimePeriodLast24Hours:
		dur, _ := time.ParseDuration("-24h")
		return time.Now().UTC().Add(dur)
	case model.TimePeriodLast7Days:
		return time.Now().UTC().AddDate(0, 0, -7)
	case model.TimePeriodLast30Days:
		return time.Now().UTC().AddDate(0, 0, -30)
	case model.TimePeriodToday:
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	case model.TimePeriodDaily:
		selectedDate := defaultIfNil(timeFilter.From)
		return time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), 0, 0, 0, 0, time.UTC)
	case model.TimePeriodMonthly:
		selectedDate := defaultIfNil(timeFilter.From)
		return time.Date(selectedDate.Year(), selectedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	case model.TimePeriodMonthToDate:
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	case model.TimePeriodYearToDate:
		now := time.Now().UTC()
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	case model.TimePeriodAllTime:
		return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	default:
		return time.Now().UTC()
	}
}

func TimeFilterToValue(timeFilter model.TimeFilter) time.Time {
	switch timeFilter.TimePeriod {
	case model.TimePeriodCustom:
		return defaultIfNil(timeFilter.To)
	case model.TimePeriodLastHour:
	case model.TimePeriodLast24Hours:
	case model.TimePeriodLast7Days:
	case model.TimePeriodLast30Days:
	case model.TimePeriodMonthToDate:
	case model.TimePeriodYearToDate:
	case model.TimePeriodAllTime:
		return time.Now().UTC()
	case model.TimePeriodToday:
		tomorrow := time.Now().UTC().AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
	case model.TimePeriodDaily:
		nextDate := defaultIfNil(timeFilter.From).AddDate(0, 0, 1)
		return time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, time.UTC)
	case model.TimePeriodMonthly:
		nextMonth := defaultIfNil(timeFilter.From).AddDate(0, 1, 0)
		return time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	default:
		return time.Now().UTC()
	}
	return time.Now().UTC()
}

func defaultIfNil(checkedTime *time.Time) time.Time {
	if checkedTime != nil {
		return time.Now().UTC()
	}
	return (*checkedTime).UTC()
}
