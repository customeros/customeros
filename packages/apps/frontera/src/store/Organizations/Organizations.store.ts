import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { when, runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import {
  Filter,
  SortBy,
  Pagination,
  Organization,
  SortingDirection,
  OrganizationInput,
  OrganizationStage,
} from '@graphql/types';

import { OrganizationStore } from './Organization2.store';

export class OrganizationsStore implements GroupStore<Organization> {
  version = 0;
  isLoading = false;
  totalElements = 0;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, OrganizationStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Organization>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Organizations',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: OrganizationStore,
    });

    when(
      () => this.isBootstrapped && this.totalElements > 0,
      async () => {
        await this.bootstrapRest();
      },
    );
  }

  async bootstrap() {
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

      this.load(dashboardView_Organizations.content);
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
          this.load(dashboardView_Organizations.content);
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

  toComputedArray<T extends Store<Organization>>(
    compute: (arr: Store<Organization>[]) => T[],
  ) {
    const arr = this.toArray();

    return compute(arr);
  }

  create = async (
    payload?: OrganizationInput,
    options?: { onSucces?: (serverId: string) => void },
  ) => {
    const newOrganization = new OrganizationStore(this.root, this.transport);
    const tempId = newOrganization.value.metadata.id;
    let serverId = '';

    if (payload) {
      merge(newOrganization.value, payload);
    }

    this.value.set(tempId, newOrganization);

    try {
      const { organization_Create } = await this.transport.graphql.request<
        CREATE_ORGANIZATION_RESPONSE,
        CREATE_ORGANIZATION_PAYLOAD
      >(CREATE_ORGANIZATION_MUTATION, {
        input: {
          name: 'Unnamed',
          relationship: newOrganization.value.relationship,
          stage: newOrganization.value.stage,
        },
      });

      runInAction(() => {
        serverId = organization_Create.metadata.id;

        newOrganization.value.metadata.id = serverId;

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
  };

  hide = async (ids: string[]) => {
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
  };

  merge = async (primaryId: string, mergeIds: string[]) => {
    mergeIds.forEach((id) => {
      this.value.delete(id);
    });

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
  };

  updateStage = (ids: string[], stage: OrganizationStage) => {
    ids.forEach((id) => {
      this.value.get(id)?.update((org) => {
        org.stage = stage;

        return org;
      });
    });
  };
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
        metadata {
          id
          created
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
        socialMedia {
          id
          url
        }
        employees
        yearFounded
        accountDetails {
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
