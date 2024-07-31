import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Td = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g fillRule='evenodd'>
      <path fill='#002664' d='M0 0h214v480H0z' />
      <path fill='#c60c30' d='M426 0h214v480H426z' />
      <path fill='#fecb00' d='M214 0h212v480H214z' />
    </g>
  </svg>
);
