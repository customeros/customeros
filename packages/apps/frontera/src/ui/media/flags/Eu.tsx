import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Eu = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <g id='eu-d'>
        <g id='eu-b'>
          <path id='eu-a' d='m0-1-.3 1 .5.1z' />
          <use xlinkHref='#eu-a' transform='scale(-1 1)' />
        </g>
        <g id='eu-c'>
          <use xlinkHref='#eu-b' transform='rotate(72)' />
          <use xlinkHref='#eu-b' transform='rotate(144)' />
        </g>
        <use xlinkHref='#eu-c' transform='scale(-1 1)' />
      </g>
    </defs>
    <path fill='#039' d='M0 0h640v480H0z' />
    <g fill='#fc0' transform='translate(320 242.3)scale(23.7037)'>
      <use y='-6' width='100%' height='100%' xlinkHref='#eu-d' />
      <use y='6' width='100%' height='100%' xlinkHref='#eu-d' />
      <g id='eu-e'>
        <use x='-6' width='100%' height='100%' xlinkHref='#eu-d' />
        <use
          width='100%'
          height='100%'
          xlinkHref='#eu-d'
          transform='rotate(-144 -2.3 -2.1)'
        />
        <use
          width='100%'
          height='100%'
          xlinkHref='#eu-d'
          transform='rotate(144 -2.1 -2.3)'
        />
        <use
          width='100%'
          height='100%'
          xlinkHref='#eu-d'
          transform='rotate(72 -4.7 -2)'
        />
        <use
          width='100%'
          height='100%'
          xlinkHref='#eu-d'
          transform='rotate(72 -5 .5)'
        />
      </g>
      <use
        width='100%'
        height='100%'
        xlinkHref='#eu-e'
        transform='scale(-1 1)'
      />
    </g>
  </svg>
);
