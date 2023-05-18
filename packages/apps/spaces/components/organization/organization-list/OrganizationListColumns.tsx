import React from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableCell, TableHeaderCell } from '@spaces/atoms/table';
import {
  ActionColumn,
  AddressTableCell,
  ContactTableCell,
  FinderCell,
  FinderMergeItemTableHeader,
  OrganizationTableCell,
} from '@spaces/finder/finder-table';
import { LinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { OrganizationAvatar } from '@spaces/molecules/organization-avatar/OrganizationAvatar';
import { useRecoilState } from 'recoil';
import { finderOrganizationTableSortingState } from '../../../state/finderOrganizationTable';
import { SortingDirection } from '../../../graphQL/__generated__/generated';

const SortableCell = () => {
  const [sort, setSortingState] = useRecoilState(
    finderOrganizationTableSortingState,
  );
  return (
    <TableHeaderCell
      label='Company'
      subLabel='Branch'
      // sortable
      hasAvatar
      onSort={(direction: SortingDirection) => {
        setSortingState({ direction, column: 'NAME' });
      }}
      direction={sort.direction}
    />
  );
};

export const organizationListColumns: Array<Column> = [
  {
    id: 'finder-table-column-organization-name',
    width: '25%',
    label: <SortableCell />,
    template: (organization: any) => {
      const hasSubsidiaries = !!organization.subsidiaries?.length;
      if (hasSubsidiaries) {
        return (
          <LinkCell
            label={organization.subsidiaries[0].organization.name || 'Unnamed'}
            subLabel={organization.name}
            url={`/organization/${organization.id}`}
          >
            {<OrganizationAvatar organizationId={organization.id} />}
          </LinkCell>
        );
      }
      return (
        <LinkCell
          label={organization.name || 'Unnamed'}
          subLabel={''}
          url={`/organization/${organization.id}`}
        >
          {<OrganizationAvatar organizationId={organization.id} />}
        </LinkCell>
      );
    },
  },
  {
    id: 'finder-table-column-domain-website',
    width: '25%',
    label: <TableHeaderCell label='Domain' subLabel='Website' />,

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
    isLast:true,
    template: (organization: any) => {
      return <AddressTableCell  locations={organization.locations} />;
    },
  },
  // {
  //   id: 'finder-table-column-actions',
  //   width: '5%',
  //   label: <ActionColumn scope={'MERGE_ORG'} />,
  //   subLabel: '',
  //   template: () => {
  //     return <div style={{ display: 'none' }} />;
  //   },
  // },
];
