import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const St = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#12ad2b' d='M0 0h640v480H0z' />
    <path fill='#ffce00' d='M0 137.1h640V343H0z' />
    <path fill='#d21034' d='M0 0v480l240-240' />
    <g id='st-c' transform='translate(351.6 240)scale(.34286)'>
      <g id='st-b'>
        <path
          id='st-a'
          fill='#000001'
          d='M0-200V0h100'
          transform='rotate(18 0 -200)'
        />
        <use
          width='100%'
          height='100%'
          xlinkHref='#st-a'
          transform='scale(-1 1)'
        />
      </g>
      <use
        width='100%'
        height='100%'
        xlinkHref='#st-b'
        transform='rotate(72)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#st-b'
        transform='rotate(144)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#st-b'
        transform='rotate(-144)'
      />
      <use
        width='100%'
        height='100%'
        xlinkHref='#st-b'
        transform='rotate(-72)'
      />
    </g>
    <use
      x='700'
      width='100%'
      height='100%'
      xlinkHref='#st-c'
      transform='translate(-523.2)'
    />
  </svg>
);
