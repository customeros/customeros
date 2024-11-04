import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { TableIdType, TableViewType, type TableViewDef } from '@graphql/types';

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
    const tableViewDefStore = this.value.get(id);

    if (!tableViewDefStore && this.isBootstrapped) {
      const defaultPresetId = this.defaultPreset;
      const navigateToDefaultPreset =
        window.location.pathname.includes('finder');

      if (defaultPresetId && navigateToDefaultPreset) {
        const defaultTableViewDefStore = this.value.get(defaultPresetId);

        if (defaultTableViewDefStore) {
          runInAction(() => {
            const url = new URL(window.location.href);

            url.searchParams.set('preset', defaultPresetId);
            window.history.replaceState(null, '', url.toString());
          });

          return defaultTableViewDefStore;
        }
      }
    }

    return tableViewDefStore;
  }

  toArray(): TableViewDefStore[] {
    return Array.from(this.value)?.flatMap(
      ([, tableViewDefStore]) => tableViewDefStore,
    );
  }

  get defaultPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Customers && t.value.isPreset,
    )?.value.id;
  }

  get opportunitiesPreset() {
    return this?.toArray().find(
      (t) =>
        t.value.tableId === TableIdType.Opportunities &&
        t.value.isShared &&
        t.value.isPreset,
    )?.value.id;
  }

  get opportunitiesTablePreset() {
    return this?.toArray().find(
      (t) =>
        t.value.tableId === TableIdType.OpportunitiesRecords &&
        !t.value.isShared &&
        t.value.isPreset,
    )?.value.id;
  }

  get targetsPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Targets && t.value.isPreset,
    )?.value.id;
  }

  get organizationsPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Organizations && t.value.isPreset,
    )?.value.id;
  }

  get upcomingInvoicesPreset() {
    return this?.toArray().find(
      (t) =>
        t.value.tableId === TableIdType.UpcomingInvoices && t.value.isPreset,
    )?.value.id;
  }

  get pastInvoicesPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.PastInvoices && t.value.isPreset,
    )?.value.id;
  }

  get contactsPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Contacts && t.value.isPreset,
    )?.value.id;
  }

  get contactsTargetPreset() {
    return this?.toArray().find(
      (t) =>
        t.value.tableId === TableIdType.ContactsForTargetOrganizations &&
        t.value.isPreset,
    )?.value.id;
  }

  get contractsPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Contracts && t.value.isPreset,
    )?.value.id;
  }

  get flowsPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.FlowActions && t.value.isPreset,
    )?.value.id;
  }

  get flowContactsPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.FlowContacts && t.value.isPreset,
    )?.value.id;
  }

  createFavorite = async (
    {
      id,
      isShared,
      name,
    }: {
      id: string;
      name?: string;
      isShared: boolean;
    },
    options?: { onSuccess?: (serverId: string) => void },
  ) => {
    const favoritePreset = this.getById(id)?.getPayloadToCopy();

    const newTableViewDef = new TableViewDefStore(this.root, this.transport);

    newTableViewDef.value = {
      ...getDefaultValue(),
      ...favoritePreset,
      name: name
        ? name
        : `Copy of ${
            favoritePreset?.tableType === TableViewType.Invoices
              ? ` ${favoritePreset?.name} Invoices`
              : favoritePreset?.name
          }`,
      isPreset: false,
      isShared,
    };

    const {
      id: _id,
      createdAt,
      updatedAt,
      defaultFilters,
      ...payload
    } = newTableViewDef.value;

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
        }, 100);
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
      defaultFilters
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
