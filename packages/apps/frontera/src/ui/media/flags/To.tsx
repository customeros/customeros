import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const To = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g strokeWidth='1pt' fillRule='evenodd'>
      <path fill='#c10000' d='M0 0h640v480H0z' />
      <path fill='#fff' d='M0 0h250v200.3H0z' />
      <g fill='#c10000'>
        <path d='M102.8 31.2h39.9v139.6h-39.8z' />
        <path d='M192.6 81v40H53V81z' />
      </g>
    </g>
  </svg>
);
