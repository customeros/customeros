import React from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableHeaderCell } from '@spaces/atoms/table';
import {
  AddressTableCell,
  FinderMergeItemTableHeader,
  OrganizationTableCell,
} from '@spaces/finder/finder-table';
import { LinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { OrganizationActionColumn } from './OrganizationActionColumn';
//
// const SortableCell = () => {
//   const [sort, setSortingState] = useRecoilState(
//     finderOrganizationTableSortingState,
//   );
//   return (
//     <TableHeaderCell
//       label='Company'
//       subLabel='Branch'
//       // sortable
//       hasAvatar
//       onSort={(direction: SortingDirection) => {
//         setSortingState({ direction, column: 'NAME' });
//       }}
//       direction={sort.direction}
//     />
//   );
// };

export const organizationListColumns: Array<Column> = [
  {
    id: 'finder-table-column-organization-name',
    width: '25%',
    label: <FinderMergeItemTableHeader label='Company' subLabel='Branch' />,
    template: (organization: any) => {
      return <OrganizationTableCell organization={organization} />;
    },
  },
  {
    id: 'finder-table-column-domain-website',
    width: '25%',
    label: <TableHeaderCell label='Domain' subLabel='Website' />,

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
    width: '40%',
    label: 'Location',
    subLabel: 'City, State, Country',
    isLast: true,
    template: (organization: any) => {
      return <AddressTableCell locations={organization.locations} />;
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
