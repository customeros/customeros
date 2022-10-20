package helper

import (
	"fmt"
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
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
