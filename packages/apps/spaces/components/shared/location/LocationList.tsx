import React from 'react';
import { Location as LocationItem } from './Location';
import { Button } from '@spaces/atoms/button';
import PlusCircle from '@spaces/atoms/icons/PlusCircle';
import styles from './location-list.module.scss';

interface LocationListProps {
  locations: Array<any>;
  onCreateLocation: () => void;
}

export const LocationList: React.FC<LocationListProps> = ({
  locations,
  onCreateLocation,
}) => {
  return (
    <article className={styles.locations_section}>
      <h1 className={styles.location_header}>Locations</h1>
      <ul className={styles.location_list}>
        {locations.map((location) => (
          <li
            key={`organization-location-${location.id}`}
            className={styles.location_item}
          >
            <LocationItem
              locationId={location.id}
              rawAddress={location?.rawAddress || ''}
            />
          </li>
        ))}
      </ul>

      <div className={styles.button_section}>
        <Button onClick={onCreateLocation} mode='secondary'>
          <PlusCircle height={16} />
          Add location
        </Button>
      </div>
    </article>
  );
};
