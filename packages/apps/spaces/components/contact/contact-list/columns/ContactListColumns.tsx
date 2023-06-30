import React, { FC } from 'react';
import { Column } from '@spaces/atoms/table/types';
import {
  AddressTableCell,
  ContactTableCell,
  EmailTableCell,
  FinderMergeItemTableHeader,
} from '@spaces/finder/finder-table';
import { LinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { OrganizationAvatar } from '@spaces/molecules/organization-avatar/OrganizationAvatar';
import { ContactActionColumn } from './ContactActionColumn';
import { FinderContactTableSortingState } from '../../../../state/finderTables';
import { useRecoilState } from 'recoil';
import { TableHeaderCell } from '@spaces/atoms/table';
import { SortableCell } from '@spaces/atoms/table/table-cells/SortableCell';
import { Email } from '../../../../graphQL/__generated__/generated';
import { finderContactsGridData } from '../../../../state';

const ContactSortableCell: FC<{
  column: FinderContactTableSortingState['column'];
}> = ({ column }) => {
  const [contactsGridData, setContactsGridData] = useRecoilState(
    finderContactsGridData,
  );
  return (
    <SortableCell
      column={column}
      sort={contactsGridData.sortBy}
      setSortingState={(sortingState: any) => {
        setContactsGridData((prevState: any) => ({
          ...prevState,
          sortBy: sortingState,
        }));
      }}
    />
  );
};

export const contactListColumns: Array<Column> = [
  {
    id: 'finder-table-column-contact-name',
    width: '25%',
    label: (
      <FinderMergeItemTableHeader label='Contact' subLabel='' withIcon>
        <ContactSortableCell column='CONTACT' />
      </FinderMergeItemTableHeader>
    ),
    template: (contact: any) => {
      return <ContactTableCell contact={contact} />;
    },
  },
  {
    id: 'finder-table-column-email',
    width: '20%',
    label: (
      <TableHeaderCell label='Email' subLabel=''>
        <ContactSortableCell column='EMAIL' />
      </TableHeaderCell>
    ),
    template: (c: any) => {
      if (!c) {
        return <span>-</span>;
      }
      return (
        <EmailTableCell emails={c?.emails.filter((e: Email) => !!e.email)} />
      );
    },
  },
  {
    id: 'finder-table-organization-position',
    width: '20%',
    label: (
      <TableHeaderCell label='Organization' subLabel='Position' withIcon>
        <ContactSortableCell column='ORGANIZATION' />
      </TableHeaderCell>
    ),
    template: (c: any) => {
      if (
        !c.jobRoles ||
        c.jobRoles.length === 0 ||
        c.jobRoles[0].organization === null
      ) {
        return <span>-</span>;
      }

      return (
        <LinkCell
          label={c.jobRoles[0].organization?.name || 'Unnamed'}
          subLabel={c.jobRoles[0].jobTitle ?? '-'}
          url={`/organization/${c.jobRoles[0].organization.id ?? undefined}`}
        >
          <OrganizationAvatar
            name={c.jobRoles[0].organization?.name || 'Unnamed'}
          />
        </LinkCell>
      );
    },
  },
  {
    id: 'finder-table-column-org',
    width: '25%',
    label: (
      <TableHeaderCell label='Location' subLabel='Address'>
        <ContactSortableCell column='LOCATION' />
      </TableHeaderCell>
    ),
    isLast: true,
    template: (c: any) => {
      return <AddressTableCell locations={c.locations} />;
    },
  },
  {
    id: 'finder-table-column-actions',
    width: '10%',
    label: <ContactActionColumn />,
    subLabel: '',
    template: () => {
      return <div style={{ display: 'none' }} />;
    },
  },
];
