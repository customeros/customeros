import { action, reaction, computed, observable } from 'mobx';
import { OrganizationService } from '@domain/services/organization/organizations.service';

export class EditLatestOrganizationActive {
  @observable accessor search: string = '';

  private service = new OrganizationService();

  constructor() {
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.executeSearch = this.executeSearch.bind(this);

    reaction(() => this.search, this.executeSearch);
  }

  @action
  setSearchTerm(searchTerm: string) {
    this.search = searchTerm;
  }

  @computed
  get searchTerm() {
    return this.search;
  }

  public async executeSearch() {
    await this.service.searchTenant(this.search);
  }
}
