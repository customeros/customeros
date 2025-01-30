import { RootStore } from '@store/root';
import { OrganizationService } from '@domain/services';
import { action, computed, reaction, observable } from 'mobx';

import { SelectOption } from '@shared/types/SelectOptions';

export class EditOwnerOrganizationUsecase {
  @observable public accessor searchTerm = '';
  @observable public accessor isMenuOpen = false;
  @observable public accessor initialUsers: { label: string; value: string }[] =
    [];

  private root = RootStore.getInstance();
  private organizationService = new OrganizationService();
  private orgId: string;

  constructor(orgId: string) {
    this.orgId = orgId;
    this.execute = this.execute.bind(this);
    this.toggleMenu = this.toggleMenu.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.computeInitialOptions = this.computeInitialOptions.bind(this);

    reaction(
      () => this.root.users.tenantUsers.length,

      this.computeInitialOptions,
    );
    this.computeInitialOptions();
  }

  @computed
  get organization() {
    if (!this.orgId) return;

    return this.root.organizations.getById(this.orgId);
  }

  @computed
  get selectedUser() {
    return (
      this.organization?.owner && {
        value: this.organization?.owner?.id,
        label: this.organization?.owner?.name,
      }
    );
  }

  @action
  private computeInitialOptions() {
    this.initialUsers = this.root.users.tenantUsers
      .filter(
        (e) =>
          Boolean(e.value.firstName) ||
          Boolean(e.value.lastName) ||
          Boolean(e.value.name),
      )

      .map((user) => {
        return {
          value: user.id,
          label: user.name,
        };
      })
      .sort((a, b) => a.label.localeCompare(b.label));
  }

  @computed
  get userOptions() {
    return this.initialUsers.filter((user) =>
      user.label.toLowerCase().includes(this.searchTerm.toLowerCase()),
    );
  }

  @action
  public setSearchTerm(searchTerm: string) {
    this.searchTerm = searchTerm;
  }

  @action
  public toggleMenu(isOpen: boolean) {
    this.isMenuOpen = isOpen;
  }

  @action
  public reset() {
    this.setSearchTerm('');
  }

  @action
  public close() {
    this.reset();
    this.root.ui.commandMenu.setOpen(false);
  }

  @action
  execute(owner: SelectOption) {
    this.organizationService.editOwner(this.orgId, owner.value);
  }
}
