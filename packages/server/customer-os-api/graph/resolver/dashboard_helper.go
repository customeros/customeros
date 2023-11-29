package resolver

import (
	"fmt"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func getPeriod(period *model.DashboardPeriodInput) (time.Time, time.Time, error) {
	if period == nil {
		now := time.Now().UTC()

		//last 12 months including current month
		startDate := now.AddDate(-1, 1, 0)
		endDate := now

		return startDate, endDate, nil
	} else {
		if period.Start.After(period.End) {
			return time.Time{}, time.Time{}, fmt.Errorf("start date must be before end date")
		}
		return period.Start, period.End, nil
	}
}
