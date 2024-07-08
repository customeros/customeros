import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Mm = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fecb00' d='M0 0h640v480H0z' />
    <path fill='#34b233' d='M0 160h640v320H0z' />
    <path fill='#ea2839' d='M0 320h640v160H0z' />
    <g transform='translate(320 256.9)scale(176.87999)'>
      <path id='mm-a' fill='#fff' d='m0-1 .3 1h-.6z' />
      <use
        xlinkHref='#mm-a'
        width='100%'
        height='100%'
        transform='rotate(-144)'
      />
      <use
        xlinkHref='#mm-a'
        width='100%'
        height='100%'
        transform='rotate(-72)'
      />
      <use
        xlinkHref='#mm-a'
        width='100%'
        height='100%'
        transform='rotate(72)'
      />
      <use
        xlinkHref='#mm-a'
        width='100%'
        height='100%'
        transform='rotate(144)'
      />
    </g>
  </svg>
);
