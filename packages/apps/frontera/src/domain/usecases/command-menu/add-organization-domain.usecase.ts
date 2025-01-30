import { action, observable } from 'mobx';
import { OrganizationService } from '@domain/services';
import { Organization } from '@store/Organizations/Organization.dto';
import { OrganizationRepository } from '@infra/repositories/organization/organization.repository.ts';

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
  private service = new OrganizationService();
  private repository = OrganizationRepository.getInstance();

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
      const { checkDomain } = await this.repository.checkDomain({
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

      if (!checkDomain.allowedForOrganization) {
        this.error = 'We do not allow adding this domain';

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
    if (!this?.entity) {
      console.error('AddOrganizationDomainCase: No entity provided');

      return;
    }

    this.service.addDomain(this.entity, {
      domain: this.inputValue,
      primary: this.validationDetails?.primary || false,
      primaryDomain: this.validationDetails?.primaryDomain,
    });
    this.inputValue = '';
    this.error = '';
  }
}
