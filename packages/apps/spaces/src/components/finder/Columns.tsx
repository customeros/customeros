import {
  DashboardTableAddressCell,
  DashboardTableCell,
} from '../ui-kit/atoms/table/dashboard-table-header-label/DashboardTableCell';
import React from 'react';
import { LocationFragment } from '../../graphQL/generated';
import { Column } from '../ui-kit/atoms/table/types';

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
          <DashboardTableCell
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
        <DashboardTableCell
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
      return (c.contact?.emails || []).map((data: any, index: number) => (
        <div style={{ display: 'flex', flexWrap: 'wrap' }} key={data.id}>
          <DashboardTableCell
            className='lowercase'
            label={data.email}
            url={`/contact/${data.id}`}
          />
          {c?.contact?.emails.length - 1 !== index && ', '}
        </div>
      ));
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

      return c?.contact?.locations.map((data: LocationFragment) => (
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
