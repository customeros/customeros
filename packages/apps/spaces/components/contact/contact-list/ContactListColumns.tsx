import React from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableCell } from '@spaces/atoms/table';
import { getContactDisplayName } from '../../../utils';
import {AddressTableCell, EmailTableCell, FinderCell} from "@spaces/finder/finder-table";

export const contactListColumns: Array<Column> = [
  {
    id: 'finder-table-column-contact',
    width: '25%',
    label: <FinderCell label='Name' />,

    template: (contact: any) => {
      return (
        <TableCell
          label={getContactDisplayName(contact)}
          url={`/contact/${contact.id}`}
        />
      );
    },
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
];
