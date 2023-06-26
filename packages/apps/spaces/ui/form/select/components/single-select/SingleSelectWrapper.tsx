import React, { PropsWithChildren } from 'react';
import styles from '../select.module.scss';
import { useSingleSelect } from './SingleSelect';

export const SingleSelectWrapper = ({
  children,
  isHidden,
}: PropsWithChildren<{ isHidden?: boolean }>) => {
  const { getWrapperProps } = useSingleSelect();

  return (
    <div
      {...getWrapperProps()}
      className={styles.dropdownWrapper}
      style={{ visibility: isHidden ? 'hidden' : 'visible' }}
    >
      {children}
    </div>
  );
};
