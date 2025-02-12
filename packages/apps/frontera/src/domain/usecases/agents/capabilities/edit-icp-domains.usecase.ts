import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { Agent } from '@store/Agents/Agent.dto';
import { when, action, reaction, computed, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service';
import { Organization } from '@store/Organizations/Organization.dto';
import { OrganizationRepository } from '@infra/repositories/organization';
import { SearchGlobalOrganizationsQuery } from '@infra/repositories/organization/queries/searchGlobalOrganizations.generated';

import { validateUrl } from '@utils/url';
import { CapabilityType } from '@graphql/types';
import { SortingDirection, ComparisonOperator } from '@graphql/types';

type GlobalOrganization =
  SearchGlobalOrganizationsQuery['globalOrganizations_Search'][number];

type Status = 'LOADING' | 'ERROR' | 'IDLE' | 'VALIDATING_DOMAIN';

export class EditIcpDomainsUsecase {
  private root = RootStore.getInstance();
  private service = OrganizationRepository.getInstance();
  private agentService = new AgentService();
  private agentId: string = '';
  @observable accessor icpCompanyExamplesError: string = '';

  @observable private accessor searchedIds: string[] = [];
  @observable public accessor icpCompanyExamples: Set<string> = new Set();

  @observable accessor searchTerm = '';
  @observable accessor errorMessage = '';
  @observable accessor status: Status = 'IDLE';
  @observable accessor domainValidationError = '';
  @observable accessor domainValidationMessage = '';
  @observable accessor options: AddSearchOrganizationOption[] = [];

  constructor() {
    this.execute = this.execute.bind(this);
    this.select = this.select.bind(this);
    this.selectCustom = this.selectCustom.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.searchGlobal = this.searchGlobal.bind(this);
    this.reset = this.reset.bind(this);

    when(() => !!this.root.session.sessionToken, this.searchGlobal);
    reaction(() => this.searchTerm, this.execute);
  }

  @action
  init(agentId: string) {
    const span = Tracer.span('EditIcpDomainsUsecase.init', {
      icpCompanyExamples: this.icpCompanyExamples,
    });

    this.agentId = agentId;

    const agent = this.root.agents.getById(agentId);

    if (!agent) {
      console.error(
        'EditIcpDomainsUsecase.init: Agent not found. aborting execution',
      );

      return;
    }

    const capability = agent.value.capabilities.find(
      (c) => c.type === CapabilityType.IcpQualify,
    );

    if (!capability) {
      console.error(
        'EditIcpDomainsUsecase.init: Capability not found. aborting execution',
      );

      return;
    }

    const config = Agent.parseConfig(capability.config);

    if (!config) {
      console.error(
        'EditIcpDomainsUsecase.init: Could not parse config. aborting',
      );

      return;
    }

    if (!config.icpCompanyExamples) {
      console.error(
        'EditIcpDomainsUsecase.init: icpCompanyExamples not found in config. aborting execution',
      );

      return;
    }

    this.icpCompanyExamples = Array.isArray(config?.icpCompanyExamples?.value)
      ? new Set(config.icpCompanyExamples.value)
      : new Set();
    this.icpCompanyExamplesError = config.icpCompanyExamples.error as string;
    span.end({
      icpCompanyExamples: config.icpCompanyExamples.value,
      icpCompanyExamplesError: this.icpCompanyExamplesError,
    });
  }

  @computed
  get isLoading() {
    return this.status === 'LOADING';
  }

  @computed
  get isValidatingDomain() {
    return this.status === 'VALIDATING_DOMAIN';
  }

  @computed
  get mixedOptions() {
    const tenantOrgs = this.searchedIds.reduce((acc, id) => {
      const entity = this.root.organizations.getById(id);

      if (entity) {
        acc.push(new AddSearchOrganizationOption(entity));
      }

      return acc;
    }, [] as AddSearchOrganizationOption[]);

    if (!this.searchTerm) {
      return this.options;
    }

    return [...tenantOrgs, ...this.options];
  }

  @action
  public async select(website: string) {
    const span = Tracer.span('EditIcpDomainsUsecase.select', {
      website,
      currentlySelected: Array.from(this.icpCompanyExamples),
    });

    if (this.icpCompanyExamples.has(website)) {
      this.icpCompanyExamples.delete(website);
    } else {
      this.icpCompanyExamples.add(website);
    }

    const agent = this.root.agents.getById(this.agentId);

    if (agent) {
      agent.setCapabilityConfig(
        CapabilityType.IcpQualify,
        'icpCompanyExamples',
        Array.from(this.icpCompanyExamples),
      );

      const [res, err] = await this.agentService.saveAgent(agent);

      if (res) {
        agent.put(res?.agent_Save);
        this.reset();
        this.init(this.agentId);
      }

      if (err) {
        console.error(
          'EditIcpDomainsUsecase.select: Error saving agent. aborting execution',
        );
      }
    }

    span.end({
      icpCompanyExamples: Array.from(this.icpCompanyExamples),
    });
  }

  @action
  public async selectCustom(website: string) {
    if (website) {
      const isUrl = validateUrl(this.searchTerm);

      if (!isUrl) {
        this.domainError('Invalid URL');

        return;
      }

      try {
        await this.validateDomain();

        if (!this.domainValidationError) {
          this.select(website);

          return;
        }
      } catch (err) {
        this.domainError('Could not validate domain');
      }
    }
    this.reset();
  }

  @action
  private setOptions(options: AddSearchOrganizationOption[]) {
    this.options = options.filter((o) => !o.tenantOrganizationId);
  }

  @action
  private setSearchedIds(ids: string[]) {
    this.searchedIds = ids;
  }

  @action
  public reset() {
    this.searchTerm = '';
    this.errorMessage = '';
    this.status = 'IDLE';
    this.searchedIds = [];
    this.domainValidationError = '';
    this.domainValidationMessage = '';
  }

  @action
  public setSearchTerm(searchTerm: string) {
    this.searchTerm = searchTerm;
  }

  @action
  private loading() {
    this.status = 'LOADING';
  }

  @action
  private idle() {
    this.status = 'IDLE';
  }

  @action
  private error(error: string) {
    this.status = 'ERROR';
    this.errorMessage = error;
  }

  @action
  private validatingDomain() {
    this.status = 'VALIDATING_DOMAIN';
  }

  @action
  private domainError(error: string) {
    this.status = 'IDLE';
    this.domainValidationError = error;
  }

  @action
  private domainMessage(message: string) {
    this.status = 'IDLE';
    this.domainValidationMessage = message;
  }

  private async searchGlobal() {
    try {
      this.loading();

      const { globalOrganizations_Search } =
        await this.service.searchGlobalOrganizations({
          searchTerm: this.searchTerm,
          limit: 30,
        });

      this.setOptions(
        globalOrganizations_Search.map(
          (data) => new AddSearchOrganizationOption(data),
        ) ?? [],
      );
      this.idle();
    } catch (_err) {
      this.error('Could not search companies');
    }
  }

  private async searchTenant() {
    try {
      this.loading();

      const { ui_organizations_search } =
        await this.service.searchOrganizations({
          limit: 30,
          sort: {
            by: 'ORGANIZATIONS_NAME',
            caseSensitive: false,
            direction: SortingDirection.Asc,
          },
          where: {
            OR: [
              {
                filter: {
                  property: 'ORGANIZATIONS_NAME',
                  value: this.searchTerm,
                  caseSensitive: false,
                  includeEmpty: false,
                  operation: ComparisonOperator.Contains,
                },
              },
              {
                filter: {
                  property: 'ORGANIZATIONS_WEBSITE',
                  value: this.searchTerm,
                  caseSensitive: false,
                  includeEmpty: false,
                  operation: ComparisonOperator.Contains,
                },
              },
            ],
          },
        });

      const results = ui_organizations_search?.ids ?? [];

      if (results.length === 0) {
        this.idle();
        this.setSearchedIds([]);

        return;
      }

      await this.root.organizations.retrieve(results);
      this.setSearchedIds(results);
      this.idle();
    } catch (_err) {
      this.error('Could not search companies');
    }
  }

  private async validateDomain() {
    try {
      this.validatingDomain();

      const { organization_CheckWebsite } = await this.service.checkWebsite({
        website: this.searchTerm,
      });

      const isAccepted = organization_CheckWebsite?.accepted ?? false;
      const isPrimary = organization_CheckWebsite?.primary ?? false;
      const foundPrimaryDomain = organization_CheckWebsite?.primaryDomain;
      const foundGlobalOrganization =
        organization_CheckWebsite?.globalOrganization;

      if (!isAccepted) {
        this.domainError(`We don't allow hosting or profile-based domains`);

        return;
      }

      if (!isPrimary && foundPrimaryDomain) {
        if (foundGlobalOrganization) {
          if (foundGlobalOrganization?.organizationId) {
            this.root.organizations.retrieve([
              foundGlobalOrganization.organizationId,
            ]);
            this.searchedIds.push(foundGlobalOrganization.organizationId);
          } else {
            this.setOptions([
              new AddSearchOrganizationOption(foundGlobalOrganization),
            ]);
          }
        }
        this.domainMessage(`This domain is linked to ${foundPrimaryDomain}`);

        return;
      }

      this.domainError('');
      this.domainMessage('');
    } catch (_err) {
      this.error('Could not check website');
    }
  }

  public async execute() {
    const span = Tracer.span('EditIcpDomainsUsecase.execute', {
      searchTerm: this.searchTerm,
    });

    if (this.searchTerm.length === 0) {
      this.setSearchedIds([]);
      this.domainError('');
      this.domainMessage('');
    }

    const isValidUrl = validateUrl(this.searchTerm);

    await Promise.all([this.searchGlobal(), this.searchTenant()]);

    if (isValidUrl) {
      await this.validateDomain();
    }
    span.end();
  }
}

class AddSearchOrganizationOption {
  id: string;
  name: string;
  website: string | null;
  logoUrl: string | null;
  iconUrl: string | null;
  tenantOrganizationId: string;
  source: 'global' | 'tenant' = 'global';

  constructor(data: GlobalOrganization | Organization) {
    if (data instanceof Organization) {
      this.id = data.id;
      this.name = data.name;
      this.source = 'tenant';
      this.website = data.value.website ?? null;
      this.logoUrl = data.value.logoUrl ?? null;
      this.iconUrl = data.value.iconUrl ?? null;
      this.tenantOrganizationId = data.id;
    } else {
      this.id = data.id;
      this.name = data.name;
      this.source = 'global';
      this.website = data.website ?? null;
      this.logoUrl = data.logoUrl ?? null;
      this.iconUrl = data.iconUrl ?? null;
      this.tenantOrganizationId = data.organizationId ?? '';
    }
  }
}
