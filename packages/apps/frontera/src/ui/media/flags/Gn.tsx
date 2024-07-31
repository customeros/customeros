import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gn = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g strokeWidth='1pt' fillRule='evenodd'>
      <path fill='red' d='M0 0h213.3v480H0z' />
      <path fill='#ff0' d='M213.3 0h213.4v480H213.3z' />
      <path fill='#090' d='M426.7 0H640v480H426.7z' />
    </g>
  </svg>
);
