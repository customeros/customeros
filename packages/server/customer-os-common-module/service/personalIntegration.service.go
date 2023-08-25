package service

import (
	"fmt"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type Integration string

const (
	CALCOM Integration = "calcom"
)
const CalComHeader = "x-cal-signature-256"

func SignatureCheck(hSignature string, body []byte, personalIntegrationRepository repository.PersonalIntegrationRepository, tenant, email, integration string) error {
	if hSignature != "" {
		result := personalIntegrationRepository.FindIntegration(tenant, email, integration)

		if result.Error != nil {
			return fmt.Errorf("SignatureCheck error: %v", result.Error.Error())
		}

		secret := result.Result.(*entity.PersonalIntegration)

		if secret == nil {
			return fmt.Errorf("SignatureCheck error: no information found")
		}
		cSignature := utils.Hmac(body, []byte(secret.Secret))
		if hSignature != *cSignature {
			return fmt.Errorf("SignatureCheck error: signature mismatch")
		}
		return nil
	} else {
		return fmt.Errorf("SignatureCheck error: signature header not found")
	}
}
