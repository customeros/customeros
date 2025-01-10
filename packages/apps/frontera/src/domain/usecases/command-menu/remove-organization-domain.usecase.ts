import { action, observable } from 'mobx';
import { Organization } from '@store/Organizations/Organization.dto.ts';

export class RemoveOrganizationDomainCase {
  @observable accessor domain: string = '';
  @observable accessor entity: Organization | null = null;

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

    this.entity.draft();

    this.entity.value.domainsDetails.splice(
      this.entity.value.domainsDetails.findIndex(
        (e) => e.domain === this.domain,
      ),
      1,
    );
    this.entity?.commit();
  }
}
