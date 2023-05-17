import React from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableCell } from '@spaces/atoms/table';
import {
  ActionColumn,
  AddressTableCell,
  ContactTableCell,
  FinderCell,
  FinderMergeItemTableHeader,
  OrganizationTableCell,
} from '@spaces/finder/finder-table';

export const organizationListColumns: Array<Column> = [
  {
    id: 'finder-table-column-organization-name',
    width: '25%',
    label: (
      <FinderMergeItemTableHeader
        mergeMode='MERGE_ORG'
        label='Name'
        subLabel={''}
      />
    ),
    template: (o: any) => {
      return <OrganizationTableCell organization={o} />;
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
    width: '45%',
    label: 'Location',
    subLabel: 'City, State, Country',
    template: (organization: any) => {
      return <AddressTableCell locations={organization.locations} />;
    },
  },
  {
    id: 'finder-table-column-actions',
    width: '5%',
    label: <ActionColumn scope={'MERGE_ORG'} />,
    subLabel: '',
    template: () => {
      return <div style={{ display: 'none' }} />;
    },
  },
];
