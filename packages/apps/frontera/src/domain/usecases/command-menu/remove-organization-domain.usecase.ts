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
    console.log('🏷️ ----- this.isPrimary: ', this.isPrimary);

    if (this.isPrimary) {
      this.entity.draft();

      this.entity.value.domainsDetails =
        this.entity.value.domainsDetails.filter(
          (e) => e.primaryDomain === this.domain,
        );
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
