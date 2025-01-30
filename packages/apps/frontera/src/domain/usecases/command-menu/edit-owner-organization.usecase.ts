import { action } from 'mobx';
import { OrganizationService } from '@domain/services';

export class EditOwnerUseCase {
  private organizationService = new OrganizationService();
  private orgId: string;

  constructor(orgId: string) {
    this.orgId = orgId;
  }

  @action
  execute(owner: string) {
    this.organizationService.editOwner(this.orgId, owner);
  }
}
