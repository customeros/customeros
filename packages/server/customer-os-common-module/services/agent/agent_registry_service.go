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
	Agent        AgentMetadata   `toml:"agent"`
	Goal         GoalConfig      `toml:"goal"`
	Watchers     WatchersConfig  `toml:"watchers"`
	Capabilities CapabilityTypes `toml:"capabilities"`
	Plays        Plays           `toml:"plays"`
}

type AgentMetadata struct {
	Version string `toml:"version"`
	Name    string `toml:"name"`
	Type    string `toml:"type"`
	Icon    string `toml:"icon"`
}

type GoalConfig struct {
	Goal             string   `toml:"goal"`
	CompletionEvents []string `toml:"completion_events"`
}

type WatchersConfig struct {
	Events []string `toml:"events"`
}

type CapabilityTypes struct {
	Types []string `toml:"types"`
}

type Plays map[string][]string

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryService.processAgentConfigFile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentConfig, err := r.getAgentConfig(ctx, filename)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	agentType, err := enum.GetAgentType(agentConfig.Agent.Type)
	if err != nil {
		span.LogKV("agentType", agentConfig.Agent.Type)
		err := errors.New("Not a valid agent type")
		tracing.TraceErr(span, err)
		return err
	}

	// Create agent registry entity
	dbAgent := postgres_entity.AgentRegistry{
		Type:             agentType,
		AgentName:        agentConfig.Agent.Name,
		Goal:             agentConfig.Goal.Goal,
		CompletionEvents: agentConfig.Goal.CompletionEvents,
		ListenerEvents:   agentConfig.Watchers.Events,
		Capabilities:     agentConfig.Capabilities.Types,
		Version:          agentConfig.Agent.Version,
		Filename:         filename,
		Icon:             agentConfig.Agent.Icon,
		IsActive:         true,
	}

	// Create agent plays
	var plays []postgres_entity.AgentPlay
	for triggerEvent, capabilities := range agentConfig.Plays {
		play := postgres_entity.AgentPlay{
			AgentType:    agentType,
			TriggerEvent: triggerEvent,
			Capabilities: capabilities,
		}
		plays = append(plays, play)
	}

	// Check if agent already exists in DB
	existingAgent, err := r.postgresRepositories.AgentRegistryRepository.FindByType(ctx, dbAgent.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if existingAgent == nil {
		// Create new agent with plays
		_, err = r.postgresRepositories.AgentRegistryRepository.Create(ctx, dbAgent, plays)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		// Update existing agent with new plays
		dbAgent.ID = existingAgent.ID
		_, err = r.postgresRepositories.AgentRegistryRepository.Update(ctx, dbAgent, plays)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (r *agentRegistryService) getAgentConfigFiles(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryService.getAgentConfigFile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return r.s3client.ListFiles(ctx, BUCKET)
}

func (r *agentRegistryService) getAgentConfig(ctx context.Context, filename string) (*Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryService.getAgentConfig")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRegistryService.loadAgentConfig")
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
