import { Transport } from '@infra/transport';

import AddDomainDocument from './mutations/addDomain.graphql';
import CheckDomainDocument from './queries/checkDomain.graphql';
import RemoveDomainDocument from './mutations/removeDomain.graphql';
import RemoveDomainsDocument from './mutations/removeDomains.graphql';
import {
  CheckDomainQuery,
  CheckDomainQueryVariables,
} from './queries/checkDomain.generated';
import {
  AddDomainMutation,
  AddDomainMutationVariables,
} from './mutations/addDomain.generated';
import {
  RemoveDomainMutation,
  RemoveDomainMutationVariables,
} from './mutations/removeDomain.generated';
import {
  RemoveDomainsMutation,
  RemoveDomainsMutationVariables,
} from './mutations/removeDomains.generated';

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
}
