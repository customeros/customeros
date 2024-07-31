import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Hn = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#18c3df' d='M0 0h640v480H0z' />
    <path fill='#fff' d='M0 160h640v160H0z' />
    <g id='hn-c' fill='#18c3df' transform='translate(320 240)scale(26.66665)'>
      <g id='hn-b'>
        <path id='hn-a' d='m-.3 0 .5.1L0-1z' />
        <use
          width='100%'
          height='100%'
          xlinkHref='#hn-a'
          transform='scale(-1 1)'
        />
      </g>
      <use
        width='100%'
        height='100%'
        xlinkHref='#hn-b'
        transform='rotate(72)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#hn-b'
        transform='rotate(-72)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#hn-b'
        transform='rotate(144)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#hn-b'
        transform='rotate(-144)'
      />
    </g>
    <use
      width='100%'
      height='100%'
      xlinkHref='#hn-c'
      transform='translate(133.3 -42.7)'
    />
    <use
      width='100%'
      height='100%'
      xlinkHref='#hn-c'
      transform='translate(133.3 37.3)'
    />
    <use
      width='100%'
      height='100%'
      xlinkHref='#hn-c'
      transform='translate(-133.3 -42.7)'
    />
    <use
      width='100%'
      height='100%'
      xlinkHref='#hn-c'
      transform='translate(-133.3 37.3)'
    />
  </svg>
);
