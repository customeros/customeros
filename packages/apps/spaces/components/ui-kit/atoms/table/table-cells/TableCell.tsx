import styles from './table-cells.module.scss';
import Link from 'next/link';
import React, { ReactNode } from 'react';
import classNames from 'classnames';

export const TableCell = ({
  label,
  subLabel,
  children,
  className,
}: {
  label: string | ReactNode;
  subLabel?: string | ReactNode;
  className?: string;
  children?: ReactNode;
}) => {
  return (
    <div className={classNames(styles.cell)}>
      {children}

      <div
        className={classNames({ [styles.textContent]: children })}
        style={{ width: '100%' }}
      >
        <span className={classNames(className, styles.cellData)}>{label}</span>
        {subLabel && (
          <span className={classNames(styles.subLabel, styles.cellData)}>
            {subLabel}
          </span>
        )}
      </div>
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
        <span className={classNames(className, styles.cellData)}>{label}</span>
        {subLabel && (
          <span className={classNames(styles.subLabel, styles.cellData)}>
            {subLabel}
          </span>
        )}
      </div>
    </Link>
  );
};

export const ExternalLinkCell = ({
  url,
  className,
}: {
  url: string;
  className?: string;
}) => {
  const createSocialLink = (link: string) => {
    if (link.includes('http')) {
      return link;
    } else return 'https://' + link;
  };
  return (
    <a
      href={createSocialLink(url)}
      rel='noopener noreferrer'
      target='_blank'
      className={classNames(styles.cell, styles.linkCell)}
    >
      <span className={classNames(className, styles.cellData)}>{url}</span>
    </a>
  );
};

export const DashboardTableAddressCell = ({
  country = '',
  region = '',
  locality,
  name,
  street,
  postalCode,
  zip,
  houseNumber,
  rawAddress,
  children,
}: {
  country?: string | null;
  region?: string | null;
  locality?: string | null;
  zip?: string | null;
  postalCode?: string | null;
  houseNumber?: string | null;
  rawAddress?: string | null;
  street?: string | null;
  highlight?: string;
  name?: string | null;
  children?: ReactNode;
}) => {
  if (rawAddress) {
    return (
      <div className={classNames(styles.addressContainer, styles.rawAddress)}>
        {rawAddress}
        {children}
      </div>
    );
  }

  return (
    <div className={styles.addressContainer}>
      {name && name}

      <div className={styles.addressFields}>
        {locality && (
          <div className={`${styles.addressLocality}`}>{locality}</div>
        )}

        {locality && (country || region) && <div>,&nbsp;</div>}

        <div className={`${styles.addressRegion}`}>
          {region && region}
          {(zip || postalCode) && `, ${zip || postalCode}`}
          {country && `, ${country}`}
        </div>
      </div>

      {(street || houseNumber) && (
        <div className={`${styles.addressRegion}`}>
          {street} {houseNumber}
        </div>
      )}

      {children}
    </div>
  );
};
