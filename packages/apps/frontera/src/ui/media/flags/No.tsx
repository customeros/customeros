import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const No = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#ed2939' d='M0 0h640v480H0z' />
    <path fill='#fff' d='M180 0h120v480H180z' />
    <path fill='#fff' d='M0 180h640v120H0z' />
    <path fill='#002664' d='M210 0h60v480h-60z' />
    <path fill='#002664' d='M0 210h640v60H0z' />
  </svg>
);
