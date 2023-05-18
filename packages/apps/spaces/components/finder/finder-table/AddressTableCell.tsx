import React, { useRef } from 'react';
import { DashboardTableAddressCell } from '@spaces/atoms/table/table-cells/TableCell';
import { Button } from '@spaces/atoms/button';
import { OverlayPanel } from '@spaces/atoms/overlay-panel';
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
        // highlight={searchTern}
      />
    );
  }

  const displayedLocation = locations[0];

  return (
    <div>
      <Button
        mode='text'
        // @ts-expect-error revisit
        onClick={(e) => op?.current?.toggle(e)}
        style={{ padding: 0 }}
      >
        <DashboardTableAddressCell
          key={displayedLocation.id}
          locality={displayedLocation?.locality}
          region={displayedLocation?.region}
          name={displayedLocation?.name}
          {...displayedLocation}
          // highlight={searchTern}
        />
        <span className={styles.showMoreLocationsIcon}>(...)</span>
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
              {...data}
            />
          ))}
        </ul>
      </OverlayPanel>
    </div>
  );
};
