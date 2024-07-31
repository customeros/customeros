import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Jo = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='jo-a'>
        <path fillOpacity='.7' d='M-117.8 0h682.6v512h-682.6z' />
      </clipPath>
    </defs>
    <g clipPath='url(#jo-a)' transform='translate(110.5)scale(.9375)'>
      <g strokeWidth='1pt' fillRule='evenodd'>
        <path fill='#000001' d='M-117.8 0h1024v170.7h-1024z' />
        <path fill='#fff' d='M-117.8 170.7h1024v170.6h-1024z' />
        <path fill='#090' d='M-117.8 341.3h1024V512h-1024z' />
        <path fill='red' d='m-117.8 512 512-256-512-256z' />
        <path
          fill='#fff'
          d='m24.5 289 5.7-24.9H4.7l23-11-15.9-19.9 23 11 5.6-24.8 5.7 24.9L69 233.2l-16 19.9 23 11H50.6l5.7 24.9-15.9-20z'
        />
      </g>
    </g>
  </svg>
);
