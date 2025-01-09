import { action, observable } from 'mobx';
import { Organization } from '@store/Organizations/Organization.dto.ts';

export class RemoveOrganizationDomainCase {
  @observable accessor isPrimary: boolean = false;
  @observable accessor domain: string = '';
  @observable accessor entity: Organization | null = null;

  constructor() {
    this.setIsPrimary = this.setIsPrimary.bind(this);
  }

  @action
  setEntity(entity: Organization) {
    this.entity = entity;
  }

  @action
  setIsPrimary(isPrimary: boolean) {
    this.isPrimary = isPrimary;
  }

  @action
  setDomain(domain: string) {
    this.domain = domain;
  }

  @action
  submit() {
    if (!this.entity) return;

    if (this.isPrimary) {
      this.entity.draft();

      // Remove all domains except the primary one
      for (let i = this.entity.value.domainsDetails.length - 1; i >= 0; i--) {
        if (this.entity.value.domainsDetails[i].primaryDomain === this.domain) {
          this.entity.value.domainsDetails.splice(i, 1);
        }
      }

      this.entity?.commit();
    } else {
      this.entity.draft();

      this.entity.value.domainsDetails.splice(
        this.entity.value.domainsDetails.findIndex(
          (e) => e.domain === this.domain,
        ),
        1,
      );
      this.entity?.commit();
    }

    this.isPrimary = false;
  }
}
