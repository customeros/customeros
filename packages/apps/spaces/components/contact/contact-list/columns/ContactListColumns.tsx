import React from 'react';
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

// const SortableCell = () => {
//   const [sort, setSortingState] = useRecoilState(
//     finderContactTableSortingState,
//   );
//   return (
//     <TableHeaderCell
//       label='Name'
//       // sortable
//       hasAvatar
//       onSort={(direction: SortingDirection) => {
//         setSortingState({ direction, column: 'NAME' });
//       }}
//       direction={sort.direction}
//     />
//   );
// };

export const contactListColumns: Array<Column> = [
  {
    id: 'finder-table-column-contact-name',
    width: '25%',
    label: <FinderMergeItemTableHeader label='Contact' subLabel='' />,
    template: (contact: any) => {
      return <ContactTableCell contact={contact} />;
    },
  },
  {
    id: 'finder-table-column-email',
    width: '20%',
    label: 'Email',
    template: (c: any) => {
      if (!c?.contact) {
        return <span>-</span>;
      }
      return <EmailTableCell emails={c.contact?.emails} />;
    },
  },
  {
    id: 'finder-table-organization-position',
    width: '20%',
    label: 'Organization',
    subLabel: 'Position',
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
            organizationId={c.jobRoles[0].organization.id}
            name={c.jobRoles[0].organization?.name || 'Unnamed'}
          />
        </LinkCell>
      );
    },
  },
  {
    id: 'finder-table-column-org',
    width: '25%',
    label: 'Location',
    subLabel: 'Address',
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
