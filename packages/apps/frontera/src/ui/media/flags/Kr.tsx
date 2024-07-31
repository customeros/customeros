import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Kr = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='kr-a'>
        <path fillOpacity='.7' d='M-95.8-.4h682.7v512H-95.8z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      clipPath='url(#kr-a)'
      transform='translate(89.8 .4)scale(.9375)'
    >
      <path fill='#fff' d='M-95.8-.4H587v512H-95.8Z' />
      <g transform='rotate(-56.3 361.6 -101.3)scale(10.66667)'>
        <g id='kr-c'>
          <path
            id='kr-b'
            fill='#000001'
            d='M-6-26H6v2H-6Zm0 3H6v2H-6Zm0 3H6v2H-6Z'
          />
          <use y='44' width='100%' height='100%' xlinkHref='#kr-b' />
        </g>
        <path d='M0 17v10' stroke='#fff' />
        <path fill='#cd2e3a' d='M0-12a12 12 0 0 1 0 24Z' />
        <path fill='#0047a0' d='M0-12a12 12 0 0 0 0 24A6 6 0 0 0 0 0Z' />
        <circle r='6' cy='-6' fill='#cd2e3a' />
      </g>
      <g transform='rotate(-123.7 191.2 62.2)scale(10.66667)'>
        <use width='100%' height='100%' xlinkHref='#kr-c' />
        <path stroke='#fff' d='M0-23.5v3M0 17v3.5m0 3v3' />
      </g>
    </g>
  </svg>
);
