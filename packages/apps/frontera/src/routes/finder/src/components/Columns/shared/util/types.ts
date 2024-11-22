import { ColumnDef } from '@tanstack/table-core';
import { ContactStore } from '@store/Contacts/Contact.store';
import { InvoiceStore } from '@store/Invoices/Invoice.store';
import { Organization } from '@store/Organizations/Organization.dto';

export type MergedColumnDefs = ColumnDef<
  Organization | ContactStore | InvoiceStore,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  any
>[];
