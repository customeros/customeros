import { Transport } from '@store/transport';

import GetOpportunityDocument from './getOpportunity.graphql';
import GetOpportunitiesDocument from './getOpportunities.graphql';
import CreateOpportunityDocument from './createOpportunity.graphql';
import UpdateOpportunityDocument from './updateOpportunity.graphql';
import UpdateOpportunityOwnerDocument from './updateOpportunityOwner.graphql';
import UpdateOpportunityRenewalDocument from './updateOpportunityRenewal.graphql';
import UpdateOpportunityToCloseWonDocument from './updateOpportunityToCloseWon.graphql';
import UpdateOpportunityToCloseLostDocument from './updateOpportunityToCloseLost.graphql';
import {
  OpportunityQuery,
  OpportunityQueryVariables,
} from './getOpportunity.generated';
import {
  GetOpportunitiesQuery,
  GetOpportunitiesQueryVariables,
} from './getOpportunities.generated';
import {
  CreateOpportunityMutation,
  CreateOpportunityMutationVariables,
} from './createOpportunity.generated';
import {
  UpdateOpportunityMutation,
  UpdateOpportunityMutationVariables,
} from './updateOpportunity.generated';
import {
  UpdateOpportunityOwnerMutation,
  UpdateOpportunityOwnerMutationVariables,
} from './updateOpportunityOwner.generated';
import {
  UpdateOpportunityRenewalMutation,
  UpdateOpportunityRenewalMutationVariables,
} from './updateOpportunityRenewal.generated';
import {
  UpdateOpportunityToCloseWonMutation,
  UpdateOpportunityToCloseWonMutationVariables,
} from './updateOpportunityToCloseWon.generated';
import {
  UpdateOpportunityToCloseLostMutation,
  UpdateOpportunityToCloseLostMutationVariables,
} from './updateOpportunityToCloseLost.generated';

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

  async createOpportunity(
    variables: CreateOpportunityMutationVariables,
  ): Promise<CreateOpportunityMutation> {
    return this.transport.graphql.request<
      CreateOpportunityMutation,
      CreateOpportunityMutationVariables
    >(CreateOpportunityDocument, variables);
  }

  async updateOpportunity(
    variables: UpdateOpportunityMutationVariables,
  ): Promise<UpdateOpportunityMutation> {
    return this.transport.graphql.request<
      UpdateOpportunityMutation,
      UpdateOpportunityMutationVariables
    >(UpdateOpportunityDocument, variables);
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

  async updateOpportunityToCloseLost(
    variables: UpdateOpportunityToCloseLostMutationVariables,
  ): Promise<UpdateOpportunityToCloseLostMutation> {
    return this.transport.graphql.request<
      UpdateOpportunityToCloseLostMutation,
      UpdateOpportunityToCloseLostMutationVariables
    >(UpdateOpportunityToCloseLostDocument, variables);
  }

  async updateOpportunityToCloseWon(
    variables: UpdateOpportunityToCloseWonMutationVariables,
  ): Promise<UpdateOpportunityToCloseWonMutation> {
    return this.transport.graphql.request<
      UpdateOpportunityToCloseWonMutation,
      UpdateOpportunityToCloseWonMutationVariables
    >(UpdateOpportunityToCloseWonDocument, variables);
  }
}
