import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Is = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='is-a'>
        <path fillOpacity='.7' d='M0 0h640v480H0z' />
      </clipPath>
    </defs>
    <g fillRule='evenodd' strokeWidth='0' clipPath='url(#is-a)'>
      <path fill='#003897' d='M0 0h666.7v480H0z' />
      <path
        fill='#fff'
        d='M0 186.7h186.7V0h106.6v186.7h373.4v106.6H293.3V480H186.7V293.3H0z'
      />
      <path
        fill='#d72828'
        d='M0 213.3h213.3V0h53.4v213.3h400v53.4h-400V480h-53.4V266.7H0z'
      />
    </g>
  </svg>
);
