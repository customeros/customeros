import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const De = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fc0' d='M0 320h640v160H0z' />
    <path fill='#000001' d='M0 0h640v160H0z' />
    <path fill='red' d='M0 160h640v160H0z' />
  </svg>
);
