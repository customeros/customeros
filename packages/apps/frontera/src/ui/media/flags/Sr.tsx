import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Sr = ({ className, ...props }: IconProps) => (
  <svg
    viewBox='0 0 640 480'
    fill='none'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path fill='#377e3f' d='M.1 0h640v480H.1z' />
    <path fill='#fff' d='M.1 96h640v288H.1z' />
    <path fill='#b40a2d' d='M.1 144h640v192H.1z' />
    <path
      fill='#ecc81d'
      d='m320 153.2 56.4 173.6-147.7-107.3h182.6L263.6 326.8z'
    />
  </svg>
);
