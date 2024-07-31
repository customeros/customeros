import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const La = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='la-a'>
        <path fillOpacity='.7' d='M0 0h640v480H0z' />
      </clipPath>
    </defs>
    <g fillRule='evenodd' clipPath='url(#la-a)'>
      <path fill='#ce1126' d='M-40 0h720v480H-40z' />
      <path fill='#002868' d='M-40 119.3h720v241.4H-40z' />
      <path
        fill='#fff'
        d='M423.4 240a103.4 103.4 0 1 1-206.8 0 103.4 103.4 0 1 1 206.8 0'
      />
    </g>
  </svg>
);
