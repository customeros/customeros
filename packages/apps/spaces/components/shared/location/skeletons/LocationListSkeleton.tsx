import React from 'react';
import styles from '../location-list.module.scss';
import { Skeleton } from '@spaces/atoms/skeleton';

export const LocationListSkeleton: React.FC = () => {
  return (
    <article className={styles.locations_section}>
      <h1 className={styles.location_header}>Locations</h1>
      <ul className={styles.location_list}>
        {[1, 2].map((location) => (
          <li
            key={`organization-location-skeleton-${location}`}
            className={styles.location_item}
          >
            <Skeleton width='180px' height='10px' />
          </li>
        ))}
      </ul>
    </article>
  );
};
