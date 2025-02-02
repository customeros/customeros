package agent

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/aws_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
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
	Capabilities []Capability `toml:"capabilities"`
}

type Triggers struct {
	Events []string `toml:"events"`
}

type Capability struct {
	Type          string      `toml:"type"`
	ConfigDefault interface{} `toml:"config_default,omitempty"`
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

	for _, file := range agentConfigFiles {
		agentConfig, err := r.loadAgentConfig(ctx, file)
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
			Type:     agentType,
			Name:     agentConfig.Name,
			Goal:     agentConfig.Goal,
			Icon:     agentConfig.Icon,
			IsActive: true,
		}

		// process capabilities
		var dbCapabilities []postgres_entity.Capability
		for _, cap := range agentConfig.Capabilities {
			capType, err := enum.GetAgentCapability(cap.Type)
			if err != nil {
				span.LogKV("capabilityType", cap.Type)
				tracing.TraceErr(span, err)
				return err
			}

			// process capability config
			configStr := ""
			if cap.ConfigDefault != nil {
				config, err := r.processAgentConfig(ctx, cap)
				if err != nil {
					tracing.LogObjectAsJson(span, "config", config)
					tracing.TraceErr(span, err)
					return err
				}

				// convert config to JSON string
				configBytes, err := json.Marshal(config)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}
				configStr = string(configBytes)
			}

			dbCapabilities = append(dbCapabilities, postgres_entity.Capability{
				Type:   capType,
				Config: configStr,
				Active: true,
			})
		}

		dbAgent.CapabilitiesConfig = postgres_entity.CapabilitiesConfig{
			Capabilities: dbCapabilities,
		}

		// Check if agent already exists in DB
		existingAgent, err := r.postgresRepositories.AgentRegistryRepository.FindByType(ctx, dbAgent.Type)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if existingAgent != nil {
			// Update existing agent
			existingAgent.Name = dbAgent.Name
			existingAgent.Goal = dbAgent.Goal
			existingAgent.Icon = dbAgent.Icon
			existingAgent.CapabilitiesConfig = dbAgent.CapabilitiesConfig

			_, err = r.postgresRepositories.AgentRegistryRepository.Update(ctx, *existingAgent)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		} else {
			// Create new agent
			_, err = r.postgresRepositories.AgentRegistryRepository.Create(ctx, dbAgent)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
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

	for _, cap := range agent.Capabilities {
		_, err := r.processAgentConfig(ctx, cap)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return &agent, nil
}

func (r *agentRegistryService) processAgentConfig(ctx context.Context, capability Capability) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentRegistryService.processAgentConfig")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "capability", capability)

	capType, err := enum.GetAgentCapability(capability.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	_, exists := r.agentCapabilityExecutors[capType]
	if !exists {
		err = errors.New("capability not registered")
		tracing.TraceErr(span, err)
		return nil, err
	}

	typedCap, ok := agent_capability.GetTypedExecutor[any, any, any](
		r.agentCapabilityExecutors,
		capType,
	)
	if !ok {
		err = errors.New("failed to get typed capability")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Get base config for capability
	config := typedCap.GetConfig()

	// If no config in TOML, return base config
	if capability.ConfigDefault == nil {
		return config, nil
	}

	// Marshal TOML config to bytes for processing
	configBytes, err := toml.Marshal(capability.ConfigDefault)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Unmarshal into the correct config type
	err = toml.Unmarshal(configBytes, &config)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return config, nil
}
