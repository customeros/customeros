import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ee = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#1791ff' d='M0 0h640v160H0z' />
    <path fill='#000001' d='M0 160h640v160H0z' />
    <path fill='#fff' d='M0 320h640v160H0z' />
  </svg>
);
