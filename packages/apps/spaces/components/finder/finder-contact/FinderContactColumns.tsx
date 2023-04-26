import React from 'react';
import { Column } from '../../ui-kit/atoms/table/types';
import { FinderCell } from '../finder-table';
import { TableCell } from '../../ui-kit/atoms/table';
import { getContactDisplayName } from '../../../utils';

export const finderContactColumns: Array<Column> = [
  {
    id: 'finder-table-column-contact',
    width: '25%',
    label: <FinderCell label='Name' />,

    template: (contact: any) => {
      console.log(contact);
      return (
        <TableCell
          label={getContactDisplayName(contact)}
          url={`/contact/${contact.id}`}
        />
      );
    },
  },
];
