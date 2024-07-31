import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Ve = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <g id='ve-d' transform='translate(0 -36)'>
        <g id='ve-c'>
          <g id='ve-b'>
            <path id='ve-a' fill='#fff' d='M0-5-1.5-.2l2.8.9z' />
            <use
              width='180'
              height='120'
              xlinkHref='#ve-a'
              transform='scale(-1 1)'
            />
          </g>
          <use
            width='180'
            height='120'
            xlinkHref='#ve-b'
            transform='rotate(72)'
          />
        </g>
        <use
          width='180'
          height='120'
          xlinkHref='#ve-b'
          transform='rotate(-72)'
        />
        <use
          width='180'
          height='120'
          xlinkHref='#ve-c'
          transform='rotate(144)'
        />
      </g>
    </defs>
    <path fill='#cf142b' d='M0 0h640v480H0z' />
    <path fill='#00247d' d='M0 0h640v320H0z' />
    <path fill='#fc0' d='M0 0h640v160H0z' />
    <g id='ve-f' transform='matrix(4 0 0 4 320 336)'>
      <g id='ve-e'>
        <use
          width='180'
          height='120'
          xlinkHref='#ve-d'
          transform='rotate(10)'
        />
        <use
          width='180'
          height='120'
          xlinkHref='#ve-d'
          transform='rotate(30)'
        />
      </g>
      <use width='180' height='120' xlinkHref='#ve-e' transform='rotate(40)' />
    </g>
    <use
      width='180'
      height='120'
      xlinkHref='#ve-f'
      transform='rotate(-80 320 336)'
    />
  </svg>
);
