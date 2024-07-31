import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Bv = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='bv-a'>
        <path fillOpacity='.7' d='M0 0h640v480H0z' />
      </clipPath>
    </defs>
    <g strokeWidth='1pt' fillRule='evenodd' clipPath='url(#bv-a)'>
      <path fill='#fff' d='M-28 0h699.7v512H-28z' />
      <path
        fill='#d72828'
        d='M-53-77.8h218.7v276.2H-53zM289.4-.6h381v199h-381zM-27.6 320h190.4v190.3H-27.6zm319.6 2.1h378.3v188.2H292z'
      />
      <path fill='#003897' d='M196.7-25.4H261v535.7h-64.5z' />
      <path fill='#003897' d='M-27.6 224.8h698v63.5h-698z' />
    </g>
  </svg>
);
