import type { Transport } from '@infra/transport';

import { match } from 'ts-pattern';

import type { Operation } from './types';

import { RootStore } from './root';
import { ContactService } from './Contacts/__service__/Contacts.service';
import { InvoicesService } from './Invoices/__service__/Invoices.service';
import { JobRolesService } from './JobRoles/__service__/JobRoles.service';
import { MailboxesService } from './Settings/__service__/Mailboxes/Mailboxes.service';
import { TableViewDefsService } from './TableViewDefs/__services__/TableViewDef.service';
import { CustomFieldsService } from './Settings/__service__/CustomFields/CustomFields.service';

export class GraphqlService {
  private customFieldsService: CustomFieldsService;
  private invoiceService: InvoicesService;
  private contactService: ContactService;
  private mailboxService: MailboxesService;
  private tableViewDefService: TableViewDefsService;
  private jobRolesService: JobRolesService;

  constructor(private root: RootStore, private transport: Transport) {
    this.customFieldsService = CustomFieldsService.getInstance(this.transport);
    this.invoiceService = InvoicesService.getInstance(this.transport);
    this.contactService = ContactService.getInstance();
    this.mailboxService = MailboxesService.getInstance(this.transport);
    this.tableViewDefService = TableViewDefsService.getInstance();
    this.getStore = this.getStore.bind(this);
    this.jobRolesService = JobRolesService.getInstance();
  }

  public async mutate(operation: Operation) {
    if (!operation.entityId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }

    return match(operation.entity)
      .with('customFields', async () => {
        const store = this.getStore(operation, 'customFields');

        if (!store) return;

        return await this.customFieldsService.mutateOperation(operation, store);
      })
      .with('Invoices', async () => {
        const store = this.getStore(operation, 'invoices');

        if (!store) return;

        return await this.invoiceService.mutateOperation(operation, store);
      })
      .with('Contacts', async () => {
        const store = this.getStore(operation, 'contacts');

        if (!store) return;

        return await this.contactService.mutateOperation(operation, store);
      })
      .with('Mailboxes', async () => {
        const store = this.getStore(operation, 'mailboxes');

        if (!store) return;

        return await this.mailboxService.mutateOperation(operation, store);
      })
      .with('TableViewDefs', async () => {
        const store = this.getStore(operation, 'tableViewDefs');

        if (!store) return;

        return await this.tableViewDefService.mutateOperation(operation, store);
      })

      .otherwise(() => {});
  }

  // TODO: see if we need to get an active record instead of a dto
  private getStore(operation: Operation, storePath: keyof RootStore) {
    let store;

    if (operation.entity === 'Organizations') {
      // @ts-expect-error no issue
      store = this.root[storePath]?.getById(operation.entityId);
    } else {
      // @ts-expect-error no issue
      store = this.root[storePath]?.value?.get(operation.entityId);
    }

    if (!store) {
      console.error(
        `Store with id ${operation.entityId} not found. Mutations will not fire`,
      );

      return null;
    }

    return store;
  }
}
