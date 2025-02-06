package authentication

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	common_utils "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neoEntity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

type authenticationService struct {
	neo4j *neo4j_repository.Repositories
	user  interfaces.UserService
	email interfaces.EmailService
}

func NewAuthenticationService(neo4j *neo4j_repository.Repositories, user interfaces.UserService, email interfaces.EmailService) interfaces.AuthenticationService {
	return &authenticationService{
		neo4j: neo4j,
		user:  user,
		email: email,
	}
}

func (s *authenticationService) SetUserService(user interfaces.UserService) {
	s.user = user
}

func (s *authenticationService) SetEmailService(email interfaces.EmailService) {
	s.email = email
}

func (s *authenticationService) IsInitialized() bool {
	if s.neo4j == nil || s.user == nil || s.email == nil {
		return false
	}
	return true
}

func (a *authenticationService) CreateUserInTenant(ctx context.Context, txWithPostCommit *common_utils.TxWithPostCommit, tenant string, impersonating bool, authUserId, email, firstName, lastName string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationService.CreateUserInTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	authenticatedUserInTenant, err := a.neo4j.UserReadRepository.FindFirstUserWithRolesByEmail(ctx, tenant, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if authenticatedUserInTenant != nil && authenticatedUserInTenant.UserId != "" {
		span.LogKV("user_already_exists", true)
		return authenticatedUserInTenant.UserId, nil
	}

	span.LogKV("user_already_exists", false)

	roles := []string{"USER"}

	if impersonating {
		roles = append(roles, "IMPERSONATED")
	} else {
		roles = append(roles, "OWNER")
	}

	//set the tenant received in context
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: tenant,
	})

	userId, err := common_utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, a.neo4j.Neo4jDriver, a.neo4j.Database, txWithPostCommit, func(txWithPostCommit *common_utils.TxWithPostCommit) (any, error) {

		userId, err := a.user.Save(ctx, txWithPostCommit, nil, data_fields.UserFields{
			FirstName: &firstName,
			LastName:  &lastName,
			Roles:     common_utils.ToPtr(roles),
		})
		if err != nil {
			return "", err
		}

		err = a.neo4j.AuthenticationWriteRepository.LinkAuthenticationUserWithUser(ctx, *txWithPostCommit.Tx, authUserId, userId)
		if err != nil {
			return "", err
		}

		emailNode, err := a.neo4j.EmailReadRepository.GetFirstByEmail(ctx, tenant, email)
		if err != nil {
			return "", err
		}

		if emailNode == nil {
			_, err = a.email.Merge(ctx, txWithPostCommit, tenant, interfaces.EmailFields{
				Primary: true,
				Email:   email,
				Source:  neoEntity.DataSourceOpenline,
			}, &common_srv.LinkWith{
				Type: model.USER,
				Id:   userId,
			})
			if err != nil {
				return "", err
			}
		} else {
			emailEntity := mapper.MapDbNodeToEmailEntity(emailNode)

			err = a.neo4j.EmailWriteRepository.LinkWithUser(ctx, txWithPostCommit.Tx, tenant, userId, emailEntity.Id, true)
			if err != nil {
				return "", err
			}
		}

		return userId, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return userId.(string), nil
}
