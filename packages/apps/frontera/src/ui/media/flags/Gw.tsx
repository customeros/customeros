import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gw = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#ce1126' d='M0 0h220v480H0z' />
    <path fill='#fcd116' d='M220 0h420v240H220z' />
    <path fill='#009e49' d='M220 240h420v240H220z' />
    <g id='gw-b' transform='matrix(80 0 0 80 110 240)'>
      <path
        id='gw-a'
        d='M0-1v1h.5'
        fill='#000001'
        transform='rotate(18 0 -1)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#gw-a'
        transform='scale(-1 1)'
      />
    </g>
    <use
      width='100%'
      height='100%'
      xlinkHref='#gw-b'
      transform='rotate(72 110 240)'
    />
    <use
      width='100%'
      height='100%'
      xlinkHref='#gw-b'
      transform='rotate(144 110 240)'
    />
    <use
      width='100%'
      height='100%'
      xlinkHref='#gw-b'
      transform='rotate(-144 110 240)'
    />
    <use
      width='100%'
      height='100%'
      xlinkHref='#gw-b'
      transform='rotate(-72 110 240)'
    />
  </svg>
);
