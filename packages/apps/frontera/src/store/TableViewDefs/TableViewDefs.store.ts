import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { TableIdType, type TableViewDef } from '@graphql/types';

import mock from './mock.json';
import { getDefaultValue, TableViewDefStore } from './TableViewDef.store';
import { TableViewDefsService } from './__services__/TableViewDef.service';

export class TableViewDefsStore implements GroupStore<TableViewDef> {
  value: Map<string, TableViewDefStore> = new Map();
  isLoading = false;
  channel?: Channel;
  version: number = 0;
  history: GroupOperation[] = [];
  isBootstrapped = false;
  error: string | null = null;
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<TableViewDef>();
  private service: TableViewDefsService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = TableViewDefsService.getInstance(transport);
    makeAutoSyncableGroup(this, {
      channelName: 'TableViewDefs',
      ItemStore: TableViewDefStore,
      getItemId: (item) => item.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped) return;

    if (this.root.demoMode) {
      this.load(mock.data.tableViewDefs as TableViewDef[]);
      this.isBootstrapped = true;

      return;
    }

    try {
      this.isLoading = true;

      const res =
        await this.transport.graphql.request<TABLE_VIEW_DEFS_QUERY_RESULT>(
          TABLE_VIEW_DEFS_QUERY,
        );

      this.load(res?.tableViewDefs);
      runInAction(() => {
        this.isBootstrapped = true;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async invalidate() {
    try {
      this.isLoading = true;

      const { tableViewDefs } = await this.service.getTableViewDefs();

      this.load(tableViewDefs);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  getById(id: string) {
    return this.value.get(id);
  }

  toArray(): TableViewDefStore[] {
    return Array.from(this.value)?.flatMap(
      ([, tableViewDefStore]) => tableViewDefStore,
    );
  }

  get defaultPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Customers,
    )?.value.id;
  }

  get opportunitiesPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Opportunities && t.value.isShared,
    )?.value.id;
  }

  get leadsPreset() {
    return this?.toArray().find((t) => t.value.tableId === TableIdType.Leads)
      ?.value.id;
  }

  get targetsPreset() {
    return this?.toArray().find((t) => t.value.tableId === TableIdType.Nurture)
      ?.value.id;
  }

  get churnedPreset() {
    return this?.toArray().find((t) => t.value.tableId === TableIdType.Churn)
      ?.value.id;
  }

  get addressBookPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Organizations,
    )?.value.id;
  }

  createFavorite = async (
    favoritePresetId: string,
    options?: { onSuccess?: (serverId: string) => void },
  ) => {
    const favoritePreset = this.getById(favoritePresetId)?.getPayloadToCopy();

    const newTableViewDef = new TableViewDefStore(this.root, this.transport);

    newTableViewDef.value = {
      ...getDefaultValue(),
      ...favoritePreset,
      isPreset: false,
    };

    const { id: _id, createdAt, updatedAt, ...payload } = newTableViewDef.value;

    const tempId = newTableViewDef.id;
    let serverId = '';

    this.value.set(tempId, newTableViewDef);
    this.isLoading = true;

    try {
      const { tableViewDef_Create } = await this.service.createTableViewDef({
        input: {
          ...payload,
        },
      });

      runInAction(() => {
        serverId = tableViewDef_Create.id;

        newTableViewDef.value.id = serverId;

        this.value.set(serverId, newTableViewDef);
        this.value.delete(tempId);

        this.sync({
          action: 'APPEND',
          ids: [serverId],
        });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      this.isLoading = false;

      if (serverId) {
        setTimeout(() => {
          this.invalidate();
          options?.onSuccess?.(serverId);
        }, 1000);
      }
    }
  };

  archive = async (id: string, options?: { onSuccess?: () => void }) => {
    this.isLoading = true;

    const viewName = this.getById(id)?.value.name;

    try {
      const { tableViewDef_Archive } = await this.service.archiveTableViewDef({
        id,
      });

      if (tableViewDef_Archive.accepted) {
        runInAction(() => {
          this.value.delete(id);

          this.sync({
            action: 'DELETE',
            ids: [id],
          });
        });
        this.root.ui.toastSuccess(
          `${viewName} is now archived`,
          'archive-view-success',
        );
      }
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
        this.root.ui.toastError(
          `We couldn't archive ${viewName} view`,
          'archive-view-error',
        );
      });
    } finally {
      this.isLoading = false;
      options?.onSuccess?.();
    }
  };
}

type TABLE_VIEW_DEFS_QUERY_RESULT = { tableViewDefs: TableViewDef[] };
const TABLE_VIEW_DEFS_QUERY = gql`
  query tableViewDefs {
    tableViewDefs {
      id
      name
      tableType
      tableId
      order
      icon
      filters
      sorting
      columns {
        columnId
        columnType
        name
        width
        visible
        filter
      }
      isPreset
      isShared
      createdAt
      updatedAt
    }
  }
`;
