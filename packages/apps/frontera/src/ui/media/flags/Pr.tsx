import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Pr = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='pr-a'>
        <path fillOpacity='.7' d='M-37.3 0h682.7v512H-37.3z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      clipPath='url(#pr-a)'
      transform='translate(35)scale(.9375)'
    >
      <path fill='#ed0000' d='M-37.3 0h768v512h-768z' />
      <path
        fill='#fff'
        d='M-37.3 102.4h768v102.4h-768zm0 204.8h768v102.4h-768z'
      />
      <path fill='#0050f0' d='m-37.3 0 440.7 255.7L-37.3 511z' />
      <path
        fill='#fff'
        d='M156.4 325.5 109 290l-47.2 35.8 17.6-58.1-47.2-36 58.3-.4 18.1-58 18.5 57.8 58.3.1-46.9 36.3z'
      />
    </g>
  </svg>
);