package notifications

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Boostport/mjml-go"
	novu "github.com/novuhq/go-novu/lib"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type NovuProvider struct {
	NovuClient      *novu.APIClient
	TemplatePath    string
	log             logger.Logger
	emailRawContent string
	emailContent    string
}

func newNovuProvider(log logger.Logger, apiKey, templatePath string) *NovuProvider {
	return &NovuProvider{
		NovuClient:   novu.NewAPIClient(apiKey, &novu.Config{}),
		log:          log,
		TemplatePath: templatePath,
	}
}

func (np *NovuProvider) SendNotification(ctx context.Context, u *NotifiableUser, payload map[string]interface{}, workflowId string) error {
	to := map[string]interface{}{
		"lastName":     u.LastName,
		"firstName":    u.FirstName,
		"subscriberId": u.SubscriberID,
		"email":        u.Email,
	}
	data := novu.ITriggerPayloadOptions{To: to, Payload: payload}
	if containsKey(payload, "overrides") {
		overrides := payload["overrides"].(map[string]interface{})
		if len(overrides) > 0 {
			if containsKey(overrides, "email") {
				data.Overrides = overrides
			}
		}
	}

	_, err := np.NovuClient.EventApi.Trigger(ctx, workflowId, data)

	if err != nil {
		np.log.Errorf("(NotificationsSubscriber.NovuProvider.SendNotification) error: %s", err.Error())
		return err
	}

	return nil
}

func (np *NovuProvider) LoadEmailBody(ctx context.Context, workflowId string) error {
	var emailPath string
	switch workflowId {
	case WorkflowIdOrgOwnerUpdateEmail:
		if _, err := os.Stat(np.TemplatePath); os.IsNotExist(err) {
			return fmt.Errorf("(NovuProvider.LoadEmailBody) error: %s", err.Error())
		}
		emailPath = fmt.Sprintf("%s/ownership.single.mjml", np.TemplatePath)
	}
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("(NovuProvider.LoadEmailBody) Getwd: %s", err.Error())
	}

	// Build the full path to the template file
	templatePath := filepath.Join(currentDir, emailPath)
	if _, err := os.Stat(templatePath); err != nil {
		return fmt.Errorf("(NovuProvider.LoadEmailBody) error: %s", err.Error())
	}

	rawMjml, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("(NovuProvider.LoadEmailBody) error: %s", err.Error())
	}
	np.emailRawContent = string(rawMjml[:])
	return nil
}

func (np *NovuProvider) Template(ctx context.Context, replace map[string]string) (string, error) {
	_, ok := replace["{{userFirstName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing userFirstName")
	}
	_, ok = replace["{{actorFirstName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing actorFirstName")
	}
	_, ok = replace["{{actorLastName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing actorLastName")
	}
	_, ok = replace["{{orgName}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing orgName")
	}
	_, ok = replace["{{orgLink}}"]
	if !ok {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", "missing orgLink")
	}
	mjmlf := np.emailRawContent
	for k, v := range replace {
		mjmlf = strings.Replace(mjmlf, k, v, -1)
	}
	np.emailRawContent = mjmlf
	// mjmlf := strings.Replace(string(np.emailRawContent[:]), "{{userFirstName}}", userFirstName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{actorFirstName}}", actorFirstName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{actorLastName}}", actorLastName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{orgName}}", orgName, -1)
	// mjmlf = strings.Replace(mjmlf, "{{orgLink}}", orgLink, -1)

	html, err := mjml.ToHTML(context.Background(), mjmlf)
	var mjmlError mjml.Error
	if errors.As(err, &mjmlError) {
		return "", fmt.Errorf("(NovuProvider.Template) error: %s", mjmlError.Message)
	}
	np.emailContent = html
	return html, err
}

func (np *NovuProvider) GetRawContent() string {
	return np.emailRawContent
}

func containsKey[M ~map[K]V, K comparable, V any](m M, k K) bool {
	_, ok := m[k]
	return ok
}
