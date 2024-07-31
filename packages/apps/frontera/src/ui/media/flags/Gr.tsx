import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gr = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#0d5eaf' fillRule='evenodd' d='M0 0h640v53.3H0z' />
    <path fill='#fff' fillRule='evenodd' d='M0 53.3h640v53.4H0z' />
    <path fill='#0d5eaf' fillRule='evenodd' d='M0 106.7h640V160H0z' />
    <path fill='#fff' fillRule='evenodd' d='M0 160h640v53.3H0z' />
    <path fill='#0d5eaf' d='M0 0h266.7v266.7H0z' />
    <path fill='#0d5eaf' fillRule='evenodd' d='M0 213.3h640v53.4H0z' />
    <path fill='#fff' fillRule='evenodd' d='M0 266.7h640V320H0z' />
    <path fill='#0d5eaf' fillRule='evenodd' d='M0 320h640v53.3H0z' />
    <path fill='#fff' fillRule='evenodd' d='M0 373.3h640v53.4H0z' />
    <g fill='#fff' strokeWidth='1.3' fillRule='evenodd'>
      <path d='M106.7 0H160v266.7h-53.3z' />
      <path d='M0 106.7h266.7V160H0z' />
    </g>
    <path fill='#0d5eaf' d='M0 426.7h640V480H0z' />
  </svg>
);
