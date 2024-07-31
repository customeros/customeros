import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gg = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#fff' d='M0 0h640v480H0z' />
    <path fill='#e8112d' d='M256 0h128v480H256z' />
    <path fill='#e8112d' d='M0 176h640v128H0z' />
    <path
      id='gg-a'
      fill='#f9dd16'
      d='m110 286.7 23.3-23.4h210v-46.6h-210L110 193.3z'
    />
    <use
      width='36'
      height='24'
      xlinkHref='#gg-a'
      transform='rotate(90 320 240)'
    />
    <use
      width='36'
      height='24'
      xlinkHref='#gg-a'
      transform='rotate(-90 320 240)'
    />
    <use
      width='36'
      height='24'
      xlinkHref='#gg-a'
      transform='rotate(180 320 240)'
    />
  </svg>
);
