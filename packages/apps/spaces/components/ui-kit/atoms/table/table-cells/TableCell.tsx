import styles from './table-cells.module.scss';
import Link from 'next/link';
import React, { ReactNode } from 'react';
import classNames from 'classnames';

export const TableCell = ({
  label,
  subLabel,
  url,
  className,
}: {
  label: string;
  subLabel?: string | ReactNode;
  url?: string;
  className?: string;
}) => {
  return (
    <div className={styles.cell}>
      {url ? (
        <Link
          href={url}
          className={classNames(styles.link, styles.cellData)}
          title={label}
        >
          {label}
        </Link>
      ) : (
        <span title={label} className={classNames(className, styles.cellData)}>
          {label}
        </span>
      )}

      {subLabel && (
        <span className={classNames(styles.subLabel, styles.cellData)}>
          {subLabel}
        </span>
      )}
    </div>
  );
};

export const DashboardTableAddressCell = ({
  country,
  region,
  locality,
}: {
  country?: string | null;
  region?: string | null;
  locality?: string | null;
}) => {
  return (
    <div className={styles.addressContainer}>
      {locality && (
        <div className={`${styles.addressLocality}`}>{locality}</div>
      )}

      {(country || region) && (
        <div className={`${styles.addressRegion}`}>
          {region}, {country}
        </div>
      )}
    </div>
  );
};
