import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Cm = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#007a5e' d='M0 0h213.3v480H0z' />
    <path fill='#ce1126' d='M213.3 0h213.4v480H213.3z' />
    <path fill='#fcd116' d='M426.7 0H640v480H426.7z' />
    <g fill='#fcd116' transform='translate(320 240)scale(7.1111)'>
      <g id='cm-b'>
        <path id='cm-a' d='M0-8-2.5-.4 1.3.9z' />
        <use
          width='100%'
          height='100%'
          xlinkHref='#cm-a'
          transform='scale(-1 1)'
        />
      </g>
      <use
        width='100%'
        height='100%'
        xlinkHref='#cm-b'
        transform='rotate(72)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#cm-b'
        transform='rotate(144)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#cm-b'
        transform='rotate(-144)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#cm-b'
        transform='rotate(-72)'
      />
    </g>
  </svg>
);
