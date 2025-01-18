package resolver

import (
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func getPeriod(period *model.DashboardPeriodInput, now time.Time) (time.Time, time.Time, error) {
	if period == nil {
		// last 12 months including current month
		startDate := utils.LastTimeOfMonth(now.Year()-1, int(now.Month())+1)
		endDate := now

		return startDate, endDate, nil
	} else {
		if period.Start.After(period.End) {
			return time.Time{}, time.Time{}, fmt.Errorf("start date must be before end date")
		}
		return period.Start, period.End, nil
	}
}
