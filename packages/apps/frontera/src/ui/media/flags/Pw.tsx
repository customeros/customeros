import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Pw = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='pw-a'>
        <path fillOpacity='.7' d='M-70.3 0h640v480h-640z' />
      </clipPath>
    </defs>
    <g
      strokeWidth='1pt'
      fillRule='evenodd'
      clipPath='url(#pw-a)'
      transform='translate(70.3)'
    >
      <path fill='#4aadd6' d='M-173.4 0h846.3v480h-846.3z' />
      <path
        fill='#ffde00'
        d='M335.6 232.1a135.9 130.1 0 1 1-271.7 0 135.9 130.1 0 1 1 271.7 0'
      />
    </g>
  </svg>
);
