import { action, computed, observable } from 'mobx';

type GqlOperation = {
  id: string;
  name: string;
  variables: Record<string, unknown> | null;
};

type GqlResponse = {
  id: string;
  errors: string | null;
  data: Record<string, unknown> | null;
};

type View = 'operations' | 'store' | 'settings' | 'traces';

export class DevtoolsStore {
  @observable accessor gqlOperations: GqlOperation[] = [];
  @observable accessor gqlResponses: Map<string, GqlResponse> = new Map();
  @observable accessor openGqlOperationId: string | null = null;
  @observable accessor operationsSearchTerm = '';
  @observable accessor responseSearchTerm = '';
  @observable accessor view: View = 'operations';
  @observable accessor detailedStore: string | null = null;
  @observable accessor detailedEntityId: string | null = null;

  visibleStores = [
    'agents',
    'contacts',
    'contracts',
    'common.slackChannels',
    'flows',
    'jobRoles',
    'organizations',
    'tableViewDefs',
    'tags',
  ];

  constructor() {}

  @computed
  get filteredGqlOperations() {
    return this.gqlOperations.filter((op) =>
      op.name.toLowerCase().includes(this.operationsSearchTerm.toLowerCase()),
    );
  }

  @action
  addGqlOp(operation: GqlOperation) {
    this.gqlOperations.push(observable.object(operation));
  }

  @action
  addGqlRes(response: GqlResponse) {
    this.gqlResponses.set(response.id, observable.object(response));
  }

  @action
  toggleGqlOp(id?: string) {
    if (this.openGqlOperationId === id) {
      this.openGqlOperationId = null;

      return;
    }
    this.openGqlOperationId = id ?? null;
  }

  @action
  searchOperations(term: string) {
    this.operationsSearchTerm = term;
  }

  @action
  searchResponse(term: string) {
    this.responseSearchTerm = term;
  }

  @action
  toggleView(view: View) {
    this.view = view;
  }

  @action
  toggleStore(store: string | null) {
    this.detailedStore = store;
  }

  @action
  toggleEntity(id: string | null) {
    this.detailedEntityId = this.detailedEntityId === id ? null : id;
  }
}
