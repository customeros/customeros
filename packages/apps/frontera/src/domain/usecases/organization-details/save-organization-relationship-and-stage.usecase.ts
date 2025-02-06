import { action } from 'mobx';
import { OrganizationService } from '@domain/services';
import { OrganizationDatum } from '@store/Organizations/Organization.dto';

export class SaveOrganizationRelationshipAndStageUsecase {
  private orgService = new OrganizationService();
  private orgId: string;

  constructor(orgId: string) {
    this.orgId = orgId;
  }

  @action
  public execute(payload: Partial<OrganizationDatum>) {
    this.orgService.setRelationshipAndStage(payload, this.orgId);
  }
}
