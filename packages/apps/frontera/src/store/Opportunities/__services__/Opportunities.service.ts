import { Transport } from '@store/transport';

import GetOpportunityDocument from './getOpportunity.graphql';
import SaveOpportunityDocument from './saveOpportunity.graphql';
import GetOpportunitiesDocument from './getOpportunities.graphql';
import ArchiveOpportunityDocument from './archiveOpportunity.graphql';
import UpdateOpportunityOwnerDocument from './updateOpportunityOwner.graphql';
import UpdateOpportunityRenewalDocument from './updateOpportunityRenewal.graphql';
import {
  OpportunityQuery,
  OpportunityQueryVariables,
} from './getOpportunity.generated';
import {
  GetOpportunitiesQuery,
  GetOpportunitiesQueryVariables,
} from './getOpportunities.generated';
import {
  SaveOpportunityMutation,
  SaveOpportunityMutationVariables,
} from './saveOpportunity.generated';
import {
  ArchiveOpportunityMutation,
  ArchiveOpportunityMutationVariables,
} from './archiveOpportunity.generated';
import {
  UpdateOpportunityOwnerMutation,
  UpdateOpportunityOwnerMutationVariables,
} from './updateOpportunityOwner.generated';
import {
  UpdateOpportunityRenewalMutation,
  UpdateOpportunityRenewalMutationVariables,
} from './updateOpportunityRenewal.generated';

export class OpportunitiesService {
  private static instance: OpportunitiesService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  public static getInstance(transport: Transport): OpportunitiesService {
    if (!OpportunitiesService.instance) {
      OpportunitiesService.instance = new OpportunitiesService(transport);
    }

    return OpportunitiesService.instance;
  }

  async getOpportunities(
    variables: GetOpportunitiesQueryVariables,
  ): Promise<GetOpportunitiesQuery> {
    return this.transport.graphql.request<GetOpportunitiesQuery>(
      GetOpportunitiesDocument,
      variables,
    );
  }

  async getOpportunity(
    variables: OpportunityQueryVariables,
  ): Promise<OpportunityQuery> {
    return this.transport.graphql.request<OpportunityQuery>(
      GetOpportunityDocument,
      variables,
    );
  }

  async saveOpportunity(
    variables: SaveOpportunityMutationVariables,
  ): Promise<SaveOpportunityMutation> {
    return this.transport.graphql.request<
      SaveOpportunityMutation,
      SaveOpportunityMutationVariables
    >(SaveOpportunityDocument, variables);
  }

  async updateOpportunityRenewal(
    variables: UpdateOpportunityRenewalMutationVariables,
  ): Promise<UpdateOpportunityRenewalMutation> {
    return this.transport.graphql.request<
      UpdateOpportunityRenewalMutation,
      UpdateOpportunityRenewalMutationVariables
    >(UpdateOpportunityRenewalDocument, variables);
  }

  async updateOpportunityOwner(
    variables: UpdateOpportunityOwnerMutationVariables,
  ): Promise<UpdateOpportunityOwnerMutation> {
    return this.transport.graphql.request<
      UpdateOpportunityOwnerMutation,
      UpdateOpportunityOwnerMutationVariables
    >(UpdateOpportunityOwnerDocument, variables);
  }

  async archiveOpportunity(
    variables: ArchiveOpportunityMutationVariables,
  ): Promise<ArchiveOpportunityMutation> {
    return this.transport.graphql.request<
      ArchiveOpportunityMutation,
      ArchiveOpportunityMutationVariables
    >(ArchiveOpportunityDocument, variables);
  }
}
