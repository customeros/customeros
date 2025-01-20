import { action, observable } from 'mobx';
import { RootStore } from '@store/root.ts';
import { Organization } from '@store/Organizations/Organization.dto.ts';
import { OrganizationsService } from '@store/Organizations/__service__/Organizations.service.ts';

export class AddOrganizationDomainCase {
  @observable accessor inputValue: string = '';
  @observable accessor error: string = '';
  @observable accessor associatedOrg: {
    id: string;
    name: string;
  } | null = null;

  @observable accessor isValidating: boolean = false;
  @observable accessor validationDetails: null | {
    primary: boolean;
    primaryDomain: string;
  } = null;
  @observable accessor entity: Organization | null = null;
  private root = RootStore.getInstance();
  private service = OrganizationsService.getInstance();

  constructor() {
    this.setInputValue = this.setInputValue.bind(this);
    this.reset = this.reset.bind(this);
  }

  @action
  setEntity(entity: Organization) {
    this.entity = entity;
  }

  @action
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
  }

  @action
  getInputValue() {
    return this.inputValue;
  }

  @action
  reset() {
    this.inputValue = '';
    this.associatedOrg = null;
    this.error = '';
  }

  @action
  resetValidation() {
    this.associatedOrg = null;
    this.error = '';
  }

  @action
  checkIfEmpty() {
    if (this.inputValue.trim() === '') {
      this.error = 'Houston, we have a blank...';

      return;
    }
    this.error = '';
  }

  @action
  async validateDomain() {
    this.isValidating = true;

    if (this.error) {
      this.isValidating = false;

      return;
    }

    try {
      const { checkDomain } = await this.service.checkDomain({
        domain: this.inputValue,
      });

      if (
        (checkDomain.primaryDomainOrganizationId &&
          checkDomain.primaryDomainOrganizationId !== this.entity?.id) ||
        (checkDomain.domainOrganizationId &&
          checkDomain.domainOrganizationId !== this.entity?.id)
      ) {
        this.error = 'Duplicate';
        this.associatedOrg = {
          name:
            checkDomain.primaryDomainOrganizationName ||
            checkDomain.domainOrganizationName ||
            '',
          id: (checkDomain.primaryDomainOrganizationId ||
            checkDomain.domainOrganizationId) as string,
        };

        return;
      }

      if (!checkDomain.validSyntax) {
        this.error = 'This domain appears to be invalid';

        return;
      }

      if (!checkDomain.accessible) {
        this.error = 'This domain is not reachable';

        return;
      }

      this.validationDetails = {
        primary: checkDomain.primary,
        primaryDomain: checkDomain.primaryDomain,
      };
    } catch (e) {
      this.error = `Validation failed`;
    } finally {
      this.isValidating = false;
    }
  }

  @action
  submit() {
    if (!this?.entity) return;

    this.entity?.draft();
    this.entity?.value.domainsDetails.push({
      domain: this.inputValue,
      primary: this.validationDetails?.primary || false,
      primaryDomain: this.validationDetails?.primaryDomain,
    });
    this.entity?.commit();
    this.inputValue = '';
    this.error = '';
  }
}
