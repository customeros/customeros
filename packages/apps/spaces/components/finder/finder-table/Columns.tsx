import React from 'react';
import { Column } from '../../ui-kit/atoms/table/types';
import { OrganizationTableCell } from './OrganizationTableCell';
import { EmailTableCell } from './EmailTableCell';
import { AddressTableCell } from './AddressTableCell';
import { ActionColumn } from './ActionTableHeader';
import { FinderMergeItemTableHeader } from './FinderMergeItemTableHeader';
import { ContactTableCell } from './ContactTableCell';

export const columns: Array<Column> = [
  {
    id: 'finder-table-column-org',
    width: '25%',
    label: (
      <FinderMergeItemTableHeader
        mergeMode='MERGE_ORG'
        label='Organization'
        subLabel='Industry'
      />
    ),
    template: (c: any) => {
      return <OrganizationTableCell organization={c?.organization} />;
    },
  },
  {
    id: 'finder-table-column-contact',
    width: '25%',
    label: (
      <FinderMergeItemTableHeader
        mergeMode='MERGE_CONTACT'
        label='Name'
        subLabel='Role'
      />
    ),

    template: (c: any) => <ContactTableCell contact={c?.contact} />,
  },
  {
    id: 'finder-table-column-email',
    width: '25%',
    label: 'Email',
    template: (c: any) => {
      if (!c?.contact) {
        return <span>-</span>;
      }
      return <EmailTableCell emails={c.contact?.emails} />;
    },
  },
  {
    id: 'finder-table-column-address',
    width: '25%',
    label: 'Location',
    subLabel: 'City, State, Country',
    template: (c: any) => {
      return <AddressTableCell locations={c?.contact?.locations} />;
    },
  },
  {
    id: 'finder-table-column-actions',
    width: '10%',
    label: <ActionColumn />,
    subLabel: '',
    template: () => {
      return <div style={{ display: 'none' }} />;
    },
  },
];
