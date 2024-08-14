import merge from 'lodash/merge';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { SyncableGroup } from '@store/syncable-group';
import {
  when,
  action,
  computed,
  override,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import {
  relationshipStageMap,
  stageRelationshipMap,
  validRelationshipsForStage,
} from '@utils/orgStageAndRelationshipStatusMap.ts';
import {
  Tag,
  Filter,
  SortBy,
  Pagination,
  Organization,
  SortingDirection,
  OrganizationInput,
  OrganizationStage,
  OrganizationRelationship,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

import mock from './mock.json';
import {
  getDefaultValue,
  OrganizationStore,
  ORGANIZATION_QUERY,
  ORGANIZATION_QUERY_RESULT,
} from './Organization.store';

export class OrganizationsStore extends SyncableGroup<
  Organization,
  OrganizationStore
> {
  totalElements = 0;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, OrganizationStore);

    makeObservable(this, {
      maxLtv: computed,
      isFullyLoaded: computed,
      hide: action.bound,
      merge: action.bound,
      create: action.bound,
      channelName: override,
      updateStage: action.bound,
      totalElements: observable,
    });

    when(
      () =>
        this.isBootstrapped && this.totalElements > 0 && !this.root.demoMode,
      async () => {
        await this.bootstrapRest();
      },
    );
  }

  get channelName() {
    return 'Organizations';
  }

  get maxLtv() {
    return Math.max(
      ...this.toArray().map(
        (org) => Math.round(org.value.accountDetails?.ltv ?? 0) + 1,
      ),
    );
  }

  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  async bootstrapStream() {
    try {
      await this.transport.stream<Organization>('/organizations', {
        onData: (data) => {
          runInAction(() => {
            this.load([data], { getId: (data) => data.metadata.id });
          });
        },
      });

      runInAction(() => {
        this.isBootstrapped = true;
      });
    } catch (e) {
      runInAction(() => {
        console.error(e);
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(
        mock.data.dashboardView_Organizations
          .content as unknown as Organization[],
        { getId: (data) => data.metadata.id },
      );
      this.totalElements = mock.data.dashboardView_Organizations.totalElements;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { dashboardView_Organizations } =
        await this.transport.graphql.request<
          ORGANIZATIONS_QUERY_RESPONSE,
          ORGANIZATIONS_QUERY_PAYLOAD
        >(ORGANIZATIONS_QUERY, {
          pagination: { limit: 1000, page: 0 },
          sort: {
            by: 'LAST_TOUCHPOINT',
            caseSensitive: false,
            direction: SortingDirection.Desc,
          },
        });

      this.load(dashboardView_Organizations.content, {
        getId: (data) => data.metadata.id,
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = dashboardView_Organizations.totalElements;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async bootstrapRest() {
    let page = 1;

    while (this.totalElements > this.value.size) {
      try {
        const { dashboardView_Organizations } =
          await this.transport.graphql.request<
            ORGANIZATIONS_QUERY_RESPONSE,
            ORGANIZATIONS_QUERY_PAYLOAD
          >(ORGANIZATIONS_QUERY, {
            pagination: { limit: 1000, page },
            sort: {
              by: 'LAST_TOUCHPOINT',
              caseSensitive: false,
              direction: SortingDirection.Desc,
            },
          });

        runInAction(() => {
          page++;
          this.load(dashboardView_Organizations.content, {
            getId: (data) => data.metadata.id,
          });
        });
      } catch (e) {
        runInAction(() => {
          this.error = (e as Error)?.message;
        });
        break;
      }
    }
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray<T extends OrganizationStore>(
    compute: (arr: OrganizationStore[]) => T[],
  ) {
    const arr = this.toArray();

    return compute(arr);
  }

  async create(
    payload?: OrganizationInput,
    options?: { onSucces?: (serverId: string) => void },
  ) {
    const newOrganization = new OrganizationStore(
      this.root,
      this.transport,
      merge(getDefaultValue(), payload),
    );
    const tempId = newOrganization.id;
    let serverId = '';

    this.value.set(tempId, newOrganization);
    this.isLoading = true;

    try {
      const { organization_Create } = await this.transport.graphql.request<
        CREATE_ORGANIZATION_RESPONSE,
        CREATE_ORGANIZATION_PAYLOAD
      >(CREATE_ORGANIZATION_MUTATION, {
        input: {
          website: payload?.website ?? '',
          name: payload?.name ?? 'Unnamed',
          relationship: newOrganization.value.relationship,
          stage: newOrganization.value.stage,
        },
      });

      runInAction(() => {
        serverId = organization_Create.metadata.id;

        newOrganization.setId(serverId);

        this.value.set(serverId, newOrganization);
        this.value.delete(tempId);

        this.sync({
          action: 'APPEND',
          ids: [serverId],
        });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      this.isLoading = false;

      if (serverId) {
        // Invalidate the cache after 1 second to allow the server to process the data
        // invalidating immediately would cause the server to return the organization data without
        // lastTouchpoint properties populated
        setTimeout(() => {
          this.value.get(serverId)?.invalidate();
          options?.onSucces?.(serverId);
        }, 1000);
      }
    }
  }

  async hide(ids: string[]) {
    ids.forEach((id) => {
      this.value.delete(id);
    });

    try {
      this.isLoading = true;
      await this.transport.graphql.request<unknown, HIDE_ORGANIZATIONS_PAYLOAD>(
        HIDE_ORGANIZATIONS_MUTATION,
        { ids },
      );

      runInAction(() => {
        this.sync({ action: 'DELETE', ids });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async merge(
    primaryId: string,
    mergeIds: string[],
    callback?: (id: string) => void,
  ) {
    mergeIds.forEach((id) => {
      this.value.delete(id);
    });
    callback?.(primaryId);

    try {
      this.isLoading = true;
      await this.transport.graphql.request<
        unknown,
        MERGE_ORGANIZATIONS_PAYLOAD
      >(MERGE_ORGANIZATIONS_MUTATION, {
        primaryOrganizationId: primaryId,
        mergedOrganizationIds: mergeIds,
      });

      runInAction(() => {
        this.sync({ action: 'DELETE', ids: mergeIds });
        this.sync({ action: 'INVALIDATE', ids: mergeIds });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  updateTags = (ids: string[], tags: Tag[]) => {
    ids.forEach((id) => {
      this.value.get(id)?.update((organization) => {
        const organizationTagIds = new Set(
          (organization.tags ?? []).map((tag) => tag.id),
        );
        const filteredTags = tags.filter(
          (tag) => !organizationTagIds.has(tag.id),
        );

        organization.tags = [...(organization.tags ?? []), ...filteredTags];

        return organization;
      });
    });
  };

  removeTags = (ids: string[]) => {
    ids.forEach((id) => {
      this.value.get(id)?.update((organization) => {
        organization.tags = [];

        return organization;
      });
    });
  };

  updateStage = (ids: string[], stage: OrganizationStage, mutate = true) => {
    let invalidCustomerStageCount = 0;

    ids.forEach((id) => {
      this.value.get(id)?.update(
        (org) => {
          const currentRelationship = org.relationship;
          const newDefaultRelationship = stageRelationshipMap[stage];
          const validRelationships = validRelationshipsForStage[stage];

          if (
            currentRelationship &&
            validRelationships?.includes(currentRelationship)
          ) {
            org.stage = stage;
          } else if (
            currentRelationship === OrganizationRelationship.Customer
          ) {
            invalidCustomerStageCount++;

            return org; // Do not update if current relationship is Customer and new stage is not valid
          } else {
            org.stage = stage;
            org.relationship = newDefaultRelationship || org.relationship;
          }

          return org;
        },
        { mutate: mutate },
      );
    });

    if (invalidCustomerStageCount) {
      this.root.ui.toastError(
        `${invalidCustomerStageCount} customer${
          invalidCustomerStageCount > 1 ? 's' : ''
        } remain unchanged`,
        'stage-update-failed-due-to-relationship-mismatch',
      );
    }
  };

  updateRelationship = (
    ids: string[],
    relationship: OrganizationRelationship,
    mutate = true,
  ) => {
    let invalidCustomerStageCount = 0;

    ids.forEach((id) => {
      this.value.get(id)?.update(
        (org) => {
          if (
            org.relationship === OrganizationRelationship.Customer &&
            ![
              OrganizationRelationship.FormerCustomer,
              OrganizationRelationship.NotAFit,
            ].includes(relationship)
          ) {
            invalidCustomerStageCount++;

            return org; // Do not update if current is customer and new is not formet customer or not a fit
          }
          org.relationship = relationship;
          org.stage = relationshipStageMap[org.relationship];

          return org;
        },
        { mutate: mutate },
      );
    });

    if (invalidCustomerStageCount) {
      this.root.ui.toastError(
        `${invalidCustomerStageCount} customer${
          invalidCustomerStageCount > 1 ? 's' : ''
        } remain unchanged`,
        'stage-update-failed-due-to-relationship-mismatch',
      );
    }
  };

  updateHealth = (
    ids: string[],
    health: OpportunityRenewalLikelihood,
    mutate = true,
  ) => {
    ids.forEach((id) => {
      this.value.get(id)?.update(
        (org) => {
          org.accountDetails = {
            ...org.accountDetails,
            renewalSummary: {
              ...(org.accountDetails?.renewalSummary ?? {}),
              renewalLikelihood: health,
            },
          };

          return org;
        },
        { mutate: mutate },
      );
    });
  };

  async getById(id: string) {
    try {
      this.isLoading = true;

      const { organization } = await this.transport.graphql.request<
        ORGANIZATION_QUERY_RESULT,
        { id: string }
      >(ORGANIZATION_QUERY, { id });
      const newOrganization = new OrganizationStore(
        this.root,
        this.transport,
        merge(getDefaultValue(), organization),
      );

      this.value.set(organization.metadata.id, newOrganization);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}

type ORGANIZATIONS_QUERY_PAYLOAD = {
  sort?: SortBy;
  where?: Filter;
  pagination: Pagination;
};
type ORGANIZATIONS_QUERY_RESPONSE = {
  dashboardView_Organizations: {
    totalElements: number;
    totalAvailable: number;
    content: Organization[];
  };
};
const ORGANIZATIONS_QUERY = gql`
  query getOrganizations(
    $pagination: Pagination!
    $where: Filter
    $sort: SortBy
  ) {
    dashboardView_Organizations(
      pagination: $pagination
      where: $where
      sort: $sort
    ) {
      content {
        name
        note
        notes

        metadata {
          id
          created
        }
        contracts {
          metadata {
            id
          }
        }

        parentCompanies {
          organization {
            metadata {
              id
            }
            name
          }
        }
        owner {
          id
          firstName
          lastName
          name
        }
        contacts(pagination: { page: 0, limit: 100 }) {
          content {
            id
            metadata {
              id
            }
          }
        }
        stage
        description
        industry
        market
        website
        domains
        isCustomer
        logo
        icon
        relationship
        lastFundingRound
        leadSource
        valueProposition
        slackChannelId
        public

        socialMedia {
          id
          url
          followersCount
        }
        employees
        tags {
          id
          name
          createdAt
          updatedAt
          source
          appSource
          metadata {
            id
            created
            lastUpdated
            source
            sourceOfTruth
            appSource
          }
        }
        yearFounded
        accountDetails {
          ltv
          churned
          renewalSummary {
            arrForecast
            maxArrForecast
            renewalLikelihood
            nextRenewalDate
          }
          onboarding {
            status
            comments
            updatedAt
          }
        }
        locations {
          id
          name
          country
          region
          locality
          zip
          street
          postalCode
          houseNumber
          rawAddress
          locality
          countryCodeA2
          countryCodeA3
        }
        subsidiaries {
          organization {
            metadata {
              id
            }
            name
            parentCompanies {
              organization {
                name
                metadata {
                  id
                }
              }
            }
          }
        }
        parentCompanies {
          organization {
            metadata {
              id
            }
          }
        }
        lastTouchpoint {
          lastTouchPointTimelineEventId
          lastTouchPointAt
          lastTouchPointType
          lastTouchPointTimelineEvent {
            __typename
            ... on PageView {
              id
            }
            ... on Issue {
              id
              createdAt
              updatedAt
            }
            ... on LogEntry {
              id
              createdBy {
                lastName
                firstName
              }
            }
            ... on Note {
              id
              createdBy {
                firstName
                lastName
              }
            }
            ... on InteractionEvent {
              id
              channel
              eventType
              externalLinks {
                type
              }
              sentBy {
                __typename
                ... on EmailParticipant {
                  type
                  emailParticipant {
                    id
                    email
                    rawEmail
                  }
                }
                ... on ContactParticipant {
                  contactParticipant {
                    id
                    name
                    firstName
                    lastName
                  }
                }
                ... on JobRoleParticipant {
                  jobRoleParticipant {
                    contact {
                      id
                      name
                      firstName
                      lastName
                    }
                  }
                }
                ... on UserParticipant {
                  userParticipant {
                    id
                    firstName
                    lastName
                  }
                }
              }
            }
            ... on Analysis {
              id
            }
            ... on Meeting {
              id
              name
              attendedBy {
                __typename
              }
            }
            ... on Action {
              id
              actionType
              createdAt
              source
              actionType
              createdBy {
                id
                firstName
                lastName
              }
            }
          }
        }

        contracts {
          metadata {
            id
          }
        }
      }
      totalElements
      totalAvailable
    }
  }
`;
type CREATE_ORGANIZATION_PAYLOAD = {
  input: OrganizationInput;
};
type CREATE_ORGANIZATION_RESPONSE = {
  organization_Create: {
    metadata: {
      id: string;
    };
  };
};
const CREATE_ORGANIZATION_MUTATION = gql`
  mutation createOrganization($input: OrganizationInput!) {
    organization_Create(input: $input) {
      metadata {
        id
      }
    }
  }
`;
type HIDE_ORGANIZATIONS_PAYLOAD = {
  ids: string[];
};
const HIDE_ORGANIZATIONS_MUTATION = gql`
  mutation hideOrganizations($ids: [ID!]!) {
    organization_HideAll(ids: $ids) {
      result
    }
  }
`;
type MERGE_ORGANIZATIONS_PAYLOAD = {
  primaryOrganizationId: string;
  mergedOrganizationIds: string[];
};
const MERGE_ORGANIZATIONS_MUTATION = gql`
  mutation mergeOrganizations(
    $primaryOrganizationId: ID!
    $mergedOrganizationIds: [ID!]!
  ) {
    organization_Merge(
      primaryOrganizationId: $primaryOrganizationId
      mergedOrganizationIds: $mergedOrganizationIds
    ) {
      id
    }
  }
`;
