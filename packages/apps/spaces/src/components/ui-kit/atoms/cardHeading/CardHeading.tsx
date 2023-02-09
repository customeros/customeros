import React, { ReactNode } from 'react';
import styles from './card-heading.module.scss';
interface Props {
  children: ReactNode;
  subheading?: React.ReactNode | string;
}

export const CardHeading: React.FC<Props> = ({ children, subheading }) => {
  return (
    <div className={styles.headingContainer}>
      <h1 className={styles.heading}>{children}</h1>
      {subheading && (
        <span className='mr-3 overflow-hidden text-gray-500 text-sm text-overflow-ellipsis capitalise'>
          {typeof subheading === 'string'
            ? subheading.split('_').join(' ')
            : subheading}
        </span>
      )}
    </div>
  );
};
