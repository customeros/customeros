import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sy = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#000001' d='M0 0h640v480H0Z' />
    <path fill='#fff' d='M0 0h640v320H0Z' />
    <path fill='#ce1126' d='M0 0h640v160H0Z' />
    <path
      fill='#007a3d'
      d='m161 300 39-120 39 120-102-74.2h126M401 300l39-120 39 120-102-74.2h126'
    />
  </svg>
);
