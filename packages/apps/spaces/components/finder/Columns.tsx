import {
  DashboardTableAddressCell,
  TableCell,
} from '../ui-kit/atoms/table/table-cells/TableCell';
import React from 'react';
import { Column } from '../ui-kit/atoms/table/types';
import { LocationBaseDetailsFragment } from '../../graphQL/__generated__/generated';

export const columns: Array<Column> = [
  {
    width: '25%',
    label: 'Organization',
    subLabel: 'Industry',
    template: (c: any) => {
      if (c.organization) {
        const industry = (
          <span className={'capitalise'}>
            {c.organization?.industry.split('_').join(' ').toLowerCase()}
          </span>
        );
        return (
          <TableCell
            label={c.organization.name}
            subLabel={industry}
            url={`/organization/${c.organization.id}`}
          />
        );

        return c.organization.name;
      }
      if (!c?.organization) {
        return <span>-</span>;
      }
    },
  },
  {
    width: '25%',
    label: 'Name',
    subLabel: 'Role',

    template: (c: any) => {
      if (!c?.contact) {
        return <span>-</span>;
      }
      return (
        <TableCell
          label={`${c?.contact?.firstName} ${c?.contact?.lastName}`}
          subLabel={c?.contact.job}
          url={`/contact/${c?.contact.id}`}
        />
      );
    },
  },
  {
    width: '25%',
    label: 'Email',
    template: (c: any) => {
      if (!c?.contact?.emails) {
        return <span>-</span>;
      }
      return (c.contact?.emails || []).map((data: any, index: number) => {
        const label =
          c?.contact?.emails.length - 1 !== index
            ? `${data.email},`
            : data.email;
        return (
          <div style={{ display: 'flex', flexWrap: 'wrap' }} key={data.id}>
            <TableCell className='lowercase' label={label} />
          </div>
        );
      });
    },
  },
  {
    width: '25%',
    label: 'Location',
    subLabel: 'City, State, Country',
    template: (c: any) => {
      if (!c?.contact?.locations?.length) {
        return <span>-</span>;
      }

      return c?.contact?.locations.map((data: LocationBaseDetailsFragment) => (
        <DashboardTableAddressCell
          key={data.id}
          locality={data?.locality}
          region={data?.region}
          country={data?.country}
        />
      ));
    },
  },
];
