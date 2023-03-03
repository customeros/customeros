import React, { useRef } from 'react';
import { DashboardTableAddressCell } from '../../ui-kit/atoms/table/table-cells/TableCell';
import { Button } from '../../ui-kit';
import { OverlayPanel } from '../../ui-kit/atoms/overlay-panel';
import styles from './finder-table.module.scss';

export const AddressTableCell = ({
  locations = [],
}: {
  locations: Array<any>;
}) => {
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
