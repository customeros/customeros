import { match } from 'ts-pattern';
import { RootStore } from '@store/root';
import { when, action, reaction, computed, observable } from 'mobx';
import { Organization } from '@store/Organizations/Organization.dto';
import { OrganizationRepository } from '@infra/repositories/organization';
import { SearchGlobalOrganizationsQuery } from '@infra/repositories/organization/queries/searchGlobalOrganizations.generated';

import { validateUrl } from '@utils/url';
import {
  SortingDirection,
  OrganizationStage,
  ComparisonOperator,
  OrganizationRelationship,
} from '@graphql/types';

type GlobalOrganization =
  SearchGlobalOrganizationsQuery['globalOrganizations_Search'][number];

type Status =
  | 'LOADING'
  | 'ERROR'
  | 'IDLE'
  | 'VALIDATING_DOMAIN'
  | 'ADDING_ORGANIZATION';

export class AddSearchOrganizationsUsecase {
  private root = RootStore.getInstance();
  private service = OrganizationRepository.getInstance();
  @observable private accessor searchedIds: string[] = [];

  @observable accessor searchTerm = '';
  @observable accessor errorMessage = '';
  @observable accessor status: Status = 'IDLE';
  @observable accessor addingOrganizationId = '';
  @observable accessor domainValidationError = '';
  @observable accessor domainValidationMessage = '';
  @observable accessor options: AddSearchOrganizationOption[] = [];

  constructor() {
    this.execute = this.execute.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.searchGlobal = this.searchGlobal.bind(this);

    // get first 30 global suggestions on mount
    when(() => !!this.root.session.sessionToken, this.searchGlobal);
    reaction(() => this.searchTerm, this.execute);
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
  get isAddingOrganization() {
    return this.status === 'ADDING_ORGANIZATION';
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
  private addingOrganization(id: string) {
    this.addingOrganizationId = id;
    this.status = 'ADDING_ORGANIZATION';
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
    this.addingOrganizationId = '';
    this.domainValidationError = '';
    this.domainValidationMessage = '';
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

  private getViewDefDefaults() {
    const searchParams = new URLSearchParams(window.location.search);
    const preset = searchParams.get('preset');

    if (!preset) return {};

    const viewDef = this.root.tableViewDefs.getById(preset);

    if (!viewDef) return {};

    return match(viewDef.getDefaultFilters())
      .with(null, () => ({}))
      .with({ AND: [] }, () => ({}))
      .with(
        {
          AND: [
            {
              filter: {
                property: 'ORGANIZATIONS_RELATIONSHIP',
                value: ['CUSTOMER'],
              },
            },
          ],
        },
        () => ({
          relationship: OrganizationRelationship.Customer,
          stage: OrganizationStage.Onboarding,
        }),
      )
      .with(
        {
          AND: [
            {
              filter: {
                property: 'ORGANIZATIONS_STAGE',
                value: ['TARGET'],
              },
            },
            {
              filter: {
                property: 'ORGANIZATIONS_RELATIONSHIP',
                value: ['PROSPECT'],
              },
            },
          ],
        },
        () => ({
          relationship: OrganizationRelationship.Prospect,
          stage: OrganizationStage.Target,
        }),
      )
      .otherwise(() => ({}));
  }

  public async addNewOrganization(option?: AddSearchOrganizationOption) {
    if (!option) {
      const isUrl = validateUrl(this.searchTerm);

      await this.root.organizations.create({
        ...this.getViewDefDefaults(),
        name: isUrl ? undefined : this.searchTerm || 'Unnamed',
        website: !isUrl ? undefined : this.searchTerm || '',
      });
    }

    if (option?.source !== 'global') {
      this.root.ui.commandMenu.toggle('AddNewOrganization');
      this.reset();
      this.searchGlobal();

      return;
    }

    if (option) {
      const { id: globalId, ...rest } = option;

      this.addingOrganization(globalId);

      await this.root.organizations.import(Number(globalId), {
        ...this.getViewDefDefaults(),
        ...rest,
      });
      this.searchGlobal();
    }

    this.root.ui.commandMenu.toggle('AddNewOrganization');
    this.reset();
  }

  public async execute() {
    if (this.searchTerm.length === 0) {
      this.setSearchedIds([]);
      this.domainError('');
      this.domainMessage('');
    }

    const isValidUrl = validateUrl(this.searchTerm);

    // fetch global and tenant organizations
    await Promise.all([this.searchGlobal(), this.searchTenant()]);

    if (isValidUrl) {
      // if term is a valid URL, validate and check domain
      await this.validateDomain();
    }
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
