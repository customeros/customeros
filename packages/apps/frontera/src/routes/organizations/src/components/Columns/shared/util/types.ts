import { ColumnDef } from '@tanstack/table-core';
import { ContactStore } from '@store/Contacts/Contact.store';
import { OrganizationStore } from '@store/Organizations/Organization.store';

export type MergedColumnDefs = ColumnDef<
  OrganizationStore | ContactStore,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  any
>[];
