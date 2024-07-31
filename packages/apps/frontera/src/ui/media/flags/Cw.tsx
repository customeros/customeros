import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cw = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='cw-a'>
        <path fillOpacity='.7' d='M0 0h682.7v512H0z' />
      </clipPath>
      <path id='cw-b' d='m0-1 .2.7H1L.3 0l.2.7L0 .4l-.6.4.2-.7-.5-.4h.7z' />
    </defs>
    <g clipPath='url(#cw-a)' transform='scale(.94)'>
      <path fill='#002b7f' d='M0 0h768v512H0z' />
      <path fill='#f9e814' d='M0 320h768v64H0z' />
      <use
        x='2'
        y='2'
        fill='#fff'
        width='13500'
        height='9000'
        xlinkHref='#cw-b'
        transform='scale(42.67)'
      />
      <use
        x='3'
        y='3'
        fill='#fff'
        width='13500'
        height='9000'
        xlinkHref='#cw-b'
        transform='scale(56.9)'
      />
    </g>
  </svg>
);
