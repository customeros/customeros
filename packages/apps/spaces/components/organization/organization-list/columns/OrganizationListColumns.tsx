import React, { FC } from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableHeaderCell } from '@spaces/atoms/table';
import {
  AddressTableCell,
  FinderMergeItemTableHeader,
  OrganizationTableCell,
} from '@spaces/finder/finder-table';
import { LinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { OrganizationActionColumn } from './OrganizationActionColumn';
import {
  FinderOrganizationTableSortingState,
  finderOrganizationTableSortingState,
} from '../../../../state/finderTables';
import { useRecoilState } from 'recoil';
import { SortableCell } from '@spaces/atoms/table/table-cells/SortableCell';
import { ContactAvatar } from '@spaces/molecules/contact-avatar/ContactAvatar';
import { OrganizationAvatar } from '@spaces/molecules/organization-avatar/OrganizationAvatar';
import { OwnerTableCell } from '@spaces/finder/finder-table/OwnerTableCell';

const OrganizationSortableCell: FC<{
  column: FinderOrganizationTableSortingState['column'];
}> = ({ column }) => {
  const [sort, setSortingState] = useRecoilState(
    finderOrganizationTableSortingState,
  );
  return (
    <SortableCell
      column={column}
      sort={sort}
      setSortingState={setSortingState}
    />
  );
};

export const organizationListColumns: Array<Column> = [
  {
    id: 'finder-table-column-organization-name',
    width: '25%',
    label: (
      <FinderMergeItemTableHeader label='Company' subLabel='Branch'>
        <OrganizationSortableCell column='ORGANIZATION' />
      </FinderMergeItemTableHeader>
    ),
    template: (organization: any) => {
      return <OrganizationTableCell organization={organization} />;
    },
  },
  {
    id: 'finder-table-column-domain-website',
    width: '25%',
    label: (
      <TableHeaderCell label='Domain' subLabel='Website'>
        <OrganizationSortableCell column='DOMAIN' />
      </TableHeaderCell>
    ),

    template: (organization: any) => {
      return (
        <LinkCell
          label={organization.domain}
          subLabel={organization.website}
          url={`/organization/${organization.id}`}
        />
      );
    },
  },
  {
    id: 'finder-table-column-address',
    width: '25%',
    label: (
      <TableHeaderCell label='Location' subLabel='Address'>
        <OrganizationSortableCell column='LOCATION' />
      </TableHeaderCell>
    ),
    template: (organization: any) => {
      return <AddressTableCell locations={organization.locations} />;
    },
  },
  {
    id: 'finder-table-column-organization-owner',
    width: '15%',
    label: (
      <TableHeaderCell label='Owner'>
        <OrganizationSortableCell column='OWNER' />
      </TableHeaderCell>
    ),
    isLast: true,
    template: (organization: any) => {
      return <OwnerTableCell owner={organization.owner} />;
    },
  },
  {
    id: 'finder-table-column-actions',
    width: '10%',
    label: <OrganizationActionColumn />,
    subLabel: '',
    template: () => {
      return <div style={{ display: 'none' }} />;
    },
  },
];
