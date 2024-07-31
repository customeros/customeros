import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const In = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#f93' d='M0 0h640v160H0z' />
    <path fill='#fff' d='M0 160h640v160H0z' />
    <path fill='#128807' d='M0 320h640v160H0z' />
    <g transform='matrix(3.2 0 0 3.2 320 240)'>
      <circle r='20' fill='#008' />
      <circle r='17.5' fill='#fff' />
      <circle r='3.5' fill='#008' />
      <g id='in-d'>
        <g id='in-c'>
          <g id='in-b'>
            <g id='in-a' fill='#008'>
              <circle r='.9' transform='rotate(7.5 -8.8 133.5)' />
              <path d='M0 17.5.6 7 0 2l-.6 5z' />
            </g>
            <use
              width='100%'
              height='100%'
              xlinkHref='#in-a'
              transform='rotate(15)'
            />
          </g>
          <use
            width='100%'
            height='100%'
            xlinkHref='#in-b'
            transform='rotate(30)'
          />
        </g>
        <use
          width='100%'
          height='100%'
          xlinkHref='#in-c'
          transform='rotate(60)'
        />
      </g>
      <use
        width='100%'
        height='100%'
        xlinkHref='#in-d'
        transform='rotate(120)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#in-d'
        transform='rotate(-120)'
      />
    </g>
  </svg>
);
