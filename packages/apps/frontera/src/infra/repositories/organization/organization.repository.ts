import { Transport } from '@infra/transport';

import AddTagDocument from './mutations/addTag.graphql';
import AddDomainDocument from './mutations/addDomain.graphql';
import RemoveTagDocument from './mutations/removeTag.graphql';
import CheckDomainDocument from './queries/checkDomain.graphql';
import CheckWebsiteDocument from './queries/checkWebsite.graphql';
import RemoveDomainDocument from './mutations/removeDomain.graphql';
import RemoveSocialDocument from './mutations/removeSocial.graphql';
import AddSubsidiaryDocument from './mutations/addSubsidiary.graphql';
import RemoveDomainsDocument from './mutations/removeDomains.graphql';
import FlagWrongFieldDocument from './mutations/flagWrongField.graphql';
import GetOrganizationsDocument from './queries/getOrganizations.graphql';
import RemoveSubsidiaryDocument from './mutations/removeSubsidiary.graphql';
import SaveOrganizationDocument from './mutations/saveOrganization.graphql';
import HideOrganizationsDocument from './mutations/hideOrganizations.graphql';
import SearchOrganizationsDocument from './queries/searchOrganizations.graphql';
import ImportOrganizationDocument from './mutations/importOrganization.graphql';
import MergeOrganizationsDocument from './mutations/mergeOrganizations.graphql';
import GetOrganizationsByIdsDocument from './queries/getOrganizationsByIds.graphql';
import UpdateOnboardingStatusDocument from './mutations/updateOnboardingStatus.graphql';
import SearchGlobalOrganizationsDocument from './queries/searchGlobalOrganizations.graphql';
import GetArchivedOrganizationsAfterDocument from './queries/getArchivedOrganizations.graphql';
import {
  CheckDomainQuery,
  CheckDomainQueryVariables,
} from './queries/checkDomain.generated';
import UpdateAllOpportunityRenewalsDocument from './mutations/updateAllOpportunityRenewals.graphql';
import {
  AddDomainMutation,
  AddDomainMutationVariables,
} from './mutations/addDomain.generated';
import {
  CheckWebsiteQuery,
  CheckWebsiteQueryVariables,
} from './queries/checkWebsite.generated';
import {
  RemoveSocialMutation,
  RemoveSocialMutationVariables,
} from './mutations/removeSocial.generated';
import {
  RemoveDomainMutation,
  RemoveDomainMutationVariables,
} from './mutations/removeDomain.generated';
import {
  RemoveDomainsMutation,
  RemoveDomainsMutationVariables,
} from './mutations/removeDomains.generated';
import {
  GetOrganizationsQuery,
  GetOrganizationsQueryVariables,
} from './queries/getOrganizations.generated';
import {
  FlagWrongFieldMutation,
  FlagWrongFieldMutationVariables,
} from './mutations/flagWrongField.generated';
import {
  AddTagsToOrganizationMutation,
  AddTagsToOrganizationMutationVariables,
} from './mutations/addTag.generated';
import {
  SaveOrganizationMutation,
  SaveOrganizationMutationVariables,
} from './mutations/saveOrganization.generated';
import {
  SearchOrganizationsQuery,
  SearchOrganizationsQueryVariables,
} from './queries/searchOrganizations.generated';
import {
  HideOrganizationsMutation,
  HideOrganizationsMutationVariables,
} from './mutations/hideOrganizations.generated';
import {
  MergeOrganizationsMutation,
  MergeOrganizationsMutationVariables,
} from './mutations/mergeOrganizations.generated';
import {
  ImportOrganizationMutation,
  ImportOrganizationMutationVariables,
} from './mutations/importOrganization.generated';
import {
  GetOrganizationsByIdsQuery,
  GetOrganizationsByIdsQueryVariables,
} from './queries/getOrganizationsByIds.generated';
import {
  RemoveTagFromOrganizationMutation,
  RemoveTagFromOrganizationMutationVariables,
} from './mutations/removeTag.generated';
import {
  UpdateOnboardingStatusMutation,
  UpdateOnboardingStatusMutationVariables,
} from './mutations/updateOnboardingStatus.generated';
import {
  AddSubsidiaryToOrganizationMutation,
  AddSubsidiaryToOrganizationMutationVariables,
} from './mutations/addSubsidiary.generated';
import {
  SearchGlobalOrganizationsQuery,
  SearchGlobalOrganizationsQueryVariables,
} from './queries/searchGlobalOrganizations.generated';
import {
  GetArchivedOrganizationsAfterQuery,
  GetArchivedOrganizationsAfterQueryVariables,
} from './queries/getArchivedOrganizations.generated';
import {
  RemoveSubsidiaryToOrganizationMutation,
  RemoveSubsidiaryToOrganizationMutationVariables,
} from './mutations/removeSubsidiary.generated';
import {
  BulkUpdateOpportunityRenewalMutation,
  BulkUpdateOpportunityRenewalMutationVariables,
} from './mutations/updateAllOpportunityRenewals.generated';

export class OrganizationRepository {
  static instance: OrganizationRepository | null = null;
  private transport = Transport.getInstance();

  public static getInstance() {
    if (!OrganizationRepository.instance) {
      OrganizationRepository.instance = new OrganizationRepository();
    }

    return OrganizationRepository.instance;
  }

  async checkDomain(payload: CheckDomainQueryVariables) {
    return this.transport.graphql.request<
      CheckDomainQuery,
      CheckDomainQueryVariables
    >(CheckDomainDocument, payload);
  }

