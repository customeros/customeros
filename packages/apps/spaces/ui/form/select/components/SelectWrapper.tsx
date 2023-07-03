import React, { CSSProperties, PropsWithChildren } from 'react';
import { useSelect } from '../useSelect';
import styles from './select.module.scss';

export const SelectWrapper = ({
  children,
  isHidden,
  customStyles = {},
}: PropsWithChildren<{
  isHidden?: boolean;
  customStyles?: CSSProperties | undefined;
}>) => {
  const { getWrapperProps } = useSelect();

  return (
    <div
      {...getWrapperProps()}
      className={styles.dropdownWrapper}
      style={{ visibility: isHidden ? 'hidden' : 'visible', ...customStyles }}
    >
      {children}
    </div>
  );
};
