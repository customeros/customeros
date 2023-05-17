import React from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableCell } from '@spaces/atoms/table';
import {
  ActionColumn,
  AddressTableCell,
  ContactTableCell,
  EmailTableCell,
  FinderMergeItemTableHeader,
} from '@spaces/finder/finder-table';
import { uuid4 } from '@sentry/utils';

export const contactListColumns: Array<Column> = [
  {
    id: 'finder-table-column-contact',
    width: '25%',
    label: (
      <FinderMergeItemTableHeader
        mergeMode='MERGE_CONTACT'
        label='Name'
        subLabel={''}
      />
    ),
    template: (contact: any) => {
      return <ContactTableCell contact={contact} />;
    },
  },
  {
    id: 'finder-table-column-email',
    width: '20%',
    label: 'Email',
    template: (c: any) => {
      if (!c?.contact) {
        return <span key={uuid4()}>-</span>;
      }
      return <EmailTableCell emails={c.contact?.emails} />;
    },
  },
  {
    id: 'finder-table-organization-position',
    width: '25%',
    label: 'Organization',
    subLabel: 'Position',
    template: (c: any) => {
      if (!c.jobRoles || c.jobRoles.length === 0) {
        return <span>-</span>;
      }
      return (
        <TableCell
          label={c.jobRoles[0].organization?.name ?? 'Unnamed'}
          subLabel={c.jobRoles[0].jobTitle ?? '-'}
          url={`/organization/${c.jobRoles[0].organization.id}`}
        />
      );
    },
  },
  {
    id: 'finder-table-column-org',
    width: '25%',
    label: 'Location',
    subLabel: 'City, State, Country',
    template: (c: any) => {
      return <AddressTableCell locations={c?.contact?.locations} />;
    },
  },
  {
    id: 'finder-table-column-actions',
    width: '5%',
    label: <ActionColumn scope={'MERGE_CONTACT'} />,
    subLabel: '',
    template: () => {
      return <div style={{ display: 'none' }} />;
    },
  },
];
