import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Am = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#d90012' d='M0 0h640v160H0z' />
    <path fill='#0033a0' d='M0 160h640v160H0z' />
    <path fill='#f2a800' d='M0 320h640v160H0z' />
  </svg>
);
