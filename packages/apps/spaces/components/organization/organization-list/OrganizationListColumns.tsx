import React from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableCell } from '@spaces/atoms/table';
import {AddressTableCell, FinderCell} from "@spaces/finder/finder-table";

export const organizationListColumns: Array<Column> = [
  {
    id: 'finder-table-column-name-industry',
    width: '25%',
    label: <FinderCell label='Name' subLabel='Industry' />,

    template: (organization: any) => {
      return (
        <TableCell
          label={organization.name && organization.name !== '' ? organization.name : 'Unnamed'}
          subLabel={organization.industry}
          url={`/organization/${organization.id}`}
        />
      );
    },
  },
  {
    id: 'finder-table-column-domain-website',
    width: '25%',
    label: <FinderCell label='Domain' subLabel='Website' />,

    template: (organization: any) => {
      return (
        <TableCell
          label={organization.domain}
          subLabel={organization.website}
          url={`/organization/${organization.id}`}
        />
      );
    },
  },
  {
    id: 'finder-table-column-address',
    width: '50%',
    label: 'Location',
    subLabel: 'City, State, Country',
    template: (organization: any) => {
      return <AddressTableCell locations={organization.locations} />;
    },
  },
];
