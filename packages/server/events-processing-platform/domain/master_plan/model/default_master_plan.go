package model

import (
	"time"

	"github.com/google/uuid"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

func NewDefaultMasterPlan(id string, currentTime *time.Time) MasterPlan {
	m1ID := uuid.New().String()
	m2ID := uuid.New().String()
	m3ID := uuid.New().String()
	return MasterPlan{
		ID:           id,
		Name:         "Default Master Plan",
		Retired:      false,
		CreatedAt:    *currentTime,
		UpdatedAt:    *currentTime,
		SourceFields: commonmodel.Source{},
		Milestones: map[string]MasterPlanMilestone{
			m1ID: NewDefaultMasterPlanMilestone(m1ID, "Milestone 1", []string{"item1", "item2"}, *currentTime, 0),
			m2ID: NewDefaultMasterPlanMilestone(m2ID, "Milestone 2", []string{"item3", "item4"}, *currentTime, 1),
			m3ID: NewDefaultMasterPlanMilestone(m3ID, "Milestone 3", []string{"item5", "item6"}, *currentTime, 2),
		},
	}
}

func NewDefaultMasterPlanMilestone(id, name string, items []string, currentTime time.Time, order int64) MasterPlanMilestone {
	return MasterPlanMilestone{
		ID:            id,
		Name:          name,
		Retired:       false,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
		SourceFields:  commonmodel.Source{},
		Optional:      false,
		Order:         order,
		DurationHours: 24,
		Items:         items,
	}
}
