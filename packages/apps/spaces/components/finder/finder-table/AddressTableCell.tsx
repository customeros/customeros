import React, { useRef } from 'react';
import { DashboardTableAddressCell } from '@spaces/atoms/table/table-cells/TableCell';
import { OverlayPanel } from '@spaces/atoms/overlay-panel';
import styles from './finder-table.module.scss';

export const AddressTableCell = ({
  locations = [],
}: {
  locations: Array<any>;
}) => {
  const op = useRef(null);

  const locationsCount: number | undefined = locations.length;
  const displayedLocation = locations[0];
  const hiddenLocations = locations.filter(
    (loc) =>
      loc.id !== displayedLocation.id &&
      (loc.rawAddress ||
        loc.region ||
        loc.country ||
        loc.street ||
        loc.postalCode ||
        loc.houseNumber ||
        loc.zip ||
        loc.name),
  );
  if (!locationsCount) {
    return <span>-</span>;
  }

  if (hiddenLocations.length === 0) {
    return (
      <DashboardTableAddressCell
        key={locations[0].id}
        locality={locations[0]?.locality}
        region={locations[0]?.region}
        name={locations[0]?.name}
        rawAddress={locations[0]?.rawAddress}
        // highlight={searchTern}
      />
    );
  }

  return (
    <div className={styles.addressCellWrapper}>
      <div
        role='button'
        //@ts-expect-error ignore
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
        >
          <span className={styles.showMoreLocationsIcon}>
            +{hiddenLocations.length}
          </span>
        </DashboardTableAddressCell>
      </div>

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
          {hiddenLocations.map((data) => (
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
