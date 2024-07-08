import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ax = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='ax-a'>
        <path fillOpacity='.7' d='M106.3 0h1133.3v850H106.3z' />
      </clipPath>
    </defs>
    <g clipPath='url(#ax-a)' transform='matrix(.56472 0 0 .56482 -60 -.1)'>
      <path fill='#0053a5' d='M0 0h1300v850H0z' />
      <g fill='#ffce00'>
        <path d='M400 0h250v850H400z' />
        <path d='M0 300h1300v250H0z' />
      </g>
      <g fill='#d21034'>
        <path d='M475 0h100v850H475z' />
        <path d='M0 375h1300v100H0z' />
      </g>
    </g>
  </svg>
);
