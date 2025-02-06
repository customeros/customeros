import { action } from 'mobx';
import { OrganizationService } from '@domain/services';

import { OrganizationRelationship } from '@shared/types/__generated__/graphql.types';

export class EditOrganizationRelationshipUseCase {
  private orgService = new OrganizationService();
  private orgId: string;

  constructor(orgId: string) {
    this.orgId = orgId;
  }

  @action
  public execute(relationship: OrganizationRelationship) {
    this.orgService.setRelationshipAndStage(
      { relationship: relationship },
      this.orgId,
    );
  }
}
