import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sd = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='sd-a'>
        <path fillOpacity='.7' d='M0 0h682.7v512H0z' />
      </clipPath>
    </defs>
    <g
      strokeWidth='1pt'
      fillRule='evenodd'
      clipPath='url(#sd-a)'
      transform='scale(.9375)'
    >
      <path fill='#000001' d='M0 341.3h1024V512H0z' />
      <path fill='#fff' d='M0 170.6h1024v170.7H0z' />
      <path fill='red' d='M0 0h1024.8v170.7H0z' />
      <path fill='#009a00' d='M0 0v512l341.3-256z' />
    </g>
  </svg>
);
