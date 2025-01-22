import set from 'lodash/set';
import omit from 'lodash/omit';
import merge from 'lodash/merge';
import { P, match } from 'ts-pattern';
import { Entity } from '@store/record';
import { Filter, FilterItem } from '@store/types';
import { action, computed, observable, runInAction } from 'mobx';

import {
  SortBy,
  TableIdType,
  TableViewType,
  SortingDirection,
} from '@graphql/types';

import { type TableViewDefStore } from './TableViewDef.store';
import { TableViewDefsQuery } from './__services__/getTableViewDefs.generated';

export type TableViewDefDatum = TableViewDefsQuery['tableViewDefs'][number];

export class TableViewDef extends Entity<TableViewDefDatum> {
  @observable accessor value: TableViewDefDatum = TableViewDef.default();

  constructor(store: TableViewDefStore, data: TableViewDefDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  @computed
  get id() {
    return this.value.id;
  }

  set id(value: string) {
    runInAction(() => {
      this.value.id = value;
    });
  }

  @computed
  get name() {
    return this.value.name;
  }

  public reorderColumn(sourceColumnId: number, targetColumnId: number) {
    this.draft();

    const fromIndex = this.value.columns.findIndex(
      (c) => c.columnId === sourceColumnId,
    );
    const toIndex = this.value.columns.findIndex(
      (c) => c.columnId === targetColumnId,
    );
    const column = this.value.columns[fromIndex];

    this.value.columns.splice(fromIndex, 1);
    this.value.columns.splice(toIndex, 0, column);

    this.commit();
  }

  public orderColumnsByVisibility() {
    this.draft();

    const prevLastVisibleIndex = [
      ...this.value.columns.map((c) => c.visible),
    ].lastIndexOf(true);

    const orderedColumns = this.value.columns.sort((a, b) => {
      if (a.visible === b.visible) return 0;
      if (a.visible) return -1;

      return 1;
    });

    const currentLastVisibleIndex = orderedColumns
      .map((c) => c.visible)
      .lastIndexOf(true);

    if (prevLastVisibleIndex === currentLastVisibleIndex) {
      this.commit();
    }

    this.value.columns.sort((a, b) => {
      if (a.visible === b.visible) return 0;
      if (a.visible) return -1;

      return 1;
    });

    this.commit();
  }

  @action
  public setColumnName(columnId: number, name: string) {
    const columnIdx = this.value.columns.findIndex(
      (c) => c.columnId === columnId,
    );

    this.draft();
    this.value.columns[columnIdx].name = name;
    this.commit();
  }

  @action
  public setColumnSize(columnType: string, size: number) {
    const columnIdx = this.value.columns.findIndex(
      (c) => c.columnType === columnType,
    );

    if (columnIdx !== -1) {
      this.draft();
      this.value.columns[columnIdx].width = size;
      this.commit();
    }

    // this.debouncedSave();
  }

  public getFilters() {
    try {
      return match(this.value.filters)
        .with(P.string.includes('AND'), (data) => JSON.parse(data))
        .otherwise(() => null);
    } catch (err) {
      console.error('Error parsing filters', err);

      return null;
    }
  }

  getDefaultFilters() {
    try {
      return match(this.value.defaultFilters)
        .with(P.string.includes('AND'), (data) => JSON.parse(data))
        .otherwise(() => null);
    } catch (err) {
      console.error('Error parsing default filters', err);

      return null;
    }
  }

  public getSorting() {
    try {
      return match(this.value.sorting)
        .with(P.string.includes('id'), (data) => JSON.parse(data))
        .otherwise(() => null);
    } catch (err) {
      console.error('Error parsing sorting', err);

      return null;
    }
  }

  public toSearchPayload(): { sort: SortBy; where: Filter | null } {
    const activeFilters =
      this.getFilters()
        ?.AND?.filter((f: Filter) => f?.filter && 'value' in f.filter)
        .map?.((f: Filter) => ({
          filter: omit(f.filter, 'active') as Filter['filter'],
        })) ?? [];

    const defaultFilters =
      this.getDefaultFilters()
        ?.AND?.filter((f: Filter) => f?.filter && 'value' in f.filter)
        .map?.((f: Filter) => ({
          filter: omit(f.filter, 'active') as Filter['filter'],
        })) ?? [];

    const where = {
      AND: [...defaultFilters, ...activeFilters],
    };

    const viewDefSorting = this.getSorting();
    const sort = {
      by: viewDefSorting.id,
      direction: viewDefSorting.desc
        ? SortingDirection.Desc
        : SortingDirection.Asc,
    };

    return { where, sort };
  }

  @action
  appendDefaultFilter(filter: FilterItem) {
    let draft = this.getDefaultFilters() as Filter;

    if (draft) {
      (draft as Filter).AND?.push({ filter });
    } else {
      draft = { AND: [{ filter }] };
    }

    this.draft();
    this.value.filters = JSON.stringify(draft);
    this.commit();
  }

  @action
  appendFilter(filter: FilterItem) {
    let draft = this.getFilters() as Filter;

    if (draft) {
      (draft as Filter).AND?.push({ filter });
    } else {
      draft = { AND: [{ filter }] };
    }

    this.draft();
    this.value.filters = JSON.stringify(draft);
    this.commit();
  }

  public getFilter(id: string) {
    const filters = this.getFilters();

    return (filters?.AND as Filter[])?.find((f) => f.filter?.property === id)
      ?.filter;
  }

  @action
  public setDefaultFilters(filter: FilterItem) {
    const draft = this.getDefaultFilters();

    if (!draft) {
      this.appendDefaultFilter({ ...filter, active: true });
    }
    const foundIndex = (draft.AND as Filter[])?.findIndex(
      (f) => f.filter?.property === filter.property,
    );

    if (foundIndex !== -1) {
      draft.AND[foundIndex].filter = filter;

      this.draft();
      this.value.filters = JSON.stringify(draft);
      this.commit();
    } else {
      this.appendDefaultFilter({ ...filter, active: true });
    }
  }

  @action
  public removeFilter(id: string, index?: number) {
    const draft = this.getFilters();

    if (draft) {
      if (index !== undefined) {
        draft.AND?.splice(index, 1);
      } else {
        draft.AND = (draft.AND as Filter[])?.filter(
          (f) => f.filter?.property !== id,
        );
      }

      this.draft();
      this.value.filters = JSON.stringify(draft);
      this.commit();
    }
  }

  @action
  public removeFilters() {
    this.draft();
    this.value.filters = JSON.stringify({ AND: [] });
    this.commit();
  }

  @action
  public toggleFilter(filter: FilterItem) {
    const draft = this.getFilters();

    if (draft) {
      const foundFilter = (draft.AND as Filter[])?.find(
        (f) => f.filter?.property === filter.property,
      )?.filter;

      if (foundFilter) {
        set(foundFilter, 'active', !filter?.active);

        this.draft();
        this.value.filters = JSON.stringify(draft);
        this.commit();
      } else {
        this.appendFilter({ ...filter, active: true });
      }
    }
  }

  @action
  public setFilterv2(filter: FilterItem, index: number) {
    const draft = this.getFilters();

    if (!draft) {
      this.appendFilter({ ...filter, active: true });

      return;
    }

    if (draft.AND && draft.AND[index]) {
      draft.AND[index].filter = filter;
    } else {
      draft.AND?.push({ filter });
    }

    this.draft();
    this.value.filters = JSON.stringify(draft);
    this.commit();
  }

  @action
  public setFilter(filter: FilterItem) {
    const draft = this.getFilters();

    if (!draft) {
      this.appendFilter({ ...filter, active: true });

      return;
    }
    const foundIndex = (draft.AND as Filter[])?.findIndex(
      (f) => f.filter?.property === filter.property,
    );

    if (foundIndex !== -1) {
      draft.AND![foundIndex]!.filter = filter;

      this.draft();
      this.value.filters = JSON.stringify(draft);
      this.commit();
    } else {
      this.appendFilter({ ...filter, active: true });
    }
  }

  @action
  public setPropertyFilter(property: string) {
    const draft = this.getFilters();

    if (!draft) {
      this.appendFilter({
        property,
        active: false,
        value: undefined,
      });

      return;
    }

    const foundIndex = (draft.AND as Filter[])?.findIndex(
      (f) => f.filter?.property === property,
    );

    if (foundIndex !== -1) {
      draft.AND![foundIndex].filter = { property, active: false };

      this.draft();
      this.value.filters = JSON.stringify(draft);
      this.commit();
    } else {
      this.appendFilter({
        property,
        active: false,
        value: undefined,
      });
    }
  }

  @action
  public setSorting(columndId: string, isDesc: boolean) {
    const draft = this.getSorting() as { id: string; desc: boolean };

    if (!draft) {
      this.draft();
      this.value.sorting = JSON.stringify({ id: columndId, desc: isDesc });
      this.commit();

      return;
    }

    draft.id = columndId;
    draft.desc = isDesc;

    this.draft();
    this.value.sorting = JSON.stringify(draft);
    this.commit();
  }

  public getPayloadToCopy = () => {
    return omit(this.value, 'id', 'createdAt', 'updatedAt');
  };

  public hasFilters() {
    return this.getFilters()?.AND?.length > 0;
  }

  public static default(
    payload?: Partial<TableViewDefDatum>,
  ): TableViewDefDatum {
    return merge(
      {
        __typename: 'TableViewDef',
        id: crypto.randomUUID(),
        name: '',
        columns: [],
        defaultFilters: '',
        filters: '',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        icon: '',
        isPreset: false,
        isShared: false,
        order: 0,
        sorting: '',
        tableId: TableIdType.Organizations,
        tableType: TableViewType.Organizations,
      },
      payload,
    );
  }
}
