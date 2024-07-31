import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Se = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 640 480'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#005293' d='M0 0h640v480H0z' />
    <path fill='#fecb00' d='M176 0v192H0v96h176v192h96V288h368v-96H272V0z' />
  </svg>
);
