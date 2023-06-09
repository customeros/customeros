import React, { PropsWithChildren } from 'react';
import { useSelect } from '@spaces/atoms/select/useSelect';
import styles from './select.module.scss';

export const SelectWrapper = ({ children }: PropsWithChildren) => {
  const { getWrapperProps } = useSelect();

  return (
    <div {...getWrapperProps()} className={styles.dropdownWrapper}>
      {children}
    </div>
  );
};
