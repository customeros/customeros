import React from 'react';
import { Location as LocationItem } from './Location';
import { Button } from '@spaces/atoms/button';
import PlusCircle from '@spaces/atoms/icons/PlusCircle';
import styles from './location-list.module.scss';
import { Location as COSLocation } from '@spaces/graphql';

type TLocation = Omit<
  COSLocation,
  'appSource' | 'source' | 'sourceOfTruth' | 'createdAt' | 'updatedAt'
>;
interface LocationListProps {
  locations: Array<TLocation>;
  onCreateLocation: () => void;
  onRemoveLocation: (locationId: string) => void;
  isEditMode: boolean;
}

const getLocationString = (location: TLocation) => {
  if (location.rawAddress) {
    return location.rawAddress;
  }
  const addressComponents = [
    location?.country || '',
    location?.zip || location?.postalCode || '',
    location?.street || '',
    location?.houseNumber || '',
  ];

  const formattedAddress = addressComponents.map((item, index) =>
    index !== 0 && item ? `, ${item}` : item,
  );

  return formattedAddress.join('').trim();
};
export const LocationList: React.FC<LocationListProps> = ({
  locations,
  onCreateLocation,
  onRemoveLocation,
  isEditMode,
}) => {
  return (
    <article className={styles.locations_section}>
      <h1 className={styles.location_header}>Locations</h1>
      {!locations.length && !isEditMode && (
        <div className={styles.location_item}>
          This company has no locations
        </div>
      )}
      <ul className={styles.location_list}>
        {locations.map((location) => (
          <li
            key={`organization-location-${location.id}`}
            className={styles.location_item}
          >
            <LocationItem
              isEditMode={isEditMode}
              locationId={location.id}
              onRemoveLocation={onRemoveLocation}
              locationString={getLocationString(location)}
            />
          </li>
        ))}
      </ul>

      {isEditMode && (
        <div className={styles.button_section}>
          <Button onClick={onCreateLocation} mode='secondary'>
            <PlusCircle height={16} />
            Add location
          </Button>
        </div>
      )}
    </article>
  );
};
