import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gm = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='gm-a'>
        <path fillOpacity='.7' d='M0-48h640v480H0z' />
      </clipPath>
    </defs>
    <g
      strokeWidth='1pt'
      fillRule='evenodd'
      clipPath='url(#gm-a)'
      transform='translate(0 48)'
    >
      <path fill='red' d='M0-128h640V85.3H0z' />
      <path fill='#fff' d='M0 85.3h640V121H0z' />
      <path fill='#009' d='M0 120.9h640V263H0z' />
      <path fill='#fff' d='M0 263.1h640v35.6H0z' />
      <path fill='#090' d='M0 298.7h640V512H0z' />
    </g>
  </svg>
);
