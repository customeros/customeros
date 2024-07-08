import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Kw = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='kw-a'>
        <path fillOpacity='.7' d='M0 0h682.7v512H0z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      strokeWidth='1pt'
      clipPath='url(#kw-a)'
      transform='scale(.9375)'
    >
      <path fill='#fff' d='M0 170.6h1024v170.7H0z' />
      <path fill='#f31830' d='M0 341.3h1024V512H0z' />
      <path fill='#00d941' d='M0 0h1024v170.7H0z' />
      <path fill='#000001' d='M0 0v512l255.4-170.7.6-170.8z' />
    </g>
  </svg>
);
