import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Fo = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <defs>
      <clipPath id='fo-a'>
        <path fillOpacity='.7' d='M-78 32h640v480H-78z' />
      </clipPath>
    </defs>
    <g
      fillRule='evenodd'
      strokeWidth='0'
      clipPath='url(#fo-a)'
      transform='translate(78 -32)'
    >
      <path fill='#fff' d='M-78 32h663.9v480H-78z' />
      <path
        fill='#003897'
        d='M-76 218.7h185.9V32H216v186.7h371.8v106.6H216V512H109.9V325.3h-186z'
      />
      <path
        fill='#d72828'
        d='M-76 245.3h212.4V32h53.1v213.3H588v53.4H189.5V512h-53V298.7H-76z'
      />
    </g>
  </svg>
);
