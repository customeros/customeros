import styles from './table-cells.module.scss';
import Link from 'next/link';
import React, { ReactNode } from 'react';
import classNames from 'classnames';
import { Highlight } from '../../highlight';

export const TableCell = ({
  label,
  subLabel,
  url,
  className,
}: {
  label: string | ReactNode;
  subLabel?: string | ReactNode;
  url?: string;
  className?: string;
}) => {
  return (
    <div className={classNames(styles.cell, styles.columnCell)}>
      {url ? (
        <Link href={url} className={classNames(styles.link, styles.cellData)}>
          {label}
        </Link>
      ) : (
        <span className={classNames(className, styles.cellData)}>{label}</span>
      )}

      {subLabel && (
        <div className={classNames(styles.subLabel, styles.cellData)}>
          {subLabel}
        </div>
      )}
    </div>
  );
};
export const LinkCell = ({
  label,
  subLabel,
  url,
  className,
  children,
}: {
  label: string | ReactNode;
  subLabel?: string | ReactNode;
  url: string;
  className?: string;
  children?: ReactNode;
}) => {
  return (
    <Link href={url} className={classNames(styles.cell, styles.linkCell)}>
      {children}
      <div className={classNames({ [styles.textContent]: children })}>
        <div className={classNames(className, styles.cellData)}>{label}</div>
        {subLabel && (
          <span className={classNames(styles.subLabel, styles.cellData)}>
            {subLabel}
          </span>
        )}
      </div>
    </Link>
  );
};

export const DashboardTableAddressCell = ({
  country = '',
  region = '',
  locality,
  name,
  highlight = '',
}: {
  country?: string | null;
  region?: string | null;
  locality?: string | null;
  highlight?: string;
  name?: string | null;
}) => {
  return (
    <div className={styles.addressContainer}>
      {name && <Highlight text={name} highlight={highlight} />}

      <div className={styles.addressFields}>
        {locality && (
          <div className={`${styles.addressLocality}`}>
            <Highlight text={locality} highlight={highlight} />
          </div>
        )}

        {locality && (country || region) && <div>,&nbsp;</div>}

        {country || region ? (
          <div className={`${styles.addressRegion}`}>
            <Highlight text={region || ''} highlight={highlight} />{' '}
            {country && ', '}
            <Highlight text={country || ''} highlight={highlight} />
          </div>
        ) : (
          ''
        )}
      </div>
    </div>
  );
};
