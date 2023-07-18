'use client';
import { createColumnHelper } from '@tanstack/react-table';

import { Table } from '@ui/presentation/_Table/Table';
import { THead } from '@ui/presentation/_Table/THead';
import { MenuItem } from '@ui/overlay/Menu';
import { Organization } from '@graphql/types';

type OrgTableRow = Pick<Organization, 'name' | 'industry' | 'domain'>;

const columnHelper = createColumnHelper<OrgTableRow>();

const columns = [
  columnHelper.accessor('name', {
    header: (props) => (
      <THead<OrgTableRow> title='Name' subTitle='SubName' {...props} />
    ),
    cell: (props) => <>{props.getValue()}</>,
    skeleton: () => <div>loading...</div>,
  }),
  columnHelper.accessor('industry', {
    header: (props) => (
      <THead<OrgTableRow> title='Industry' subTitle='Subindustry' {...props} />
    ),
    cell: (props) => <>{props.getValue()}</>,
    skeleton: () => <div>loading...</div>,
  }),
  columnHelper.accessor('domain', {
    header: (props) => (
      <THead<OrgTableRow> title='Domain' subTitle='Subdomain' {...props} />
    ),
    cell: (props) => <>{props.getValue()}</>,
    skeleton: () => <div>loading...</div>,
  }),
];

const mock: OrgTableRow[] = [
  {
    name: 'Org 1',
    industry: 'Industry 1',
    domain: 'Source 1',
  },
  {
    name: 'Org 2',
    industry: 'Industry 2',
    domain: 'Source 2',
  },
  {
    name: 'Org 3',
    industry: 'Industry 3',
    domain: 'Source 3',
  },
  {
    name: 'Org 4',
    industry: 'Industry 4',
    domain: 'Source 4',
  },
  {
    name: 'Org 5',
    industry: 'Industry 5',
    domain: 'Source 5',
  },
];

export default function OrganizationsPage() {
  return (
    <>
      <Table<OrgTableRow>
        data={mock}
        columns={columns}
        isLoading={false}
        totalItems={5}
        enableTableActions
        enableRowSelection
        renderTableActions={() => {
          return (
            <>
              <MenuItem>Item 1</MenuItem>
              <MenuItem>Item 2</MenuItem>
              <MenuItem>Item 3</MenuItem>
            </>
          );
        }}
      />
    </>
  );
}
