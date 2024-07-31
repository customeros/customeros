import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Jp = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='jp-a'>
        <path fillOpacity='.7' d='M-88 32h640v480H-88z' />
      </clipPath>
    </defs>
    <g
      strokeWidth='1pt'
      fillRule='evenodd'
      clipPath='url(#jp-a)'
      transform='translate(88 -32)'
    >
      <path fill='#fff' d='M-128 32h720v480h-720z' />
      <circle
        r='194.9'
        cx='523.1'
        cy='344.1'
        fill='#bc002d'
        transform='translate(-168.4 8.6)scale(.76554)'
      />
    </g>
  </svg>
);
