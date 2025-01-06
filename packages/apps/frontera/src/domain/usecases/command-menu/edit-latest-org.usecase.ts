import { RootStore } from '@store/root';
import { reaction, observable } from 'mobx';
import { OrganizationService } from '@domain/services/organization/organizations.service';

export class EditLatestOrganizationActive {
  @observable accessor search: string = '';
  @observable private accessor searchedIds: string[] = [];
  private root = RootStore.getInstance();

  private service = new OrganizationService();

  constructor() {
    this.setSearchTerm = this.setSearchTerm.bind(this);

    reaction(() => this.search, this.executeSearch);
  }

  setSearchTerm(searchTerm: string) {
    this.search = searchTerm;
  }

  get searchTerm() {
    return this.search;
  }

  public async executeSearch() {
    await this.service.searchTenant(this.search);
  }
}
