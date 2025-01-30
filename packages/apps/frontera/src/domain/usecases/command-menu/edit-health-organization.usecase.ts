import { action } from 'mobx';
import { OrganizationService } from '@domain/services';

import { OpportunityRenewalLikelihood } from '@shared/types/__generated__/graphql.types';

export class EditOrganizationHealthUseCase {
  private orgService = new OrganizationService();
  private orgId: string;

  constructor(orgId: string) {
    this.orgId = orgId;
  }

  @action
  public execute(
    renewalSummaryRenewalLikelihood: OpportunityRenewalLikelihood,
  ) {
    this.orgService.updateOrganizationHealth(
      this.orgId,
      renewalSummaryRenewalLikelihood,
    );
  }
}
