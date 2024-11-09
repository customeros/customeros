package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"golang.org/x/net/context"
)

type Integration string

const (
	CALCOM Integration = "calcom"
)
const CalComHeader = "x-cal-signature-256"

func SignatureCheck(hSignature string, body []byte, personalIntegrationRepository postgresRepository.PersonalIntegrationRepository, tenant, email, integration string) error {
	if hSignature != "" {
		result := personalIntegrationRepository.FindIntegration(context.TODO(), tenant, email, integration)

		if result.Error != nil {
			return fmt.Errorf("SignatureCheck error: %v", result.Error.Error())
		}

		secret := result.Result.(*postgresEntity.PersonalIntegration)

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
