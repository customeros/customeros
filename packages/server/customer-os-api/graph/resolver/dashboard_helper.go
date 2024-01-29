package resolver

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func getPeriod(period *model.DashboardPeriodInput, now time.Time) (time.Time, time.Time, error) {
	if period == nil {
		//last 12 months including current month
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
