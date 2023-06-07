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
import { Organization } from '@spaces/graphql';
import {
  FinderOrganizationTableSortingState,
  finderOrganizationTableSortingState,
} from '../../../../state/finderTables';
import { useRecoilState } from 'recoil';
import { SortableCell } from '@spaces/atoms/table/table-cells/SortableCell';
import { OwnerTableCell } from '@spaces/finder/finder-table/OwnerTableCell';
import { OrganizationRelationship } from '../../organization-details/relationship/OrganizationRelationship';
import { RelationshipStage } from '../../organization-details/stage/RelationshipStage';

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

export const organizationListColumns: Array<Column<Organization>> = [
  {
    id: 'finder-table-column-organization-name',
    width: '20%',
    label: (
      <FinderMergeItemTableHeader label='Company' subLabel='Branch'>
        <OrganizationSortableCell column='ORGANIZATION' />
      </FinderMergeItemTableHeader>
    ),
    template: (organization) => {
      return <OrganizationTableCell organization={organization} />;
    },
  },
  {
    id: 'finder-table-column-organization-relationship',
    width: '15%',
    label: (
      <FinderMergeItemTableHeader label='Relationship' subLabel='Stage'>
        <OrganizationSortableCell column='RELATIONSHIP' />
      </FinderMergeItemTableHeader>
    ),
    template: (organization) => (
      <>
        <OrganizationRelationship
          organizationId={organization.id}
          defaultValue={organization.relationships?.[0]}
        />
        <RelationshipStage defaultValue={'ACTIVE'} />
      </>
    ),
  },
  {
    id: 'finder-table-column-domain-website',
    width: '15%',
    label: (
      <TableHeaderCell label='Domain' subLabel='Website'>
        <OrganizationSortableCell column='DOMAIN' />
      </TableHeaderCell>
    ),

    template: (organization) => {
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
    width: '20%',
    label: (
      <TableHeaderCell label='Location' subLabel='Address'>
        <OrganizationSortableCell column='LOCATION' />
      </TableHeaderCell>
    ),
    template: (organization) => {
      return <AddressTableCell locations={organization.locations} />;
    },
  },
  {
    id: 'finder-table-column-organization-owner',
    width: '20%',
    label: (
      <TableHeaderCell label='Owner'>
        <OrganizationSortableCell column='OWNER' />
      </TableHeaderCell>
    ),
    isLast: true,
    template: (organization) => {
      return (
        <OwnerTableCell
          owner={organization.owner}
          organizationId={organization.id}
        />
      );
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
