import React, { ReactNode } from 'react';
import styles from './ribbon.module.scss';

interface RibbonProps {
  children: ReactNode;
  top?: number;
}

export const Ribbon: React.FC<RibbonProps> = ({ children }) => {
  return <div className={styles.ribbon}>{children}</div>;
};
