import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const GbSct = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#0065bd' d='M0 0h640v480H0z' />
    <path
      stroke='#fff'
      strokeWidth='.6'
      d='m0 0 5 3M0 3l5-3'
      transform='scale(128 160)'
    />
  </svg>
);
