import { ColumnDef } from '@tanstack/table-core';
import { ContactStore } from '@store/Contacts/Contact.store';
import { InvoiceStore } from '@store/Invoices/Invoice.store';
import { OrganizationStore } from '@store/Organizations/Organization.store';

export type MergedColumnDefs = ColumnDef<
  OrganizationStore | ContactStore | InvoiceStore,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  any
>[];
