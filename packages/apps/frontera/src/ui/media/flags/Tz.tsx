import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Tz = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='tz-a'>
        <path fillOpacity='.7' d='M10 0h160v120H10z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      strokeWidth='1pt'
      clipPath='url(#tz-a)'
      transform='matrix(4 0 0 4 -40 0)'
    >
      <path fill='#09f' d='M0 0h180v120H0z' />
      <path fill='#090' d='M0 0h180L0 120z' />
      <path fill='#000001' d='M0 120h40l140-95V0h-40L0 95z' />
      <path
        fill='#ff0'
        d='M0 91.5 137.2 0h13.5L0 100.5zM29.3 120 180 19.5v9L42.8 120z'
      />
    </g>
  </svg>
);
