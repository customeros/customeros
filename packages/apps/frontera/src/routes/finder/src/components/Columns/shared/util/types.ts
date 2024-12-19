import { ColumnDef } from '@tanstack/table-core';
import { Contact } from '@store/Contacts/Contact.dto';
import { InvoiceStore } from '@store/Invoices/Invoice.store';
import { Organization } from '@store/Organizations/Organization.dto';

export type MergedColumnDefs = ColumnDef<
  Organization | Contact | InvoiceStore,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  any
>[];
