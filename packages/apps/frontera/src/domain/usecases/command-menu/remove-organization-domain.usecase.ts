import { action, observable } from 'mobx';
import { OrganizationService } from '@domain/services';
import { Organization } from '@store/Organizations/Organization.dto.ts';

export class RemoveOrganizationDomainCase {
  @observable accessor domain: string = '';
  @observable accessor entity: Organization | null = null;
  private service = new OrganizationService();

  @action
  setEntity(entity: Organization) {
    this.entity = entity;
  }

  @action
  setDomain(domain: string) {
    this.domain = domain;
  }

  @action
  submit() {
    if (!this?.entity) return;

    this.service.removeDomain(this.entity, this.domain);
  }
}
