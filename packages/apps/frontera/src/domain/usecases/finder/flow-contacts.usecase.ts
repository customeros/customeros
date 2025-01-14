import { RootStore } from '@store/root';
import { inPlaceSort } from 'fast-sort';
import { action, reaction, observable } from 'mobx';
import { Contact } from '@store/Contacts/Contact.dto';
import { indexAndSearch } from '@store/Contacts/__views__/util';
import { getContactSortFn } from '@store/Contacts/__views__/sortFns';
import { getContactFilterFns } from '@store/Contacts/__views__/filterFns';

export class FlowContactsUsecase {
  @observable private accessor flowId: string = '';
  private cache: Set<string> = new Set();
  private root = RootStore.getInstance();

  constructor() {
    reaction(() => {
      const preset = this.root.tableViewDefs.flowContactsPreset;

      if (!preset) return;
      if (!this.flowId) return;

      const viewDef = this.root.tableViewDefs.getById(preset);
      const flow = this.root.flows.value.get(this.flowId);

      const viewDefVersion = viewDef?.version;
      const flowVersion = flow?.version;

      return `${viewDefVersion}-${flowVersion}-${this.root.contacts.value.size}`;
    }, this.execute);

    // autorun(() => {
    //   this.execute(this.flowId, preset);
    // });
  }

  @action
  setFlowId(flowId: string) {
    this.flowId = flowId;
  }

  public execute = async () => {
    const flow = this.root.flows.value.get(this.flowId);
    const preset = this.root.tableViewDefs.flowContactsPreset;
    const viewDef = this.root.tableViewDefs.getById(preset!);

    if (!flow || !viewDef) return;

    const participants = flow.value.participants.map((p) => p.entityId);

    // const cacheEntry = this.makeCacheEntry(
    //   flowId,
    //   participants,
    //   JSON.stringify(viewDef.toRaw()),
    // );

    // console.log('cacheEntry', cacheEntry);

    // if (!this.cache.has(cacheEntry)) return;

    // this.cache.add(cacheEntry);
    await this.root.contacts.retrieve(participants);
    this.updateStoreView(this.flowId);
  };

  private updateStoreView(flowId: string) {
    const preset = this.root.tableViewDefs.flowContactsPreset;

    if (!preset) return;

    const viewDef = this.root.tableViewDefs.getById(preset);

    if (!viewDef) return;

    const defaultFilters = getContactFilterFns(viewDef.getDefaultFilters());
    const activeFilters = getContactFilterFns(viewDef.getFilters(), flowId);
    const sorting = JSON.parse(viewDef.value.sorting);

    this.root.contacts.setView(flowId, (data) => {
      const columnId = sorting?.id as string;
      const isDesc = sorting?.desc as boolean;

      const filteredIdsWithSortValues = (data as Contact[]).reduce(
        (acc, curr) => {
          if (!curr || !flowId) return acc;
          if (!curr.flowsIds?.includes(flowId)) return acc;

          if (
            defaultFilters.every((fn) => fn(curr)) &&
            activeFilters.every((fn) => fn(curr))
          ) {
            const sortValue = getContactSortFn(columnId)(curr);

            acc.push({ record: curr, sortValue });
          }

          return acc;
        },
        [] as {
          record: Contact;
          sortValue:
            | number
            | boolean
            | Date
            | null
            | (string | undefined)[]
            | string
            | undefined;
        }[],
      );

      let sorted = inPlaceSort(filteredIdsWithSortValues)
        [isDesc ? 'desc' : 'asc']((entry) => entry.sortValue)
        .map((entry) => entry.record);

      const searchTerm = this.root.contacts.getSearchTermByView(preset);

      if (searchTerm) {
        sorted = indexAndSearch(sorted, searchTerm);
      }

      return sorted;
    });
  }

  private makeCacheEntry(
    flowId: string,
    participants: string[],
    viewDef: string,
  ) {
    return [flowId, participants.join('-'), viewDef].join('-');
  }
}
