package helper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/graph/model"
	"gorm.io/gorm"
)

func AddDataFilter(column string, action model.Operation, value string, db *gorm.DB) *gorm.DB {
	if column != "" {
		switch action {
		case model.OperationEquals:
			whereQuery := fmt.Sprintf("%s = ?", column)
			db = db.Where(whereQuery, value)
			break
		case model.OperationContains:
			whereQuery := fmt.Sprintf("%s LIKE ?", column)
			db = db.Where(whereQuery, "%"+value+"%")
			break
		}
	}
	return db
}

func AddTimeFilter(timeFilter model.TimeFilter, column string, db *gorm.DB) *gorm.DB {

	whereQueryFrom := fmt.Sprintf("%s >= ?", column)
	db = db.Where(whereQueryFrom, TimeFilterFromValue(timeFilter))

	whereQueryTo := fmt.Sprintf("%s <= ?", column)
	db = db.Where(whereQueryTo, TimeFilterToValue(timeFilter))

	return db
}
