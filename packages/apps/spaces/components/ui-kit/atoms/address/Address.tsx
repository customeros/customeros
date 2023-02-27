import React from 'react';
import styles from './address.module.scss';
// import {Link as AddressInterface} from "../../../models/Link";
//
// interface Props extends Omit<AddressInterface, 'id'> {
//     mode?: 'default' | 'light'
// }
export const Address = ({
  country,
  state,
  city,
  address,
  address2,
  zip,
  phone,
  fax,
  mode = 'default',
}: any) => {
  return (
    <div className={styles.addressContainer}>
      {address && (
        <div className={`${styles.address} ${styles[mode]}`}>{address}</div>
      )}
      {address2 && (
        <div className={`${styles.address} ${styles[mode]}`}>{address2}</div>
      )}

      {(city || state || zip) && (
        <div className={`${styles.address} ${styles[mode]}`}>
          {city}, {state} {zip}
        </div>
      )}

      {country && (
        <div className={`${styles.country} ${styles[mode]}`}>{country}</div>
      )}

      {(phone || fax) && (
        <div className={styles.phoneAndFaxContainer}>
          {phone && (
            <div className={styles.phone}>
              <span> Phone: </span> {phone}
            </div>
          )}
          {fax && (
            <div className={styles.fax}>
              <span> Fax: </span> {fax}
            </div>
          )}
        </div>
      )}
    </div>
  );
};
