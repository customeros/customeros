import React from 'react';
import { Column } from '@spaces/atoms/table/types';
import { TableCell, TableHeaderCell } from '@spaces/atoms/table';
import { AddressTableCell, EmailTableCell } from '@spaces/finder/finder-table';
import { uuid4 } from '@sentry/utils';
import { useRecoilState } from 'recoil';
import { finderContactTableSortingState } from '../../../state/finderOrganizationTable';
import { SortingDirection } from '../../../graphQL/__generated__/generated';
import { LinkCell } from '@spaces/atoms/table/table-cells/TableCell';
import { getContactDisplayName } from '../../../utils';
import { ContactAvatar } from '@spaces/molecules/contact-avatar/ContactAvatar';

const SortableCell = () => {
  const [sort, setSortingState] = useRecoilState(
    finderContactTableSortingState,
  );
  return (
    <TableHeaderCell
      label='Name'
      // sortable
      hasAvatar
      onSort={(direction: SortingDirection) => {
        setSortingState({ direction, column: 'NAME' });
      }}
      direction={sort.direction}
    />
  );
};

export const contactListColumns: Array<Column> = [
  {
    id: 'finder-table-column-contact-name',
    width: '25%',
    label: <SortableCell />,
    template: (contact: any) => {
      return (
        <LinkCell
          label={getContactDisplayName(contact)}
          url={`/contact/${contact.id}`}
        >
          {<ContactAvatar contactId={contact.id} />}
        </LinkCell>
      );
    },
  },
  {
    id: 'finder-table-column-email',
    width: '20%',
    label: 'Email',
    template: (c: any) => {
      if (!c?.contact) {
        return <span key={uuid4()}>-</span>;
      }
      return <EmailTableCell emails={c.contact?.emails} />;
    },
  },
  {
    id: 'finder-table-organization-position',
    width: '25%',
    label: 'Organization',
    subLabel: 'Position',
    template: (c: any) => {
      if (!c.jobRoles || c.jobRoles.length === 0) {
        return <span>-</span>;
      }
      return (
        <TableCell
          label={c.jobRoles[0].organization?.name ?? 'Unnamed'}
          subLabel={c.jobRoles[0].jobTitle ?? '-'}
          url={`/organization/${c.jobRoles[0].organization?.id ?? undefined}`}
        />
      );
    },
  },
  {
    id: 'finder-table-column-org',
    width: '25%',
    label: 'Location',
    subLabel: 'City, State, Country',
    isLast: true,
    template: (c: any) => {
      return <AddressTableCell locations={c?.contact?.locations} />;
    },
  },
  // {
  //   id: 'finder-table-column-actions',
  //   width: '5%',
  //   label: <ActionColumn scope={'MERGE_CONTACT'} />,
  //   subLabel: '',
  //   template: () => {
  //     return <div style={{ display: 'none' }} />;
  //   },
  // },
];