  async removeDomain(payload: RemoveDomainMutationVariables) {
    return this.transport.graphql.request<
      RemoveDomainMutation,
      RemoveDomainMutationVariables
    >(RemoveDomainDocument, payload);
  }

  async addDomain(payload: AddDomainMutationVariables) {
    return this.transport.graphql.request<
      AddDomainMutation,
      AddDomainMutationVariables
    >(AddDomainDocument, payload);
  }

  async removeDomains(payload: RemoveDomainsMutationVariables) {
    return this.transport.graphql.request<
      RemoveDomainsMutation,
      RemoveDomainsMutationVariables
    >(RemoveDomainsDocument, payload);
  }

  async searchOrganizations(payload: SearchOrganizationsQueryVariables) {
    return this.transport.graphql.request<
      SearchOrganizationsQuery,
      SearchOrganizationsQueryVariables
    >(SearchOrganizationsDocument, payload);
  }

  async importOrganization(payload: ImportOrganizationMutationVariables) {
    return this.transport.graphql.request<
      ImportOrganizationMutation,
      ImportOrganizationMutationVariables
    >(ImportOrganizationDocument, payload);
  }

  async searchGlobalOrganizations(
    payload: SearchGlobalOrganizationsQueryVariables,
  ) {
    return this.transport.graphql.request<
      SearchGlobalOrganizationsQuery,
      SearchGlobalOrganizationsQueryVariables
    >(SearchGlobalOrganizationsDocument, payload);
  }

  async getOrganizations(payload: GetOrganizationsQueryVariables) {
    return this.transport.graphql.request<
      GetOrganizationsQuery,
      GetOrganizationsQueryVariables
    >(GetOrganizationsDocument, payload);
  }

  async getOrganizationsByIds(payload: GetOrganizationsByIdsQueryVariables) {
    return this.transport.graphql.request<
      GetOrganizationsByIdsQuery,
      GetOrganizationsByIdsQueryVariables
    >(GetOrganizationsByIdsDocument, payload);
  }

  async getOrganization(id: string) {
    const { ui_organizations } = await this.getOrganizationsByIds({
      ids: [id],
    });

    return ui_organizations[0];
  }

  async checkWebsite(payload: CheckWebsiteQueryVariables) {
    return this.transport.graphql.request<
      CheckWebsiteQuery,
      CheckWebsiteQueryVariables
    >(CheckWebsiteDocument, payload);
  }

  async getArchivedOrganizationsAfter(
    payload: GetArchivedOrganizationsAfterQueryVariables,
  ) {
    return this.transport.graphql.request<
      GetArchivedOrganizationsAfterQuery,
      GetArchivedOrganizationsAfterQueryVariables
    >(GetArchivedOrganizationsAfterDocument, payload);
  }

  async saveOrganization(payload: SaveOrganizationMutationVariables) {
    return this.transport.graphql.request<
      SaveOrganizationMutation,
      SaveOrganizationMutationVariables
    >(SaveOrganizationDocument, payload);
  }

  async hideOrganizations(payload: HideOrganizationsMutationVariables) {
    return this.transport.graphql.request<
      HideOrganizationsMutation,
      HideOrganizationsMutationVariables
    >(HideOrganizationsDocument, payload);
  }

  async mergeOrganizations(payload: MergeOrganizationsMutationVariables) {
    return this.transport.graphql.request<
      MergeOrganizationsMutation,
      MergeOrganizationsMutationVariables
    >(MergeOrganizationsDocument, payload);
  }

  async flagWrongField(payload: FlagWrongFieldMutationVariables) {
    return this.transport.graphql.request<
      FlagWrongFieldMutation,
      FlagWrongFieldMutationVariables
    >(FlagWrongFieldDocument, payload);
  }

  async addTag(payload: AddTagsToOrganizationMutationVariables) {
    return this.transport.graphql.request<
      AddTagsToOrganizationMutation,
      AddTagsToOrganizationMutationVariables
    >(AddTagDocument, payload);
  }

  async removeTag(payload: RemoveTagFromOrganizationMutationVariables) {
    return this.transport.graphql.request<
      RemoveTagFromOrganizationMutation,
      RemoveTagFromOrganizationMutationVariables
    >(RemoveTagDocument, payload);
  }

  async removeSocial(payload: RemoveSocialMutationVariables) {
    return this.transport.graphql.request<
      RemoveSocialMutation,
      RemoveSocialMutationVariables
    >(RemoveSocialDocument, payload);
  }

  async updateAllOpportunityRenewals(
    payload: BulkUpdateOpportunityRenewalMutationVariables,
  ) {
    return this.transport.graphql.request<
      BulkUpdateOpportunityRenewalMutation,
      BulkUpdateOpportunityRenewalMutationVariables
    >(UpdateAllOpportunityRenewalsDocument, payload);
  }

  async addSubsidiary(payload: AddSubsidiaryToOrganizationMutationVariables) {
    return this.transport.graphql.request<
      AddSubsidiaryToOrganizationMutation,
      AddSubsidiaryToOrganizationMutationVariables
    >(AddSubsidiaryDocument, payload);
  }

  async removeSubsidiary(
    payload: RemoveSubsidiaryToOrganizationMutationVariables,
  ) {
    return this.transport.graphql.request<
      RemoveSubsidiaryToOrganizationMutation,
      RemoveSubsidiaryToOrganizationMutationVariables
    >(RemoveSubsidiaryDocument, payload);
  }

  async updateOnboardingStatus(
    payload: UpdateOnboardingStatusMutationVariables,
  ) {
    return this.transport.graphql.request<
      UpdateOnboardingStatusMutation,
      UpdateOnboardingStatusMutationVariables
    >(UpdateOnboardingStatusDocument, payload);
  }
}
