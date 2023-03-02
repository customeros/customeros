import React from 'react';
import { Column } from '../../ui-kit/atoms/table/types';
import { IconButton, Times } from '../../ui-kit/atoms';
import { TableHeaderCell } from '../../ui-kit/atoms/table';
import { OrganizationTableCell } from './OrganizationTableCell';
import { FinderCell } from './FinderTableCell';
import { EmailTableCell } from './EmailTableCell';
import { AddressTableCell } from './AddressTableCell';
import { ActionColumn } from './ActionTableHeader';
import { OrganizationTableHeader } from './OrganizationTableHeader';

export const columns: Array<Column> = [
  {
    width: '25%',
    label: <OrganizationTableHeader />,
    template: (c: any) => {
      return <OrganizationTableCell organization={c?.organization} />;
    },
  },
  {
    width: '25%',
    label: 'Name',
    subLabel: 'Role',

    template: (c: any) => {
      if (!c?.contact) {
        return <span>-</span>;
      }
      const name = `${c?.contact?.firstName} ${c?.contact?.lastName} ${
        c?.contact?.name || ''
      }`;
      const displayName = name.trim().length ? name : 'Unnamed';
      return (
        <FinderCell
          label={displayName}
          subLabel={c?.contact.job}
          url={`/contact/${c?.contact.id}`}
        />
      );
    },
  },
  {
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
    width: '25%',
    label: 'Location',
    subLabel: 'City, State, Country',
    template: (c: any) => {
      return <AddressTableCell locations={c?.contact?.locations} />;
    },
  },
  {
    width: '10%',
    label: <ActionColumn />,
    subLabel: '',
    template: () => {
      return <div />;
    },
  },
];
