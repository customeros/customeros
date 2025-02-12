import type { RootStore } from '@store/root';
import type { Transport } from '@infra/transport';

import { runInAction } from 'mobx';
import { Store } from '@store/_store';

import { TableIdType, TableViewType } from '@graphql/types';

import { TableViewDef, type TableViewDefDatum } from './TableViewDef.dto';
import { TableViewDefsService } from './__services__/TableViewDef.service';

export class TableViewDefStore extends Store<TableViewDefDatum, TableViewDef> {
  private service: TableViewDefsService = TableViewDefsService.getInstance();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'TableViewDefs',
      getId: (data) => data.id,
      factory: TableViewDef,
    });

    this.hydrate();
  }

  get defaultPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.Organizations && t.value.isPreset,
    )?.value.id;
  }

  get customersPreset() {
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

  get contactsFlowsPreset() {
    return this?.toArray().find(
      (t) => t.value.tableId === TableIdType.FlowContacts,
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

  get customPresets() {
    return this.toArray().filter((p) => !p.value.isShared && !p.value.isPreset);
  }

  get teamPresets() {
    return this.toArray().filter(
      (p) =>
        p.value.isShared &&
        !p.value.isPreset &&
        p.value.tableId !== TableIdType.FlowContacts,
    );
  }

  get flowContactsPresets() {
    return this.toArray().filter(
      (p) => p.value.tableId === TableIdType.FlowContacts,
    );
  }

  public getById(id: string) {
    const tableViewDefStore = this.value.get(id);

    if (!tableViewDefStore && this.isBootstrapped) {
      const defaultPresetId = this.defaultPreset;
      const navigateToDefaultPreset =
        window.location.pathname.includes('finder');

      if (defaultPresetId && navigateToDefaultPreset) {
        const defaultTableViewDefStore = this.value.get(defaultPresetId);

        if (defaultTableViewDefStore) {
          const url = new URL(window.location.href);

          url.searchParams.set('preset', defaultPresetId);
          window.history.replaceState(null, '', url.toString());

          return defaultTableViewDefStore;
        }
      }
    }

    return tableViewDefStore ?? null;
  }

  public getActivePreset() {
    const url = new URLSearchParams(window.location.search);
    const preset = url.get('preset');

    return preset ?? null;
  }

  public async bootstrap() {
    if (this.isBootstrapped) return;

    try {
      this.isLoading = true;

      const { tableViewDefs } = await this.service.getTableViewDefs();

      runInAction(() => {
        tableViewDefs.forEach((raw) => {
          this.value.set(raw.id, new TableViewDef(this, raw));
        });
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

  public async invalidate() {
    try {
      this.isLoading = true;

      const { tableViewDefs } = await this.service.getTableViewDefs();

      runInAction(() => {
        tableViewDefs.forEach((raw) => {
          this.value.set(raw.id, new TableViewDef(this, raw));
        });
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

  public createFavorite = async (
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

    const newTableViewDef = new TableViewDef(
      this,
      TableViewDef.default({
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
        filters: '',
        defaultFilters: favoritePreset?.filters || '',
      }),
    );

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
    runInAction(() => {
      this.isLoading = true;
    });

    try {
      const { tableViewDef_Create } = await this.service.createTableViewDef({
        input: {
          ...payload,
          defaultFilters: newTableViewDef.value.defaultFilters,
        },
      });

      runInAction(() => {
        serverId = tableViewDef_Create.id;

        newTableViewDef.value.id = serverId;

        this.value.set(serverId, newTableViewDef);
        this.value.delete(tempId);
        this.version++;

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
      runInAction(() => {
        this.isLoading = false;
      });

      if (serverId) {
        options?.onSuccess?.(serverId);
      }
    }
  };

  public archive = async (id: string, options?: { onSuccess?: () => void }) => {
    runInAction(() => {
      this.isLoading = true;
    });

    const viewName = this.getById(id)?.value.name;

    try {
      const { tableViewDef_Archive } = await this.service.archiveTableViewDef({
        id,
      });

      if (tableViewDef_Archive.accepted) {
        runInAction(() => {
          this.value.delete(id);
          this.version++;

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
      });
      this.root.ui.toastError(
        `We couldn't archive ${viewName} view`,
        'archive-view-error',
      );
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
      options?.onSuccess?.();
    }
  };
}
