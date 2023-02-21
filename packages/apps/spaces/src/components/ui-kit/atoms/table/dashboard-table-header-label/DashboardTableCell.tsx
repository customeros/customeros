import styles from './dashboard-table-header-label.module.scss';
import Link from 'next/link';
import React, { ReactNode } from 'react';

export const DashboardTableCell = ({
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
        <Link href={url} className={styles.link}>
          {label}
        </Link>
      ) : (
        <span className={className}>{label}</span>
      )}

      {subLabel && <span className={styles.subLabel}>{subLabel}</span>}
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
