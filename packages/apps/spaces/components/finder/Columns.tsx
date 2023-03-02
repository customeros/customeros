import {
  DashboardTableAddressCell,
  TableCell,
} from '../ui-kit/atoms/table/table-cells/TableCell';
import React, { ReactNode, useRef } from 'react';
import { Column } from '../ui-kit/atoms/table/types';
import { useRecoilValue } from 'recoil';
import { finderSearchTerm } from '../../state';
import { Button, Highlight } from '../ui-kit';
import { OverlayPanel } from '../ui-kit/atoms/overlay-panel';
import styles from './finder.module.scss';

const FinderCell = ({
  label,
  subLabel,
  url,
}: {
  label: string;
  subLabel?: string | ReactNode;
  url?: string;
}) => {
  const searchTern = useRecoilValue(finderSearchTerm);

  return (
    <TableCell
      label={<Highlight text={label} highlight={searchTern} />}
      subLabel={subLabel}
      url={url}
    />
  );
};

const EmailCell = ({ emails }: any) => {
  const searchTern = useRecoilValue(finderSearchTerm);
  const op = useRef(null);

  if (!emails) {
    return <span>-</span>;
  }

  if (emails.length === 1) {
    return <FinderCell label={emails[0]?.email} />;
  }
  const primary = (emails || []).find((data: any) =>
    searchTern ? data?.email?.includes(searchTern) : data.primary,
  );

  return (
    <div>
      <Button
        role='button'
        mode='text'
        style={{ padding: 0 }}
        // @ts-expect-error revisit
        onClick={(e) => op?.current?.toggle(e)}
      >
        <FinderCell label={primary?.email || emails[0]?.email} />
        <span style={{ marginLeft: '8px' }}>(...)</span>
      </Button>
      <OverlayPanel
        ref={op}
        style={{
          maxHeight: '400px',
          height: 'fit-content',
          overflowX: 'hidden',
          overflowY: 'auto',
          bottom: 0,
        }}
      >
        <ul className={styles.adressesList}>
          {emails
            .filter((d: any) => !!d?.email)
            .map((data: any) => (
              <li
                key={data.id}
                style={{ display: 'flex' }}
                className={styles.emailList}
              >
                <FinderCell label={data.email} />
              </li>
            ))}
        </ul>
      </OverlayPanel>
    </div>
  );
};

const AdressesCell = ({ locations = [] }: { locations: Array<any> }) => {
  const op = useRef(null);

  const locationsCount: number | undefined = locations.length;
  if (!locationsCount) {
    return <span>-</span>;
  }

  if (locationsCount === 1) {
    return (
      <DashboardTableAddressCell
        key={locations[0].id}
        locality={locations[0]?.locality}
        region={locations[0]?.region}
        name={locations[0]?.name}
      />
    );
  }

  return (
    <div>
      <Button
        mode='text'
        // @ts-expect-error revisit
        onClick={(e) => op?.current?.toggle(e)}
        style={{ padding: 0 }}
      >
        <DashboardTableAddressCell
          key={locations[0].id}
          locality={locations[0]?.locality}
          region={locations[0]?.region}
          name={locations[0]?.name}
        />
        <span style={{ marginLeft: '8px' }}>(...)</span>
      </Button>
      <OverlayPanel
        ref={op}
        style={{
          maxHeight: '400px',
          height: 'fit-content',
          overflowX: 'hidden',
          overflowY: 'auto',
          bottom: 0,
        }}
      >
        <ul className={styles.adressesList}>
          {locations.map((data) => (
            <DashboardTableAddressCell
              key={data.id}
              locality={data?.locality}
              region={locations[0]?.region}
              name={data?.name}
            />
          ))}
        </ul>
      </OverlayPanel>
    </div>
  );
};

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
          <FinderCell
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
      const name = `${c?.contact?.firstName} ${c?.contact?.lastName} ${c?.contact?.name}`;
      const displayName = name.trim().length ? name : 'Unnamed';
      return (
        <FinderCell
          label={displayName}
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
      if (!c?.contact) {
        return <span>-</span>;
      }
      return <EmailCell emails={c.contact?.emails} />;
    },
  },
  {
    width: '25%',
    label: 'Location',
    subLabel: 'City, State, Country',
    template: (c: any) => {
      return <AdressesCell locations={c?.contact?.locations} />;
    },
  },
];
