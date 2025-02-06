import { action } from 'mobx';
import { OrganizationService } from '@domain/services';

import { OrganizationStage } from '@shared/types/__generated__/graphql.types';

export class EditOrganizationStageUseCase {
  private orgService = new OrganizationService();
  private orgId: string;

  constructor(orgId: string) {
    this.orgId = orgId;
  }

  @action
  public execute(stage: OrganizationStage) {
    this.orgService.setRelationshipAndStage({ stage: stage }, this.orgId);
  }
}
