import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ps = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='ps-a'>
        <path fillOpacity='.7' d='M-118 0h682.7v512H-118z' />
      </clipPath>
    </defs>
    <g clipPath='url(#ps-a)' transform='translate(110.6)scale(.9375)'>
      <g strokeWidth='1pt' fillRule='evenodd'>
        <path fill='#000001' d='M-246 0H778v170.7H-246z' />
        <path fill='#fff' d='M-246 170.7H778v170.6H-246z' />
        <path fill='#090' d='M-246 341.3H778V512H-246z' />
        <path fill='red' d='m-246 512 512-256L-246 0z' />
      </g>
    </g>
  </svg>
);
