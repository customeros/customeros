package meeting

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type meetingService struct {
	log          logger.Logger
	nylasService interfaces.NylasService
}

func NewMeetingService(log logger.Logger, nylasService interfaces.NylasService) interfaces.MeetingService {
	return &meetingService{
		log:          log,
		nylasService: nylasService,
	}
}

func (s *meetingService) IsInitialized() bool {
	return utils.IsInitialized(s)
}
