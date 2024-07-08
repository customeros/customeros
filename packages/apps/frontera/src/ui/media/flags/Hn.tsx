import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Hn = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#18c3df' d='M0 0h640v480H0z' />
    <path fill='#fff' d='M0 160h640v160H0z' />
    <g id='hn-c' fill='#18c3df' transform='translate(320 240)scale(26.66665)'>
      <g id='hn-b'>
        <path id='hn-a' d='m-.3 0 .5.1L0-1z' />
        <use
          xlinkHref='#hn-a'
          width='100%'
          height='100%'
          transform='scale(-1 1)'
        />
      </g>
      <use
        xlinkHref='#hn-b'
        width='100%'
        height='100%'
        transform='rotate(72)'
      />
      <use
        xlinkHref='#hn-b'
        width='100%'
        height='100%'
        transform='rotate(-72)'
      />
      <use
        xlinkHref='#hn-b'
        width='100%'
        height='100%'
        transform='rotate(144)'
      />
      <use
        xlinkHref='#hn-b'
        width='100%'
        height='100%'
        transform='rotate(-144)'
      />
    </g>
    <use
      xlinkHref='#hn-c'
      width='100%'
      height='100%'
      transform='translate(133.3 -42.7)'
    />
    <use
      xlinkHref='#hn-c'
      width='100%'
      height='100%'
      transform='translate(133.3 37.3)'
    />
    <use
      xlinkHref='#hn-c'
      width='100%'
      height='100%'
      transform='translate(-133.3 -42.7)'
    />
    <use
      xlinkHref='#hn-c'
      width='100%'
      height='100%'
      transform='translate(-133.3 37.3)'
    />
  </svg>
);
