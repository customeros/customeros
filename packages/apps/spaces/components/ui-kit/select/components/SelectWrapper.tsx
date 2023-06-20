import React, { PropsWithChildren } from 'react';
import { useSelect } from '../useSelect';
import styles from './select.module.scss';

export const SelectWrapper = ({
  children,
  isHidden,
}: PropsWithChildren<{ isHidden?: boolean }>) => {
  const { getWrapperProps } = useSelect();

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
