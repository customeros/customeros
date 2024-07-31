import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Dj = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='dj-a'>
        <path fillOpacity='.7' d='M-40 0h682.7v512H-40z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      clipPath='url(#dj-a)'
      transform='translate(37.5)scale(.94)'
    >
      <path fill='#0c0' d='M-40 0h768v512H-40z' />
      <path fill='#69f' d='M-40 0h768v256H-40z' />
      <path fill='#fffefe' d='m-40 0 382.7 255.7L-40 511z' />
      <path
        fill='red'
        d='M119.8 292 89 270l-30.7 22.4L69.7 256l-30.6-22.5 37.9-.3 11.7-36.3 12 36.2h37.9l-30.5 22.7 11.7 36.4z'
      />
    </g>
  </svg>
);
