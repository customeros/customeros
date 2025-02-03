package agent

import (
	"context"
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/aws_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

const (
	BUCKET     = "agent-registry"
	AWS_REGION = "eu-west-1"
)

type Agent struct {
	Version      string       `toml:"version"`
	Name         string       `toml:"name"`
	Type         string       `toml:"type"`
	Goal         string       `toml:"goal"`
	Icon         string       `toml:"icon"`
	Triggers     Triggers     `toml:"triggers"`
	Capabilities Capabilities `toml:"capabilities"`
}

type Triggers struct {
	Events []string `toml:"events"`
}

type Capabilities struct {
	Types []string `toml:"types"`
}

type agentRegistryService struct {
	postgresRepositories     *postgres_repository.Repositories
	s3client                 aws_client.S3Client
	agentCapabilityExecutors map[enum.AgentCapability]interfaces.AgentCapabilityUntyped
}

func NewAgentRegistryService(
	postgresRepositories *postgres_repository.Repositories,
	agentCapabilityExecutors map[enum.AgentCapability]interfaces.AgentCapabilityUntyped,
) interfaces.AgentRegistry {
	if postgresRepositories == nil {
		panic("postgresRepositories cannot be nil")
	}

	s3 := aws_client.NewS3Client(&aws.Config{Region: aws.String(AWS_REGION)})

	return &agentRegistryService{
		postgresRepositories:     postgresRepositories,
		s3client:                 s3,
		agentCapabilityExecutors: agentCapabilityExecutors,
	}
}

func (r *agentRegistryService) SyncRegistry(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentRegistryService.SyncRegistry")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// get list of all agent config files from S3
	agentConfigFiles, err := r.getAgentConfigFiles(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var errs error
	for _, file := range agentConfigFiles {
		err = r.processAgentConfigFile(ctx, file)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (r *agentRegistryService) processAgentConfigFile(ctx context.Context, filename string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentRegistryService.processAgentConfigFile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentConfig, err := r.getAgentConfig(ctx, filename)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	agentType, err := enum.GetAgentType(agentConfig.Type)
	if err != nil {
		span.LogKV("agentType", agentConfig.Type)
		err := errors.New("Not a valid agent type")
		tracing.TraceErr(span, err)
		return err
	}

	dbAgent := postgres_entity.AgentRegistry{
		Type:         agentType,
		Version:      agentConfig.Version,
		Filename:     filename,
		Triggers:     agentConfig.Triggers.Events,
		Capabilities: agentConfig.Capabilities.Types,
		Icon:         agentConfig.Icon,
		IsActive:     true,
	}

	// Check if agent already exists in DB
	existingAgent, err := r.postgresRepositories.AgentRegistryRepository.FindByType(ctx, dbAgent.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if existingAgent == nil {
		// Create new agent
		_, err = r.postgresRepositories.AgentRegistryRepository.Create(ctx, dbAgent)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}
	return nil
}

func (r *agentRegistryService) getAgentConfigFiles(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentRegistryService.getAgentConfigFile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return r.s3client.ListFiles(ctx, BUCKET)
}

func (r *agentRegistryService) getAgentConfig(ctx context.Context, filename string) (*Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentRegistryService.getAgentConfig")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	file, err := r.s3client.Download(ctx, BUCKET, filename)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return r.loadAgentConfig(ctx, file)
}

func (r *agentRegistryService) loadAgentConfig(ctx context.Context, file string) (*Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentRegistryService.loadAgentConfig")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var agent Agent
	_, err := toml.Decode(file, &agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &agent, nil
}
