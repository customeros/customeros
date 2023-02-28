import React, { FC, ReactNode } from 'react';
import styles from './link.module.scss';
import NextJSLink, { LinkProps } from 'next/link';
export const Link: FC<LinkProps & { children: ReactNode }> = ({
  href,
  children,
  ...rest
}) => {
  return (
    <NextJSLink href={href} className={styles.link} {...rest}>
      {children}
    </NextJSLink>
  );
};
