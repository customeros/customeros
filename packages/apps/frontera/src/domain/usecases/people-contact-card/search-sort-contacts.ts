import { action, observable } from 'mobx';

export class SearchSortContact {
  private static instance: SearchSortContact;
  @observable accessor search: string = '';
  @observable accessor sortBy: 'Created' | 'First name' | 'Updated' | 'Tenure' =
    'Created';
  @observable accessor sortDirection: 'asc' | 'desc' = 'asc';

  constructor() {
    this.setSearch = this.setSearch.bind(this);
    this.setSort = this.setSort.bind(this);
    this.setSortDirection = this.setSortDirection.bind(this);
  }

  public static getInstance(): SearchSortContact {
    if (!SearchSortContact.instance) {
      SearchSortContact.instance = new SearchSortContact();
    }

    return new SearchSortContact();
  }

  @action
  setSearch(search: string) {
    this.search = search;
  }

  @action
  setSort(sort: 'Created' | 'First name' | 'Updated' | 'Tenure') {
    this.sortBy = sort;
  }

  @action
  setSortDirection(direction: 'asc' | 'desc') {
    this.sortDirection = direction;
  }

  getSearch() {
    return this.search;
  }

  getSort() {
    return this.sortBy;
  }

  getSortDirection() {
    return this.sortDirection;
  }
}
